package domain

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertBalance(t *testing.T, customer *Customer) {
	err := customer.VerifyBalance()
	require.NoError(t, err, fmt.Sprintf("Balance verification failed for customer: %s", customer.operator.GetAccountId()))
}

func TestTransfer(t *testing.T) {
	customerA, _ := NewCustomer("customer A")
	customerB, _ := NewCustomer("customer B")
	customerA.TransferMoney(customerB, 100.00)
	assertBalance(t, customerA)
}

func TestTransferSequentially(t *testing.T) {
	customerA, _ := NewCustomerWithAmount("customer A", 100.00)
	customerB, _ := NewCustomerWithAmount("customer B", 100.00)
	customerC, _ := NewCustomerWithAmount("customer C", 100.00)

	// Transfer from A to B
	err := customerA.TransferMoney(customerB, 10)
	assert.NoError(t, err, "Transfer from A to B should succeed")

	// Verify balances after first transfer
	assertBalance(t, customerA)
	assertBalance(t, customerB)

	// Transfer from B to C
	err = customerB.TransferMoney(customerC, 10.00)
	assert.NoError(t, err, "Transfer from B to C should succeed")

	// Verify balances after second transfer
	assertBalance(t, customerB)
	assertBalance(t, customerC)
}

func TestTransferConcurrently(t *testing.T) {
	for i := 0; i < 5; i++ {
		t.Run(fmt.Sprintf("Run #%d", i+1), func(t *testing.T) {
			customerA, _ := NewCustomer("customer A")
			customerB, _ := NewCustomer("customer B")
			customerC, _ := NewCustomer("customer C")

			var startGw sync.WaitGroup
			startGw.Add(1)
			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				defer wg.Done()
				startGw.Wait()
				err := customerA.TransferMoney(customerB, 100.00)
				assert.NoError(t, err, "Transfer from A to B should succeed")
			}()
			go func() {
				defer wg.Done()
				startGw.Wait()
				err := customerB.TransferMoney(customerC, 100.00)
				assert.NoError(t, err, "Transfer from B to C should succeed")
			}()

			startGw.Done()
			wg.Wait()

			assertBalance(t, customerA)
			assertBalance(t, customerB)
			assertBalance(t, customerC)
		})
	}
}
