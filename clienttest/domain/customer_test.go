package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper function to verify customer balance
func verifyCustomerBalance(t *testing.T, customer *Customer, description string) {
	ok, err := customer.VerifyBalance()
	assert.NoError(t, err, description)
	assert.True(t, ok, description)
}

func TestTransferFromAtoB(t *testing.T) {
	customerA, _ := NewCustomer("customer A")
	customerB, _ := NewCustomer("customer B")
	customerA.TransferMoneyAsAction(customerB)

	verifyCustomerBalance(t, customerA, "Customer A should have a valid balance after transfer")
	verifyCustomerBalance(t, customerB, "Customer B should have a valid balance after receiving transfer")
}

func TestTransferFromAtoBtoC(t *testing.T) {
	customerA, _ := NewCustomer("customer A")
	customerB, _ := NewCustomer("customer B")
	customerC, _ := NewCustomer("customer C")

	// Transfer from A to B
	err := customerA.TransferMoneyAsAction(customerB)
	assert.NoError(t, err, "Transfer from A to B should succeed")

	// Verify balances after first transfer
	verifyCustomerBalance(t, customerA, "Customer A should have a valid balance after transfer to B")
	verifyCustomerBalance(t, customerB, "Customer B should have a valid balance after receiving transfer from A")

	// Transfer from B to C
	err = customerB.TransferMoneyAsAction(customerC)
	assert.NoError(t, err, "Transfer from B to C should succeed")

	// Verify balances after second transfer
	verifyCustomerBalance(t, customerB, "Customer B should have a valid balance after transfer to C")
	verifyCustomerBalance(t, customerC, "Customer C should have a valid balance after receiving transfer from B")
}

func TestSequentialActions(t *testing.T) {
	customerA, _ := NewCustomer("customer A")
	customerB, _ := NewCustomer("customer B")

	// 1. View current balance
	initialBalance, err := customerA.ViewBalance()
	assert.NoError(t, err, "Initial balance view should succeed")
	assert.Equal(t, 100.0, initialBalance, "Initial balance should be 100")

	// 2. Transfer money to other customer
	err = customerA.TransferMoneyAsAction(customerB)
	assert.NoError(t, err, "Transfer action should succeed")

	// 3. View current balance and check if it's correct
	finalBalance, err := customerA.ViewBalance()
	assert.NoError(t, err, "Final balance view should succeed")
	expectedBalance := initialBalance * 0.9 // 90% remaining after 10% transfer
	assert.Equal(t, expectedBalance, finalBalance, "Final balance should be 90% of initial")

	// Verify customer state
	verifyCustomerBalance(t, customerA, "Customer A should have valid balance after sequential actions")
	verifyCustomerBalance(t, customerB, "Customer B should have valid balance after receiving transfer")
}
