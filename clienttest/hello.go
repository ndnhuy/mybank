package main

import (
	"fmt"
	"log"
)

const baseURL = "http://localhost:8080"

// AccountInfo represents the account information returned by the API
type AccountInfo struct {
	ID      string  `json:"id"`
	Balance float64 `json:"balance"`
}

// TransferRequest represents the transfer request payload
type TransferRequest struct {
	FromAccountID string  `json:"fromAccountId"`
	ToAccountID   string  `json:"toAccountId"`
	Amount        float64 `json:"amount"`
}

// assertBalance checks if the account balance matches the expected value
func assertBalance(user *User, expectedBalance float64) error {
	latestBalance, err := user.GetAccountBalance()
	if err != nil {
		return fmt.Errorf("failed to get account balance: %w", err)
	}
	if latestBalance != expectedBalance {
		return fmt.Errorf("assertion failed: expected balance %.2f, got %.2f for account %s",
			expectedBalance, latestBalance, user.AccountId)
	}

	log.Printf("✓ Account %s has correct balance: %.2f", user.AccountId, latestBalance)
	return nil
}

func main() {
	log.Println("Starting banking API test scenario...")

	userA := NewUser(100.0)
	userB := NewUser(100.0)

	scenario := &TransferScenario{
		FromUser: userA,
		ToUser:   userB,
		Amount:   30.5,
	}
	err := scenario.Run()
	if err != nil {
		log.Fatalf("Transfer failed: %v", err)
	}

	// Assert balances
	for _, user := range []*User{userA, userB} {
		latestBalance, err := user.GetAccountBalance()
		if err != nil {
			log.Fatalf("failed to get account balance: %v", err)
		}
		expectedBalance := user.ApplyActions(scenario.actions).CurrentBalance
		if latestBalance != expectedBalance {
			log.Fatalf("assertion failed: expected balance %.2f, got %.2f for account %s",
				expectedBalance, latestBalance, user.AccountId)
		}

		log.Printf("✓ Account %s has correct balance: %.2f", user.AccountId, latestBalance)
	}

	log.Println("🎉 All assertions passed! Test scenario completed successfully.")
}
