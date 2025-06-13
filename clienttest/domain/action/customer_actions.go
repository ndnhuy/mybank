package action

import (
	"fmt"
	"log"

	"com.ndnhuy.mybank/domain"
	mybankerror "com.ndnhuy.mybank/mybankerror"
)

// Forward declare interfaces and types that will be imported
// In the real implementation, these should be imported from domain package

// ActionResult represents the result of an action execution
type ActionResult struct {
	ActionType   string
	Success      bool
	Amount       float64
	ErrorMessage string
	AccountId    string
	AccountName  string
}

// CustomerActionBase provides common functionality for all customer actions
type CustomerActionBase struct {
	executed bool
	success  bool
	err      error
	result   ActionResult
}

func (a *CustomerActionBase) Success() bool {
	return a.executed && a.success
}

func (a *CustomerActionBase) Error() error {
	return a.err
}

func (a *CustomerActionBase) GetResult() ActionResult {
	return a.result
}

// CreateAccountAction represents the action of creating a bank account
type CreateAccountAction struct {
	operator domain.BankOperator
	success  bool
	err      error
	executed bool
}

func NewCreateAccountAction(operator domain.BankOperator) *CreateAccountAction {
	return &CreateAccountAction{
		operator: operator,
	}
}

func (a *CreateAccountAction) Run() error {
	if a.executed {
		return fmt.Errorf("action already executed")
	}

	log.Printf("[%s] Creating account...", a.operator.GetName())
	_, err := a.operator.CreateAccount()

	a.executed = true
	if err != nil && err != mybankerror.AccountAlreadyCreatedError {
		a.success = false
		a.err = err
		log.Printf("[%s] Failed to create account: %v", a.operator.GetName(), err)
		return err
	}

	a.success = true
	a.err = nil
	log.Printf("[%s] Account created successfully", a.operator.GetName())
	return nil
}

func (a *CreateAccountAction) Success() bool {
	return a.executed && a.success
}

func (a *CreateAccountAction) Error() error {
	return a.err
}

// ViewBalanceAction represents the action of viewing account balance
type ViewBalanceAction struct {
	operator domain.BankOperator
	balance  float64
	success  bool
	err      error
	executed bool
}

func NewViewBalanceAction(operator domain.BankOperator) *ViewBalanceAction {
	return &ViewBalanceAction{
		operator: operator,
	}
}

func (a *ViewBalanceAction) Run() error {
	if a.executed {
		return fmt.Errorf("action already executed")
	}

	log.Printf("[%s] Viewing balance...", a.operator.GetName())
	balance, err := a.operator.GetAccountBalance()

	a.executed = true
	if err != nil {
		a.success = false
		a.err = err
		log.Printf("[%s] Failed to get balance: %v", a.operator.GetName(), err)
		return err
	}

	a.success = true
	a.err = nil
	a.balance = balance
	log.Printf("[%s] Current balance: %.2f", a.operator.GetName(), balance)
	return nil
}

func (a *ViewBalanceAction) Success() bool {
	return a.executed && a.success
}

func (a *ViewBalanceAction) Error() error {
	return a.err
}

func (a *ViewBalanceAction) GetBalance() float64 {
	return a.balance
}

// TransferMoneyAction represents the action of transferring money to another customer
type TransferMoneyAction struct {
	fromOperator domain.BankOperator
	toOperator   domain.BankOperator
	amount       float64
	success      bool
	err          error
	executed     bool
}

func NewTransferMoneyAction(from, to domain.BankOperator, amount float64) *TransferMoneyAction {
	return &TransferMoneyAction{
		fromOperator: from,
		toOperator:   to,
		amount:       amount,
	}
}

func (a *TransferMoneyAction) Run() error {
	if a.executed {
		return fmt.Errorf("action already executed")
	}

	log.Printf("[%s] Transferring %.2f to [%s]...",
		a.fromOperator.GetName(), a.amount, a.toOperator.GetName())

	// Ensure both accounts exist
	_, err := a.fromOperator.CreateAccount()
	if err != nil && err != mybankerror.AccountAlreadyCreatedError {
		a.executed = true
		a.success = false
		a.err = err
		return err
	}

	_, err = a.toOperator.CreateAccount()
	if err != nil && err != mybankerror.AccountAlreadyCreatedError {
		a.executed = true
		a.success = false
		a.err = err
		return err
	}

	// Check sufficient balance
	balance, err := a.fromOperator.GetAccountBalance()
	if err != nil {
		a.executed = true
		a.success = false
		a.err = err
		return err
	}

	if balance < a.amount {
		a.executed = true
		a.success = false
		a.err = mybankerror.ErrInsufficientBalance
		log.Printf("[%s] Insufficient balance: %.2f < %.2f",
			a.fromOperator.GetName(), balance, a.amount)
		return a.err
	}

	// Perform transfer
	err = a.fromOperator.TransferTo(a.toOperator, a.amount)
	a.executed = true

	if err != nil {
		a.success = false
		a.err = err
		log.Printf("[%s] Transfer failed: %v", a.fromOperator.GetName(), err)
		return err
	}

	a.success = true
	a.err = nil
	log.Printf("[%s] Successfully transferred %.2f to [%s]",
		a.fromOperator.GetName(), a.amount, a.toOperator.GetName())
	return nil
}

func (a *TransferMoneyAction) Success() bool {
	return a.executed && a.success
}

func (a *TransferMoneyAction) Error() error {
	return a.err
}

func (a *TransferMoneyAction) GetAmount() float64 {
	return a.amount
}

func (a *TransferMoneyAction) GetFromOperator() domain.BankOperator {
	return a.fromOperator
}

func (a *TransferMoneyAction) GetToOperator() domain.BankOperator {
	return a.toOperator
}
