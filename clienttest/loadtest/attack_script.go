package loadtest

import (
	"fmt"
	"os"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func AttackGetAccounts(rps, testDuration int) {
	fmt.Printf("Starting load test: %d RPS for %d seconds\n", rps, testDuration)
	fmt.Printf("Target URL: http://localhost:8080/accounts\n")
	fmt.Printf("Press Ctrl+C to stop early if needed\n")
	fmt.Printf("Tip: Set RPS=50 DURATION=60 to customize load parameters\n\n")

	var metrics vegeta.Metrics
	fmt.Printf("Attack in progress...")

	// Create and use Attacker instance
	attacker := NewAttacker("http://localhost:8080/accounts", "GET", rps, testDuration, &metrics)
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

	reportFile, err := os.OpenFile("accounts_loadtest_report.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to open report file: %v\n", err)
		return
	}
	defer reportFile.Close()

	timestamp := fmt.Sprintf("==== Run at %s ===\n", time.Now().Format("2006-01-02 15:04:05"))
	reportFile.WriteString(timestamp)

	reporter := vegeta.NewTextReporter(&metrics)
	reporter(reportFile)
	reportFile.WriteString("\n\n")
	fmt.Printf("Report appended to accounts_loadtest_report.txt\n")
}
