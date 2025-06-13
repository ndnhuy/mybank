package domain

// BankOperator defines the interface for user operations in the banking domain.
type BankOperator interface {
	GetAccount(accountID string) (*AccountInfo, error)
	GetAccountBalance() (float64, error)
	CreateAccount() (*AccountInfo, error)
	TransferTo(toUser BankOperator, amount float64) error
	GetExpectedBalance(actions []action) float64

	GetAccountId() string
	GetName() string
}
