package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type Config struct {
	URL         string
	Requests    int
	Concurrency int
}

type Metrics struct {
	TotalRequests   int64
	StatusCounts    map[int]int64
	StartTime       time.Time
	EndTime         time.Time
	mutex           sync.Mutex
}

func main() {
	config := parseFlags()
	if err := validateConfig(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	metrics := &Metrics{
		StatusCounts: make(map[int]int64),
		StartTime:    time.Now(),
	}

	if err := runLoadTest(config, metrics); err != nil {
		fmt.Fprintf(os.Stderr, "Error during load test: %v\n", err)
		os.Exit(1)
	}

	printReport(metrics)
}

func parseFlags() *Config {
	config := &Config{}
	
	flag.StringVar(&config.URL, "url", "", "URL of the service to test")
	flag.IntVar(&config.Requests, "requests", 0, "Total number of HTTP requests to send")
	flag.IntVar(&config.Concurrency, "concurrency", 0, "Number of concurrent workers sending requests")
	
	flag.Parse()
	
	return config
}

func validateConfig(config *Config) error {
	if config.URL == "" {
		return fmt.Errorf("--url is required")
	}
	if config.Requests <= 0 {
		return fmt.Errorf("--requests must be greater than 0")
	}
	if config.Concurrency <= 0 {
		return fmt.Errorf("--concurrency must be greater than 0")
	}
	if config.Concurrency > config.Requests {
		return fmt.Errorf("--concurrency cannot be greater than --requests")
	}
	return nil
}

func runLoadTest(config *Config, metrics *Metrics) error {
	var wg sync.WaitGroup
	requestsPerWorker := config.Requests / config.Concurrency
	remainingRequests := config.Requests % config.Concurrency
	
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	for i := 0; i < config.Concurrency; i++ {
		wg.Add(1)
		workerRequests := requestsPerWorker
		if i < remainingRequests {
			workerRequests++
		}
		
		go func(requests int) {
			defer wg.Done()
			worker(config.URL, requests, client, metrics)
		}(workerRequests)
	}
	
	wg.Wait()
	metrics.EndTime = time.Now()
	
	return nil
}

func worker(url string, requests int, client *http.Client, metrics *Metrics) {
	for i := 0; i < requests; i++ {
		resp, err := client.Get(url)
		
		atomic.AddInt64(&metrics.TotalRequests, 1)
		
		if err != nil {
			metrics.mutex.Lock()
			metrics.StatusCounts[0]++
			metrics.mutex.Unlock()
			continue
		}
		
		metrics.mutex.Lock()
		metrics.StatusCounts[resp.StatusCode]++
		metrics.mutex.Unlock()
		
		resp.Body.Close()
	}
}

func printReport(metrics *Metrics) {
	duration := metrics.EndTime.Sub(metrics.StartTime)
	
	fmt.Printf("Load Test Results:\n")
	fmt.Printf("==================\n")
	fmt.Printf("Total elapsed time: %v\n", duration)
	fmt.Printf("Total requests sent: %d\n", metrics.TotalRequests)
	
	metrics.mutex.Lock()
	if count, exists := metrics.StatusCounts[200]; exists {
		fmt.Printf("Successful responses (200): %d\n", count)
	} else {
		fmt.Printf("Successful responses (200): 0\n")
	}
	
	fmt.Printf("\nStatus Code Distribution:\n")
	for status, count := range metrics.StatusCounts {
		if status == 0 {
			fmt.Printf("  Failed requests (network errors): %d\n", count)
		} else {
			fmt.Printf("  %d: %d\n", status, count)
		}
	}
	metrics.mutex.Unlock()
}