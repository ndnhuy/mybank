package domain

import (
	"errors"
	"log"
)

type Customer struct {
	initialBalance float64
	operator       BankOperator

	events  []Event
	actions []Action
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

// ExecuteAction executes an action and stores it in the action history
func (c *Customer) ExecuteAction(action Action) error {
	err := action.Execute()
	if err != nil {
		return err
	}
	c.actions = append(c.actions, action)
	return nil
}

// TransferMoneyAsAction transfers money using the action pattern
func (c *Customer) TransferMoneyAsAction(toCustomer *Customer) error {
	action := NewTransferAction(c, toCustomer)
	return c.ExecuteAction(action)
}

// ViewBalance captures the current balance in an action
func (c *Customer) ViewBalance() (float64, error) {
	action := NewViewBalanceAction(c)
	err := c.ExecuteAction(action)
	if err != nil {
		return 0, err
	}
	return action.GetSnapshotBalance(), nil
}

func (c *Customer) onReceiving(amount float64) error {
	if amount <= 0 {
		return errors.New("invalid amount received, must be positive")
	}
	// record the amount received
	event := NewBalanceWasChanged(amount)
	c.events = append(c.events, event)
	return nil
}

func (c *Customer) VerifyBalance() (bool, error) {
	// compare the expected balance with the actual balance
	actualBalance, err := c.operator.GetAccountBalance()
	if err != nil {
		return false, err
	}
	// calculate expected balance based on recorded events
	expectedBalance := c.initialBalance
	for _, event := range c.events {
		if balanceEvent, ok := event.(*BalanceWasChanged); ok {
			expectedBalance += balanceEvent.GetChange()
		}
	}
	ok := actualBalance == expectedBalance
	if ok {
		log.Printf("[%s] Balance = %.2f. Verified ok", c.operator.GetName(), actualBalance)
	} else {
		log.Printf("[%s] Verification failed. Balance = %.2f, expected = %.2f", c.operator.GetName(), actualBalance, expectedBalance)
	}
	return actualBalance == expectedBalance, nil
}
