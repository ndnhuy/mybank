package domain 

// AccountInfo represents the account information returned by the API
type AccountInfo struct {
	ID      string  `json:"id"`
	Balance float64 `json:"balance"`
}

// TransferRequest represents the transfer request payload
type TransferRequest struct {
	FromAccountID string  `json:"fromAccountId"`
	ToAccountID   string  `json:"toAccountId"`
	Amount        float64 `json:"amount"`
}

type CreateAccountRequest struct {
	InitialBalance float64 `json:"initialBalance"`
}
