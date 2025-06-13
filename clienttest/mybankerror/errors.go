package mybankerror

import "errors"

var (
	AccountAlreadyCreatedError = errors.New("account already created for this user")
	ErrInsufficientBalance = errors.New("insufficient balance for the operation")
)
