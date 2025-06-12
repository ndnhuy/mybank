package domain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	mybankerror "com.ndnhuy.mybank/mybankerror" // Adjust import path as needed
	"com.ndnhuy.mybank/utils"
)

type BankUserImpl struct {
	InitialBalance float64
	AccountId      string
	Name           string // Optional alias for the user

	mu sync.RWMutex
}

// NewUser creates a new User with the specified initial balance
func NewUser(initialBalance float64, name string) BankUser {
	return &BankUserImpl{InitialBalance: initialBalance, Name: name}
}

func (u *BankUserImpl) GetAccount(accountID string) (*AccountInfo, error) {
	resp, err := http.Get(fmt.Sprintf("%s/accounts/%s", utils.BASE_URL, accountID))
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

func (u *BankUserImpl) GetAccountBalance() (float64, error) {
	account, err := u.GetAccount(u.AccountId)
	if err != nil {
		return 0, fmt.Errorf("failed to get account balance: %w", err)
	}
	return account.Balance, nil
}

func (u *BankUserImpl) CreateAccount() (*AccountInfo, error) {
	// validate
	if u.InitialBalance <= 0 {
		return nil, fmt.Errorf("initial balance must be greater than zero")
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	if u.AccountId != "" {
		accId := u.AccountId

		acc, err := u.GetAccount(accId)
		if err != nil {
			return nil, fmt.Errorf("failed to get existing account: %w", err)
		}
		return acc, mybankerror.AccountAlreadyCreatedError
	}

	account, err := u.createAccountRequest()
	if err != nil {
		return nil, err
	}

	u.AccountId = account.ID

	log.Printf("[%v] Created account with ID: %s and initial balance: %.2f", u.Name, account.ID, u.InitialBalance)

	return account, nil
}

func (u *BankUserImpl) createAccountRequest() (*AccountInfo, error) {
	// Create account with initial balance
	req := CreateAccountRequest{
		InitialBalance: u.InitialBalance,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transfer request: %w", err)
	}
	resp, err := http.Post(utils.BASE_URL+"/accounts", "application/json", bytes.NewReader(reqBody))
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
	return &account, nil
}

func (u *BankUserImpl) TransferTo(toUser BankUser, amount float64) error {
	transferReq := TransferRequest{
		FromAccountID: u.AccountId,
		ToAccountID:   toUser.GetAccountId(),
		Amount:        amount,
	}

	reqBody, err := json.Marshal(transferReq)
	if err != nil {
		return fmt.Errorf("failed to marshal transfer request: %w", err)
	}

	resp, err := http.Post(utils.BASE_URL+"/accounts/transfer", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to perform transfer: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("transfer failed with status: %d", resp.StatusCode)
	}

	return nil
}

func (u *BankUserImpl) GetExpectedBalance(actions []action) float64 {
	balance := u.InitialBalance
	for _, act := range actions {
		if act.accountId == u.AccountId {
			balance += act.balanceChange
		}
	}
	return balance
}

func (u *BankUserImpl) GetAccountId() string {
	return u.AccountId
}

func (u *BankUserImpl) GetName() string {
	return u.Name
}
