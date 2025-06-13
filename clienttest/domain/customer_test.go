package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransferFromAtoB(t *testing.T) {
	customerA, _ := NewCustomer("customer A")
	customerB, _ := NewCustomer("customer B")
	customerA.TransferMoney(customerB)
	ok, err := customerA.VerifyBalance()
	assert.NoError(t, err, "Customer A should have a valid balance after transfer")
	assert.True(t, ok, "Customer A should have a valid balance after transfer")

	ok, err = customerB.VerifyBalance()
	assert.NoError(t, err, "Customer B should have a valid balance after receiving transfer")
	assert.True(t, ok, "Customer B should have a valid balance after receiving transfer")
}
