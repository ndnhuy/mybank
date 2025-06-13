package domain

import (
	"testing"
)

func TestCustomerActionBasedWorkflow(t *testing.T) {
	// Create two customers
	customerA := NewCustomer()
	customerB := NewCustomer()

	// Customer A workflow: create account -> view balance -> transfer money

	// 1. Create account
	err := customerA.CreateAccount()
	if err != nil {
		t.Errorf("Failed to create account for customer A: %v", err)
	}

	// 2. View balance
	balance, err := customerA.ViewBalance()
	if err != nil {
		t.Errorf("Failed to view balance for customer A: %v", err)
	}
	if balance != 100.0 {
		t.Errorf("Expected initial balance 100.0, got %.2f", balance)
	}

	// 3. Transfer money (10% of balance = 10.0)
	transferAmount := balance * 0.1
	err = customerA.TransferMoney(customerB, transferAmount)
	if err != nil {
		t.Errorf("Failed to transfer money: %v", err)
	}

	// 4. Verify balances
	verifyResultA := customerA.VerifyBalance()
	if !verifyResultA.Ok {
		t.Errorf("Customer A balance verification failed, expected: %.2f, actual: %.2f",
			verifyResultA.Expected, verifyResultA.Actual)
	}

	verifyResultB := customerB.VerifyBalance()
	if !verifyResultB.Ok {
		t.Errorf("Customer B balance verification failed, expected: %.2f, actual: %.2f",
			verifyResultB.Expected, verifyResultB.Actual)
	}

	// Check action history
	actionsA := customerA.GetExecutedActions()
	actionsB := customerB.GetExecutedActions()

	// Customer A should have: create_account, view_balance, transfer_out
	expectedActionsA := 3
	if len(actionsA) != expectedActionsA {
		t.Errorf("Customer A should have %d actions, got %d", expectedActionsA, len(actionsA))
	}

	// Customer B should have: create_account, transfer_in
	expectedActionsB := 2
	if len(actionsB) != expectedActionsB {
		t.Errorf("Customer B should have %d actions, got %d", expectedActionsB, len(actionsB))
	}

	// Verify action types for Customer A
	if len(actionsA) >= 3 {
		if actionsA[0].ActionType != "create_account" {
			t.Errorf("Expected first action to be 'create_account', got '%s'", actionsA[0].ActionType)
		}
		if actionsA[1].ActionType != "view_balance" {
			t.Errorf("Expected second action to be 'view_balance', got '%s'", actionsA[1].ActionType)
		}
		if actionsA[2].ActionType != "transfer_out" {
			t.Errorf("Expected third action to be 'transfer_out', got '%s'", actionsA[2].ActionType)
		}
	}

	// Print action histories for debugging
	t.Logf("Customer A Action History:\n%s", customerA.GetActionHistory())
	t.Logf("Customer B Action History:\n%s", customerB.GetActionHistory())
}

func TestCustomerSequentialActions(t *testing.T) {
	// Create customer
	customer := NewCustomer()

	// Test sequence: create -> view -> view -> transfer (should fail due to no recipient)

	// 1. Create account
	err := customer.CreateAccount()
	if err != nil {
		t.Errorf("Failed to create account: %v", err)
	}

	// 2. View balance multiple times
	balance1, err := customer.ViewBalance()
	if err != nil {
		t.Errorf("Failed to view balance (1st time): %v", err)
	}

	balance2, err := customer.ViewBalance()
	if err != nil {
		t.Errorf("Failed to view balance (2nd time): %v", err)
	}

	if balance1 != balance2 {
		t.Errorf("Balance should be consistent: %.2f != %.2f", balance1, balance2)
	}

	// 3. Check that actions are tracked correctly
	actions := customer.GetExecutedActions()
	expectedActionCount := 3 // create_account + 2x view_balance
	if len(actions) != expectedActionCount {
		t.Errorf("Expected %d actions, got %d", expectedActionCount, len(actions))
	}

	// 4. Verify all actions succeeded
	for i, action := range actions {
		if !action.Success {
			t.Errorf("Action %d should have succeeded: %s - %s", i, action.ActionType, action.ErrorMessage)
		}
	}

	t.Logf("Sequential Actions Test - Action History:\n%s", customer.GetActionHistory())
}
