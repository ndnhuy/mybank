package domain

import ()

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
	} else {
		// track balance changes
		c.balanceChanges = append(c.balanceChanges, balanceChange{
			change: -transferMoney, // negative for withdrawal
		})
	}

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
	return actualBalance == expectedBalance, nil
}
