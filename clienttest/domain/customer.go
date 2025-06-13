package domain

import (
	"fmt"
	"log"

	mybankerror "com.ndnhuy.mybank/mybankerror"
)

type Customer struct {
	operator        BankOperator     // interface for user operations
	executedActions []CustomerAction // track all actions performed by this customer
	initialBalance  float64          // starting balance for calculations
}

// CustomerAction represents an action that affects the customer's balance
type CustomerAction struct {
	ActionType   string  // "create_account", "view_balance", "transfer_out", "transfer_in"
	Amount       float64 // amount involved (0 for non-monetary actions)
	AccountId    string  // account affected
	AccountName  string  // account name
	Success      bool    // whether action was successful
	ErrorMessage string  // error message if action failed
}

type VerifyResult struct {
	Ok       bool    // indicates if the verification passed
	Expected float64 // expected balance after the operation
	Actual   float64 // actual balance after the operation
}

// NewCustomer creates a new Customer instance
func NewCustomer() *Customer {
	return &Customer{
		operator:        NewUser(100.0, "Customer"), // Initial balance of 100.0, can be adjusted
		initialBalance:  100.0,
		executedActions: make([]CustomerAction, 0),
	}
}

// CreateAccount executes a create account action
func (c *Customer) CreateAccount() error {
	action := CustomerAction{
		ActionType:  "create_account",
		AccountId:   "", // Will be set after creation
		AccountName: c.operator.GetName(),
	}

	_, err := c.operator.CreateAccount()
	if err != nil && err != mybankerror.AccountAlreadyCreatedError {
		action.Success = false
		action.ErrorMessage = err.Error()
		c.executedActions = append(c.executedActions, action)
		log.Printf("[%s] Failed to create account: %v", c.operator.GetName(), err)
		return err
	}

	action.Success = true
	action.AccountId = c.operator.GetAccountId()
	c.executedActions = append(c.executedActions, action)
	log.Printf("[%s] Account created successfully", c.operator.GetName())
	return nil
}

// ViewBalance executes a view balance action
func (c *Customer) ViewBalance() (float64, error) {
	action := CustomerAction{
		ActionType:  "view_balance",
		AccountId:   c.operator.GetAccountId(),
		AccountName: c.operator.GetName(),
	}

	balance, err := c.operator.GetAccountBalance()
	if err != nil {
		action.Success = false
		action.ErrorMessage = err.Error()
		c.executedActions = append(c.executedActions, action)
		log.Printf("[%s] Failed to get balance: %v", c.operator.GetName(), err)
		return 0, err
	}

	action.Success = true
	action.Amount = balance
	c.executedActions = append(c.executedActions, action)
	log.Printf("[%s] Current balance: %.2f", c.operator.GetName(), balance)
	return balance, nil
}

// hasAccount checks if the customer already has an account created
func (c *Customer) hasAccount() bool {
	for _, action := range c.executedActions {
		if action.ActionType == "create_account" && action.Success {
			return true
		}
	}
	return false
}

// ensureAccount creates an account only if it doesn't exist
func (c *Customer) ensureAccount() error {
	if !c.hasAccount() {
		return c.CreateAccount()
	}
	return nil
}

