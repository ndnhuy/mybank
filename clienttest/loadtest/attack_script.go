package loadtest

import (
	"fmt"
	"os"
	"time"

	"com.ndnhuy.mybank/domain"
	"com.ndnhuy.mybank/utils"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func AttackGetAccounts(rps, testDuration int) {
	fmt.Printf("Starting load test: %d RPS for %d seconds\n", rps, testDuration)
	fmt.Printf("Target URL: http://localhost:8080/accounts\n")
	fmt.Printf("Press Ctrl+C to stop early if needed\n")
	fmt.Printf("Tip: Set RPS=50 DURATION=60 to customize load parameters\n\n")

	queueMetrics := NewQueueMetrics()
	fmt.Printf("Attack in progress...")

	// Create and use Attacker instance
	attacker := NewAttacker("http://localhost:8080/accounts", "GET", rps, testDuration, queueMetrics.Metrics)
	attacker.Attack()
	queueMetrics.Close()
	fmt.Printf(" completed!\n\n")

	// Print enhanced metrics report
	queueMetrics.PrintReport()

	// Save detailed report to file
	fmt.Printf("\n=== Detailed Report ===\n")
	reportFile, err := os.OpenFile("accounts_loadtest_report.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to open report file: %v\n", err)
		return
	}
	defer reportFile.Close()

	timestamp := fmt.Sprintf("==== Run at %s ===\n", time.Now().Format("2006-01-02 15:04:05"))
	reportFile.WriteString(timestamp)

	reporter := vegeta.NewTextReporter(queueMetrics.Metrics)
	reporter(reportFile)
	reportFile.WriteString("\n\n")
	fmt.Printf("Report appended to accounts_loadtest_report.txt\n")
}

// AttackTransfers simulates simultaneous money transfers between customers
func AttackTransfers(rps, testDuration int) {
	fmt.Printf("Starting transfer attack: %d RPS for %d seconds\n", rps, testDuration)
	fmt.Printf("Setting up test customers...\n")

	// Setup test customers
	sourceCustomers, destCustomers, initialTotal, err := setupTransferCustomers()
	if err != nil {
		fmt.Printf("Failed to setup customers: %v\n", err)
		return
	}
	defer cleanupTransferCustomers(append(sourceCustomers, destCustomers...))

	fmt.Printf("Created %d source customers and %d destination customers\n", len(sourceCustomers), len(destCustomers))
	fmt.Printf("Total initial balance: %.2f\n", initialTotal)
	fmt.Printf("Target URL: %s/accounts/transfer\n", utils.BASE_URL)
	fmt.Printf("Press Ctrl+C to stop early if needed\n\n")

	queueMetrics := NewQueueMetrics()
	fmt.Printf("Transfer attack in progress...")

	// Create customer-based transfer attacker
	transferTargeter := NewCustomerTransferTargeter(sourceCustomers, destCustomers)
	attacker := &Attacker{
		targeter:                 transferTargeter.targeter,
		customerTransferTargeter: transferTargeter,
		rate:                     vegeta.Rate{Freq: rps, Per: time.Second},
		duration:                 time.Duration(testDuration) * time.Second,
		attacker:                 vegeta.NewAttacker(),
		metrics:                  queueMetrics.Metrics,
	}

	attacker.Attack()
	queueMetrics.Close()
	fmt.Printf(" completed!\n\n")

	finalTotal := verifyCustomerBalances(append(sourceCustomers, destCustomers...))
	fmt.Printf("Final total balance: %.2f\n", finalTotal)
	if abs(finalTotal-initialTotal) < 0.01 {
		fmt.Printf("✅ Balance verification passed - no money lost or created\n")
	} else {
		fmt.Printf("❌ Balance verification failed - money discrepancy: %.2f\n", finalTotal-initialTotal)
	}

	// Print enhanced metrics report
	queueMetrics.PrintReport()

	// Save detailed report to file
	fmt.Printf("\n=== Detailed Report ===\n")
	reportFile, err := os.OpenFile("transfer_attack_report.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to open report file: %v\n", err)
		return
	}
	defer reportFile.Close()

	timestamp := fmt.Sprintf("==== Transfer Attack Run at %s ===\n", time.Now().Format("2006-01-02 15:04:05"))
	reportFile.WriteString(timestamp)
	reportFile.WriteString(fmt.Sprintf("Initial Balance: %.2f, Final Balance: %.2f\n", initialTotal, finalTotal))

	reporter := vegeta.NewTextReporter(queueMetrics.Metrics)
	reporter(reportFile)
	reportFile.WriteString("\n\n")
	fmt.Printf("Report appended to transfer_attack_report.txt\n")
}

// setupTransferCustomers creates test customers for transfer attacks
func setupTransferCustomers() (sourceCustomers, destCustomers []*domain.Customer, totalBalance float64, err error) {
	const numSourceCustomers = 100
	const numDestCustomers = 100
	const initialBalance = 1000.0

	// Create source customers with money
	for i := 0; i < numSourceCustomers; i++ {
		customer, err := domain.NewCustomerWithAmount(fmt.Sprintf("source-%d", i), initialBalance)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("failed to create source customer %d: %w", i, err)
		}
		sourceCustomers = append(sourceCustomers, customer)
		totalBalance += initialBalance
	}

	// Create destination customers with minimal money
	for i := 0; i < numDestCustomers; i++ {
		customer, err := domain.NewCustomerWithAmount(fmt.Sprintf("dest-%d", i), 1)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("failed to create dest customer %d: %w", i, err)
		}
		destCustomers = append(destCustomers, customer)
		totalBalance += 1
	}

	return sourceCustomers, destCustomers, totalBalance, nil
}

func verifyCustomerBalances(customers []*domain.Customer) float64 {
	total := 0.0
	allVerified := true

	for _, customer := range customers {
		err := customer.VerifyBalance()
		if err != nil {
			fmt.Printf("⚠️  %v\n", err)
			allVerified = false
		}

		// Get current balance for total calculation
		balance, err := customer.GetCurrentBalance()
		if err != nil {
			fmt.Printf("⚠️  Failed to get balance for %s: %v\n", customer.GetName(), err)
			continue
		}
		total += balance
	}

	if allVerified {
		fmt.Printf("✅ All customer balances verified successfully\n")
	} else {
		fmt.Printf("⚠️  Some customer balance verifications failed\n")
	}

	return total
}

// cleanupTransferCustomers logs customer info for cleanup (accounts would need manual cleanup)
func cleanupTransferCustomers(customers []*domain.Customer) {
	accountIDs := make([]string, len(customers))
	for i, customer := range customers {
		accountIDs[i] = customer.GetAccountID()
	}

	fmt.Printf("\nTest accounts created: %v\n", accountIDs)
	fmt.Printf("(Manual cleanup may be required if DELETE endpoint is not available)\n")
}

// abs returns absolute value of float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
