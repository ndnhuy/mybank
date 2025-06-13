package domain

import "time"

// ViewBalanceAction captures a balance snapshot
type ViewBalanceAction struct {
	customer        *Customer
	snapshotBalance float64
	timestamp       time.Time
}

func NewViewBalanceAction(customer *Customer) *ViewBalanceAction {
	return &ViewBalanceAction{
		customer:  customer,
		timestamp: time.Now(),
	}
}

func (v *ViewBalanceAction) Execute() error {
	balance, err := v.customer.operator.GetAccountBalance()
	if err != nil {
		return err
	}
	v.snapshotBalance = balance
	return nil
}

func (v *ViewBalanceAction) GetType() string {
	return "view_balance"
}

func (v *ViewBalanceAction) GetSnapshotBalance() float64 {
	return v.snapshotBalance
}
