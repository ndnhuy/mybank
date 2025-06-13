package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFoo(t *testing.T) {
	customerA, _ := NewCustomer("customer A")
	customerB, _ := NewCustomer("customer B")
	customerA.TransferMoney(customerB)
	ok, err := customerA.VerifyBalance()
	assert.NoError(t, err, "Customer A should have a valid balance after transfer")
	assert.True(t, ok, "Customer A should have a valid balance after transfer")
}
