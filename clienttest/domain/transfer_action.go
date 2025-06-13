package domain

import (
	"time"
)

// TransferAction represents a transfer to another customer
type TransferAction struct {
	customer   *Customer
	toCustomer *Customer
	amount     float64
	timestamp  time.Time
}

func NewTransferAction(customer *Customer, toCustomer *Customer) *TransferAction {
	return &TransferAction{
		customer:   customer,
		toCustomer: toCustomer,
		timestamp:  time.Now(),
	}
}

func (t *TransferAction) Execute() error {
	// Get current balance
	balance, err := t.customer.operator.GetAccountBalance()
	if err != nil {
		return err
	}

	// Calculate transfer amount (10% of balance)
	transferMoney := balance * 0.1
	t.amount = transferMoney

	// Execute the transfer
	err = t.customer.operator.TransferTo(t.toCustomer.operator, transferMoney)
	if err != nil {
		return err
	}

	// Record balance changes (keeping existing pattern)
	event := NewBalanceWasChanged(-transferMoney)
	t.customer.events = append(t.customer.events, event)
	t.toCustomer.onReceiving(transferMoney)

	return nil
}

func (t *TransferAction) GetType() string {
	return "transfer"
}

func (t *TransferAction) GetAmount() float64 {
	return t.amount
}

func (t *TransferAction) GetToCustomer() *Customer {
	return t.toCustomer
}
