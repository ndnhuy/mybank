package domain

// BankUser defines the interface for user operations in the banking domain.
type BankUser interface {
	GetAccount(accountID string) (*AccountInfo, error)
	GetAccountBalance() (float64, error)
	CreateAccount() (*AccountInfo, error)
	TransferTo(toUser BankUser, amount float64) error
	GetExpectedBalance(actions []action) float64

	GetAccountId() string
	GetName() string
}
