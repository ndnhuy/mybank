package domain

import (
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransferMoneyFromAtoB(t *testing.T) {
	log.Println("Starting banking API test scenario...")

	userA := NewUser(100.0, "UserA")
	userB := NewUser(100.0, "UserB")

	scenario := &TransferScenario{
		FromUser: userA,
		ToUser:   userB,
		Amount:   30.5,
	}
	err := scenario.Run()
	assert.NoError(t, err, "Transfer should succeed")

	// Assert balances
	for _, user := range []BankOperator{userA, userB} {
		latestBalance, err := user.GetAccountBalance()
		assert.NoError(t, err, "Should be able to get account balance for user %s", user.GetAccountId())

		expectedBalance := user.GetExpectedBalance(scenario.actions)
		assert.Equal(t, expectedBalance, latestBalance,
			"Account %s should have correct balance", user.GetAccountId())

		log.Printf("✓ Account %s has correct balance: %.2f", user.GetAccountId(), latestBalance)
	}

	log.Println("🎉 All assertions passed! Test scenario completed successfully.")
}

func TestTransferMoneyFromAtoBtoC(t *testing.T) {
	log.Println("Starting banking API test scenario...")

	userA := NewUser(100.0, "UserA")
	userB := NewUser(100.0, "UserB")
	userC := NewUser(100.0, "UserC")

	scenario1 := &TransferScenario{
		Name:     "Transfer from A to B, amount=30.5",
		FromUser: userA,
		ToUser:   userB,
		Amount:   30,
	}
	scenario2 := &TransferScenario{
		Name:     "Transfer from B to C, amount=20.0",
		FromUser: userB,
		ToUser:   userC,
		Amount:   20.0,
	}

	err := scenario1.Run()
	assert.NoError(t, err, "Scenario1: Transfer should succeed")
	err = scenario2.Run()
	assert.NoError(t, err, "Scenario2: Transfer should succeed")

	allActions := append(scenario1.actions, scenario2.actions...)

	// Assert balances
	for _, user := range []BankOperator{userA, userB, userC} {
		latestBalance, err := user.GetAccountBalance()
		assert.NoError(t, err, "Should be able to get account balance for user %s", user.GetAccountId())

		expectedBalance := user.GetExpectedBalance(allActions)
		assert.Equal(t, expectedBalance, latestBalance,
			"Account %s should have correct balance", user.GetAccountId())

		log.Printf("✓ Account %s has correct balance: %.2f", user.GetAccountId(), latestBalance)
	}

	log.Println("🎉 All assertions passed! Test scenario completed successfully.")
}

func TestTransferMoneyFromAtoBtoCConcurrently(t *testing.T) {
	for i := 0; i < 5; i++ {
		t.Run(fmt.Sprintf("Run #%d", i+1), func(t *testing.T) {
			log.Println("Starting banking API test scenario...")

			userA := NewUser(100.0, "UserA")
			userB := NewUser(100.0, "UserB")
			userC := NewUser(100.0, "UserC")

			scenario1 := &TransferScenario{
				Name:     "Transfer from A to B, amount=30",
				FromUser: userA,
				ToUser:   userB,
				Amount:   30,
			}
			scenario2 := &TransferScenario{
				Name:     "Transfer from B to C, amount=20.0",
				FromUser: userB,
				ToUser:   userC,
				Amount:   20.0,
			}

			var startGw sync.WaitGroup
			startGw.Add(1)
			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				defer wg.Done()
				startGw.Wait()
				err := scenario1.Run()
				assert.NoError(t, err, "Scenario1: Transfer should succeed")
			}()
			go func() {
				defer wg.Done()
				startGw.Wait()
				err := scenario2.Run()
				assert.NoError(t, err, "Scenario2: Transfer should succeed")
			}()

			startGw.Done()
			wg.Wait()

			allActions := append(scenario1.actions, scenario2.actions...)
			for _, a := range allActions {
				log.Println(a.String())
				log.Println(a)
			}

			// Assert balances
			for _, user := range []BankOperator{userA, userB, userC} {
				latestBalance, err := user.GetAccountBalance()
				assert.NoError(t, err, "Should be able to get account balance for user %s", user.GetAccountId())
				log.Printf("User %s has balance: %.2f", user.GetName(), latestBalance)

				expectedBalance := user.GetExpectedBalance(allActions)
				assert.Equal(t, expectedBalance, latestBalance,
					"Account %s should have correct balance", user.GetName())
			}

			log.Println("🎉 All assertions passed! Test scenario completed successfully.")
		})
	}
}
