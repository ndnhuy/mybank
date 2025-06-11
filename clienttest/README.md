# MyBank Load Testing Tool

## Usage

### Basic Usage (Default: 10 RPS for 30 seconds)
```bash
go run main.go
```

### Custom Load Parameters
```bash
# Test with 50 RPS for 60 seconds
RPS=50 DURATION=60 go run main.go

# Light load test: 5 RPS for 10 seconds
RPS=5 DURATION=10 go run main.go

# Heavy load test: 100 RPS for 2 minutes
RPS=100 DURATION=120 go run main.go
```

## Sample Output

```
=== Load Test Results ===
RPS: 10
Duration: 30s
Total Requests: 300
Success Rate: 100.00%

=== Response Time Metrics ===
Mean Response Time: 7.643627ms
50th percentile (median): 7.17289ms
95th percentile: 12.193975ms
99th percentile: 14.211151ms
Max Response Time: 14.959921ms
```

## Understanding the Metrics

- **Mean Response Time**: Average response time across all requests (primary metric)
- **50th percentile (median)**: Half of requests were faster than this
- **95th percentile**: 95% of requests were faster than this
- **99th percentile**: 99% of requests were faster than this
- **Success Rate**: Percentage of requests that returned HTTP 200

## Load Testing Best Practices

1. **Warm-up**: Run a short test first to warm up the service
2. **Realistic Load**: Start with expected production load (e.g., 10-50 RPS)
3. **Gradual Increase**: Double the load to see how performance degrades
4. **Monitor Resources**: Watch CPU, memory, and database metrics during tests
5. **Consistent Environment**: Run tests from the same machine/network for consistency