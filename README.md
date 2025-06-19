# Go Stress Test

A high-performance HTTP load testing CLI application written in Go that performs stress testing on web services by spawning concurrent workers to send HTTP requests and collect performance metrics.

## Features

- **Concurrent Load Testing**: Distributes HTTP requests across multiple concurrent workers
- **Minimal Dependencies**: Uses only Go's standard `net/http` package
- **Comprehensive Metrics**: Collects total execution time, request counts, and HTTP status codes
- **Docker Support**: Includes Dockerfile for containerized execution
- **Precise Request Distribution**: Ensures exactly the specified number of requests are sent
- **Thread-Safe Metrics**: Safe concurrent metric collection across workers

## Installation

### Prerequisites
- Go 1.19 or higher
- Docker (optional, for containerized usage)

### Build from Source

```bash
# Clone the repository
git clone <repository-url>
cd go-stress-test

# Initialize Go module (if needed)
go mod init go-stress-test

# Build the binary
go build -o stress-test main.go
```

## Usage

### Command Line Options

```bash
./stress-test --url=<URL> --requests=<N> --concurrency=<N>
```

**Required Flags:**
- `--url`: Target URL to test
- `--requests`: Total number of requests to send
- `--concurrency`: Number of concurrent workers

### Examples

```bash
# Basic load test
./stress-test --url=https://example.com --requests=1000 --concurrency=10

# High concurrency test
./stress-test --url=https://api.example.com/health --requests=5000 --concurrency=100

# Development server test
./stress-test --url=http://localhost:3000 --requests=500 --concurrency=5
```

### Run without Building

```bash
go run main.go --url=https://example.com --requests=1000 --concurrency=10
```

## Docker Usage

### Build Docker Image

```bash
docker build -t stress-test .
```

### Run with Docker

```bash
docker run stress-test --url=https://example.com --requests=1000 --concurrency=10
```

## Output

The application provides detailed metrics including:
- Total execution time
- Total requests sent
- HTTP status code distribution
- Success/failure rates
- Performance statistics

## Exit Codes

- `0`: Successful execution
- Non-zero: Error occurred (invalid arguments, network issues, etc.)

## Development

### Project Structure

The application is architected around:
- CLI argument parsing and validation
- Worker pool management for concurrent execution
- HTTP client configuration using `net/http`
- Thread-safe metrics collection system
- Report generation and formatting
- Graceful error handling

### Development Commands

```bash
# Format code
go fmt ./...

# Vet code for issues
go vet ./...

# Run tests
go test ./...

# Run specific test
go test -run <TestName>
```

### Key Design Considerations

- **Request Distribution**: Ensures exactly the specified total number of requests
- **Worker Coordination**: Prevents over/under-sending requests through proper synchronization
- **Thread Safety**: All metrics collection is safe across concurrent workers
- **HTTP/1.1 Default**: Uses standard library with HTTP/1.1 by default
- **Minimal Dependencies**: Keeps external dependencies to a minimum

## Contributing

1. Ensure code follows Go formatting standards (`go fmt`)
2. Run tests before submitting (`go test ./...`)
3. Vet code for issues (`go vet ./...`)
4. Follow the existing architecture patterns

## License

[Add your license information here]