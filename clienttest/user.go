package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type User struct {
	InitialBalance float64
	CurrentBalance float64 // Current balance of the user
	AccountId      string
}

// NewUser creates a new User with the specified initial balance
func NewUser(initialBalance float64) *User {
	return &User{InitialBalance: initialBalance, CurrentBalance: initialBalance}
}

func (u *User) GetAccount(accountID string) (*AccountInfo, error) {
	resp, err := http.Get(fmt.Sprintf("%s/accounts/%s", baseURL, accountID))
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get account failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var account AccountInfo
	if err := json.Unmarshal(body, &account); err != nil {
		return nil, fmt.Errorf("failed to unmarshal account info: %w", err)
	}

	return &account, nil
}

func (u *User) GetAccountBalance() (float64, error) {
	account, err := u.GetAccount(u.AccountId)
	if err != nil {
		return 0, fmt.Errorf("failed to get account balance: %w", err)
	}
	return account.Balance, nil
}

func (u *User) CreateAccount() (*AccountInfo, error) {
	resp, err := http.Post(baseURL+"/accounts", "application/json", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("create account failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var account AccountInfo
	if err := json.Unmarshal(body, &account); err != nil {
		return nil, fmt.Errorf("failed to unmarshal account info: %w", err)
	}

	u.AccountId = account.ID

	return &account, nil
}

// CreateAccountWithInitialBalance creates an account for this user (currently API doesn't support custom initial balance)
func (u *User) CreateAccountWithInitialBalance() (*AccountInfo, error) {
	// Note: The current API creates accounts with a fixed 100.0 initial balance
	// This method is prepared for when the API supports custom initial balances
	return u.CreateAccount()
}

func (u *User) TransferTo(fromAccountID, toAccountID string, amount float64) error {
	transferReq := TransferRequest{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        amount,
	}

	reqBody, err := json.Marshal(transferReq)
	if err != nil {
		return fmt.Errorf("failed to marshal transfer request: %w", err)
	}

	resp, err := http.Post(baseURL+"/accounts/transfer", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to perform transfer: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("transfer failed with status: %d", resp.StatusCode)
	}

	return nil
}

func (u *User) ApplyActions(actions []action) *User {
	for _, act := range actions {
		if act.accountId == u.AccountId {
			u.CurrentBalance += act.balanceChange
		}
	}
	return u
}
