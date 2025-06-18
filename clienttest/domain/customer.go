package domain

import "fmt"

type Customer struct {
	initialBalance float64
	operator       BankOperator

	balanceChanges []balanceChange
}

type balanceChange struct {
	change float64 // positive for deposit, negative for withdrawal
}

func NewCustomer(alias string) (*Customer, error) {
	operator := NewBankOperatorImpl(100.00, alias)
	_, err := operator.CreateAccount()
	if err != nil {
		return nil, err
	}

	return &Customer{
		operator:       operator,
		initialBalance: operator.InitialBalance,
	}, nil
}

func NewCustomerWithAmount(alias string, initialAmount float64) (*Customer, error) {
	operator := NewBankOperatorImpl(initialAmount, alias)
	_, err := operator.CreateAccount()
	if err != nil {
		return nil, err
	}

	return &Customer{
		operator:       operator,
		initialBalance: operator.InitialBalance,
	}, nil
}

func (c *Customer) TransferMoney(toCustomer *Customer, amount float64) error {
	transferMoney := amount
	err := c.operator.TransferTo(toCustomer.operator, transferMoney)
	if err != nil {
		return err
	} else {
		// track balance changes
		c.balanceChanges = append(c.balanceChanges, balanceChange{
			change: -transferMoney, // negative for withdrawal
		})
		toCustomer.onReceiveMoney(transferMoney) // notify recipient
	}

	return nil
}

func (c *Customer) RecordTransfer(toCustomer *Customer, amount float64) error {
	// This method is for internal tracking, not for actual transfers
	if toCustomer == nil || amount <= 0 {
		return fmt.Errorf("invalid transfer parameters")
	}

	c.balanceChanges = append(c.balanceChanges, balanceChange{
		change: -amount, // negative for withdrawal
	})
	toCustomer.onReceiveMoney(amount) // notify recipient

	return nil
}

func (c *Customer) onReceiveMoney(amount float64) {
	// track balance changes
	c.balanceChanges = append(c.balanceChanges, balanceChange{
		change: amount, // positive for deposit
	})
}

func (c *Customer) VerifyBalance() error {
	actualBalance, err := c.operator.GetAccountBalance()
	if err != nil {
		return err // error occurred, cannot verify balance
	}
	// calculate expected balance based on recorded changes
	expectedBalance := c.initialBalance
	for _, change := range c.balanceChanges {
		expectedBalance += change.change
	}
	if actualBalance != expectedBalance {
		return fmt.Errorf("[%v] balance mismatch: expected %.2f, got %.2f", c.operator.GetName(), expectedBalance, actualBalance)
	} else {
		fmt.Printf("[%v] balance verified: %.2f\n", c.operator.GetName(), actualBalance)
	}
	return nil
}

// GetAccountID returns the customer's account ID for load testing
func (c *Customer) GetAccountID() string {
	return c.operator.GetAccountId()
}

// GetCurrentBalance returns the current balance from the bank
func (c *Customer) GetCurrentBalance() (float64, error) {
	return c.operator.GetAccountBalance()
}

// GetName returns the customer's name/alias
func (c *Customer) GetName() string {
	return c.operator.GetName()
}
