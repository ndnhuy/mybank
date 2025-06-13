package domain

import "time"

// BalanceWasChanged represents a balance change event
type BalanceWasChanged struct {
	change    float64 // positive for deposit, negative for withdrawal
	timestamp time.Time
}

func NewBalanceWasChanged(change float64) *BalanceWasChanged {
	return &BalanceWasChanged{
		change:    change,
		timestamp: time.Now(),
	}
}

func (b *BalanceWasChanged) GetTimestamp() time.Time {
	return b.timestamp
}

func (b *BalanceWasChanged) GetType() string {
	return "balance_was_changed"
}

func (b *BalanceWasChanged) GetChange() float64 {
	return b.change
}
