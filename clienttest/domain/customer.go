package domain

import (
	"errors"
	"log"
)

type Customer struct {
	initialBalance float64
	operator       BankOperator

	balanceChanges []balanceChange
	errors         []error
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

func (c *Customer) TransferMoney(toCustomer *Customer) error {
	// view account balance
	balance, err := c.operator.GetAccountBalance()
	if err != nil {
		// record error
	}

	transferMoney := balance * 0.1 // transfer 10% of balance
	err = c.operator.TransferTo(toCustomer.operator, transferMoney)
	if err != nil {
		// record error
		c.errors = append(c.errors, errors.New("transfer failed: "+err.Error()))
	} else {
		// track balance changes
		c.balanceChanges = append(c.balanceChanges, balanceChange{
			change: -transferMoney, // negative for withdrawal
		})
		toCustomer.onReceiving(transferMoney)
		log.Println("Transfer successful: ", c.operator.GetName(), " transferred", transferMoney, "to", toCustomer.operator.GetName())
	}

	return nil
}

func (c *Customer) onReceiving(amount float64) error {
	if amount <= 0 {
		return errors.New("invalid amount received, must be positive")
	}
	// record the amount received
	c.balanceChanges = append(c.balanceChanges, balanceChange{
		change: amount, // positive for deposit
	})
	return nil
}

func (c *Customer) VerifyBalance() (bool, error) {
	// compare the expected balance with the actual balance
	actualBalance, err := c.operator.GetAccountBalance()
	if err != nil {
		return false, err // error occurred, cannot verify balance
	}
	// calculate expected balance based on recorded changes
	expectedBalance := c.initialBalance
	for _, change := range c.balanceChanges {
		expectedBalance += change.change
	}
	ok := actualBalance == expectedBalance
	if ok {
		log.Printf("[%s] Balance = %.2f. Verified ok", c.operator.GetName(), actualBalance)
	} else {
		log.Printf("[%s] Verification failed. Balance = %.2f, expected = %.2f", c.operator.GetName(), actualBalance, expectedBalance)
		c.errors = append(c.errors, errors.New("balance verification failed"))
	}
	return actualBalance == expectedBalance, nil
}

func (c *Customer) Errors() []error {
	return c.errors
}
