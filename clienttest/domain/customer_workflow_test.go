package domain

import (
	"testing"
)

func TestCustomerCompleteWorkflow(t *testing.T) {
	// Test complete customer workflow:
	// Customer A: create -> view -> transfer -> verify
	// Customer B: receive -> verify

	customerA := NewCustomer()
	customerB := NewCustomer()

	// Step 1: Customer A creates account
	err := customerA.CreateAccount()
	if err != nil {
		t.Fatalf("Customer A failed to create account: %v", err)
	}

	// Step 2: Customer A views balance (should be 100.0)
	balance, err := customerA.ViewBalance()
	if err != nil {
		t.Fatalf("Customer A failed to view balance: %v", err)
	}
	if balance != 100.0 {
		t.Errorf("Expected Customer A initial balance 100.0, got %.2f", balance)
	}

	// Step 3: Customer A transfers 25.0 to Customer B
	err = customerA.TransferMoney(customerB, 25.0)
	if err != nil {
		t.Fatalf("Customer A failed to transfer money: %v", err)
	}

	// Step 4: Both customers view their balances
	balanceA, err := customerA.ViewBalance()
	if err != nil {
		t.Fatalf("Customer A failed to view balance after transfer: %v", err)
	}

	balanceB, err := customerB.ViewBalance()
	if err != nil {
		t.Fatalf("Customer B failed to view balance after transfer: %v", err)
	}

	// Verify expected balances
	expectedBalanceA := 75.0  // 100 - 25
	expectedBalanceB := 125.0 // 100 + 25

	if balanceA != expectedBalanceA {
		t.Errorf("Customer A: expected balance %.2f, got %.2f", expectedBalanceA, balanceA)
	}

	if balanceB != expectedBalanceB {
		t.Errorf("Customer B: expected balance %.2f, got %.2f", expectedBalanceB, balanceB)
	}

	// Step 5: Verify balance calculations match
	resultA := customerA.VerifyBalance()
	if !resultA.Ok {
		t.Errorf("Customer A balance verification failed: expected %.2f, actual %.2f",
			resultA.Expected, resultA.Actual)
	}

	resultB := customerB.VerifyBalance()
	if !resultB.Ok {
		t.Errorf("Customer B balance verification failed: expected %.2f, actual %.2f",
			resultB.Expected, resultB.Actual)
	}

	// Step 6: Verify action histories
	actionsA := customerA.GetExecutedActions()
	actionsB := customerB.GetExecutedActions()

	// Customer A should have: create_account, view_balance, transfer_out, view_balance
	expectedActionsCountA := 4
	if len(actionsA) != expectedActionsCountA {
		t.Errorf("Customer A should have %d actions, got %d", expectedActionsCountA, len(actionsA))
		t.Logf("Customer A actions: %v", actionsA)
	}

	// Customer B should have: create_account, transfer_in, view_balance
	expectedActionsCountB := 3
	if len(actionsB) != expectedActionsCountB {
		t.Errorf("Customer B should have %d actions, got %d", expectedActionsCountB, len(actionsB))
		t.Logf("Customer B actions: %v", actionsB)
	}

	// Verify specific action types for Customer A
	if len(actionsA) >= 4 {
		expectedActionsA := []string{"create_account", "view_balance", "transfer_out", "view_balance"}
		for i, expectedAction := range expectedActionsA {
			if actionsA[i].ActionType != expectedAction {
				t.Errorf("Customer A action %d: expected '%s', got '%s'",
					i, expectedAction, actionsA[i].ActionType)
			}
		}
	}

	// Verify specific action types for Customer B
	if len(actionsB) >= 3 {
		expectedActionsB := []string{"create_account", "transfer_in", "view_balance"}
		for i, expectedAction := range expectedActionsB {
			if actionsB[i].ActionType != expectedAction {
				t.Errorf("Customer B action %d: expected '%s', got '%s'",
					i, expectedAction, actionsB[i].ActionType)
			}
		}
	}

	// Print action histories for verification
	t.Logf("=== Customer A Action History ===\n%s", customerA.GetActionHistory())
	t.Logf("=== Customer B Action History ===\n%s", customerB.GetActionHistory())

	// Additional verification: check amounts in transfer actions
	for _, action := range actionsA {
		if action.ActionType == "transfer_out" {
			expectedAmount := -25.0 // negative for outgoing
			if action.Amount != expectedAmount {
				t.Errorf("Customer A transfer_out amount: expected %.2f, got %.2f",
					expectedAmount, action.Amount)
			}
		}
	}

	for _, action := range actionsB {
		if action.ActionType == "transfer_in" {
			expectedAmount := 25.0 // positive for incoming
			if action.Amount != expectedAmount {
				t.Errorf("Customer B transfer_in amount: expected %.2f, got %.2f",
					expectedAmount, action.Amount)
			}
		}
	}
}

func TestCustomerActionBasedDesignPrinciples(t *testing.T) {
	// This test demonstrates the key principles of our action-based design:
	// 1. Each customer maintains their own action history
	// 2. Actions are immutable once executed
	// 3. Balance verification works based on action history
	// 4. Actions can be inspected and debugged

	customer := NewCustomer()

	// Initially, no actions
	actions := customer.GetExecutedActions()
	if len(actions) != 0 {
		t.Errorf("New customer should have no actions, got %d", len(actions))
	}

	// Action 1: Create account
	err := customer.CreateAccount()
	if err != nil {
		t.Fatalf("Failed to create account: %v", err)
	}

	actions = customer.GetExecutedActions()
	if len(actions) != 1 {
		t.Errorf("After create account, should have 1 action, got %d", len(actions))
	}

	if actions[0].ActionType != "create_account" {
		t.Errorf("First action should be create_account, got %s", actions[0].ActionType)
	}

	// Action 2: View balance
	balance, err := customer.ViewBalance()
	if err != nil {
		t.Fatalf("Failed to view balance: %v", err)
	}

	actions = customer.GetExecutedActions()
	if len(actions) != 2 {
		t.Errorf("After view balance, should have 2 actions, got %d", len(actions))
	}

	if actions[1].ActionType != "view_balance" {
		t.Errorf("Second action should be view_balance, got %s", actions[1].ActionType)
	}

	if actions[1].Amount != balance {
		t.Errorf("View balance action should record balance %.2f, got %.2f",
			balance, actions[1].Amount)
	}

	// Verify that all actions succeeded
	for i, action := range actions {
		if !action.Success {
			t.Errorf("Action %d should have succeeded: %s - %s",
				i, action.ActionType, action.ErrorMessage)
		}
	}

	// Verify balance calculation
	result := customer.VerifyBalance()
	if !result.Ok {
		t.Errorf("Balance verification should pass: expected %.2f, actual %.2f",
			result.Expected, result.Actual)
	}

	t.Logf("Action-based design test completed successfully!")
	t.Logf("Customer performed %d actions:", len(actions))
	t.Logf("%s", customer.GetActionHistory())
}