// TransferMoney executes a transfer money action
func (c *Customer) TransferMoney(toCustomer *Customer, amount float64) error {
	// Ensure both accounts exist
	err := c.ensureAccount()
	if err != nil {
		return err
	}

	err = toCustomer.ensureAccount()
	if err != nil {
		return err
	}

	// Record the transfer action for sender
	senderAction := CustomerAction{
		ActionType:  "transfer_out",
		Amount:      -amount, // negative for outgoing
		AccountId:   c.operator.GetAccountId(),
		AccountName: c.operator.GetName(),
	}

	// Record the transfer action for receiver
	receiverAction := CustomerAction{
		ActionType:  "transfer_in",
		Amount:      amount, // positive for incoming
		AccountId:   toCustomer.operator.GetAccountId(),
		AccountName: toCustomer.operator.GetName(),
	}

	// Check sufficient balance
	balance, err := c.operator.GetAccountBalance()
	if err != nil {
		senderAction.Success = false
		senderAction.ErrorMessage = err.Error()
		c.executedActions = append(c.executedActions, senderAction)
		return err
	}

	if balance < amount {
		senderAction.Success = false
		senderAction.ErrorMessage = "insufficient balance"
		c.executedActions = append(c.executedActions, senderAction)
		log.Printf("[%s] Insufficient balance: %.2f < %.2f", c.operator.GetName(), balance, amount)
		return mybankerror.ErrInsufficientBalance
	}

	// Perform transfer
	err = c.operator.TransferTo(toCustomer.operator, amount)
	if err != nil {
		senderAction.Success = false
		senderAction.ErrorMessage = err.Error()
		receiverAction.Success = false
		receiverAction.ErrorMessage = err.Error()
		c.executedActions = append(c.executedActions, senderAction)
		toCustomer.executedActions = append(toCustomer.executedActions, receiverAction)
		log.Printf("[%s] Transfer failed: %v", c.operator.GetName(), err)
		return err
	}

	// Mark both actions as successful
	senderAction.Success = true
	receiverAction.Success = true
	c.executedActions = append(c.executedActions, senderAction)
	toCustomer.executedActions = append(toCustomer.executedActions, receiverAction)

	log.Printf("[%s] Successfully transferred %.2f to [%s]",
		c.operator.GetName(), amount, toCustomer.operator.GetName())
	return nil
}

func (c *Customer) Transfer(toCustomer *Customer) error {
	// Create accounts if they don't exist
	err := c.ensureAccount()
	if err != nil {
		return err
	}

	err = toCustomer.ensureAccount()
	if err != nil {
		return err
	}

	// get current balance
	balance, err := c.ViewBalance()
	if err != nil {
		return err
	}

	if balance <= 0 {
		return mybankerror.ErrInsufficientBalance
	}

	// transfer 10% of balance to another customer
	transferAmount := balance * 0.1
	return c.TransferMoney(toCustomer, transferAmount)
}

// GetExecutedActions returns a copy of all executed actions
func (c *Customer) GetExecutedActions() []CustomerAction {
	actionsCopy := make([]CustomerAction, len(c.executedActions))
	copy(actionsCopy, c.executedActions)
	return actionsCopy
}

// GetActionHistory returns a summary of all executed actions
func (c *Customer) GetActionHistory() string {
	if len(c.executedActions) == 0 {
		return "No actions executed"
	}

	history := "Action History:\n"
	for i, action := range c.executedActions {
		status := "✓"
		if !action.Success {
			status = "✗"
		}
		history += fmt.Sprintf("%d. %s %s %s", i+1, status, action.ActionType, action.AccountName)
		if action.Amount != 0 {
			history += fmt.Sprintf(" (%.2f)", action.Amount)
		}
		if !action.Success {
			history += fmt.Sprintf(" - Error: %s", action.ErrorMessage)
		}
		history += "\n"
	}
	return history
}

func (c *Customer) VerifyBalance() *VerifyResult {
	// get current balance
	actual, err := c.operator.GetAccountBalance()
	if err != nil {
		log.Printf("Error getting account balance for verification: %v", err)
		return &VerifyResult{
			Ok:       false,
			Expected: 0.0,
			Actual:   actual,
		}
	}

	// calculate expected balance based on executed actions
	expected := c.calculateExpectedBalance()

	// compare actual balance with expected balance
	ok := actual == expected

	return &VerifyResult{
		Ok:       ok,
		Expected: expected,
		Actual:   actual,
	}
}

// calculateExpectedBalance calculates the expected balance based on executed actions
func (c *Customer) calculateExpectedBalance() float64 {
	balance := c.initialBalance

	for _, action := range c.executedActions {
		if action.Success && (action.ActionType == "transfer_out" || action.ActionType == "transfer_in") {
			balance += action.Amount // Amount is already signed (negative for out, positive for in)
		}
	}

	return balance
}
