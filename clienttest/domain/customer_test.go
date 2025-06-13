package domain

import (
	"testing"
)

func TestFoo(t *testing.T) {
	customerA := NewCustomer()
	customerB := NewCustomer()
	customerA.Transfer(customerB)
	verifyResult := customerA.VerifyBalance()
	if !verifyResult.Ok {
		t.Errorf("Transfer failed, expected: %v, actual: %v", verifyResult.Expected, verifyResult.Actual)
	}
}
