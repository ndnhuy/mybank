package main

import (
	"fmt"
	"os"
	"strconv"

	"com.ndnhuy.mybank/loadtest"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

const (
	DEFAULT_RPS      = 10 // Default requests per second
	DEFAULT_DURATION = 30 // Default duration in seconds
)

// getConfigFromEnv reads RPS and DURATION from environment variables or uses defaults
func getConfigFromEnv() (int, int) {
	rps := DEFAULT_RPS
	duration := DEFAULT_DURATION

	if envRPS := os.Getenv("RPS"); envRPS != "" {
		if parsed, err := strconv.Atoi(envRPS); err == nil && parsed > 0 {
			rps = parsed
		}
	}

	if envDuration := os.Getenv("DURATION"); envDuration != "" {
		if parsed, err := strconv.Atoi(envDuration); err == nil && parsed > 0 {
			duration = parsed
		}
	}

	return rps, duration
}

func main() {
	rps, testDuration := getConfigFromEnv()

	fmt.Printf("Starting load test: %d RPS for %d seconds\n", rps, testDuration)
	fmt.Printf("Target URL: http://localhost:8080/accounts\n")
	fmt.Printf("Press Ctrl+C to stop early if needed\n")
	fmt.Printf("Tip: Set RPS=50 DURATION=60 to customize load parameters\n\n")

	var metrics vegeta.Metrics
	fmt.Printf("Attack in progress...")

	// Create and use Attacker instance
	attacker := loadtest.NewAttacker("http://localhost:8080/accounts", "GET", rps, testDuration, &metrics)
	attacker.Attack()
	metrics.Close()
	fmt.Printf(" completed!\n\n")

	// Print key metrics prominently
	fmt.Printf("=== Load Test Results ===\n")
	fmt.Printf("RPS: %d\n", rps)
	fmt.Printf("Duration: %s\n", attacker.Duration())
	fmt.Printf("Total Requests: %d\n", metrics.Requests)
	fmt.Printf("Success Rate: %.2f%%\n", metrics.Success*100)
	fmt.Printf("\n=== Response Time Metrics ===\n")
	fmt.Printf("Mean Response Time: %s\n", metrics.Latencies.Mean)
	fmt.Printf("50th percentile (median): %s\n", metrics.Latencies.P50)
	fmt.Printf("95th percentile: %s\n", metrics.Latencies.P95)
	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
	fmt.Printf("Max Response Time: %s\n", metrics.Latencies.Max)

	// Check for errors
	if metrics.Success < 1.0 {
		fmt.Printf("\n⚠️  Warning: Success rate is %.2f%%. Some requests failed!\n", metrics.Success*100)
	}

	fmt.Printf("\n=== Detailed Report ===\n")
	reporter := vegeta.NewTextReporter(&metrics)
	reporter(os.Stdout)
}
