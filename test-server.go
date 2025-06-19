package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

var requestCount int64

func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/slow", handleSlow)
	http.HandleFunc("/count", handleCount)
	http.HandleFunc("/status/", handleStatus)

	fmt.Println("Test server starting on :8080")
	fmt.Println("Endpoints:")
	fmt.Println("  /          - Returns 200 OK")
	fmt.Println("  /slow      - 100ms delay")
	fmt.Println("  /count     - Request counter")
	fmt.Println("  /status/X  - Returns status code X")
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	count := atomic.AddInt64(&requestCount, 1)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Request #%d - Hello from test server!", count)
}

func handleSlow(w http.ResponseWriter, r *http.Request) {
	time.Sleep(100 * time.Millisecond)
	count := atomic.AddInt64(&requestCount, 1)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Slow request #%d completed", count)
}

func handleCount(w http.ResponseWriter, r *http.Request) {
	count := atomic.LoadInt64(&requestCount)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Total requests processed: %d", count)
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	// Extract status code from URL path
	statusCode := 200
	if len(r.URL.Path) > 8 {
		switch r.URL.Path[8:] {
		case "404":
			statusCode = 404
		case "500":
			statusCode = 500
		case "503":
			statusCode = 503
		}
	}
	
	atomic.AddInt64(&requestCount, 1)
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "Status: %d", statusCode)
}