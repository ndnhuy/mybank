package domain

// Action represents a customer action that can be executed and verified
type Action interface {
	Execute() error
	GetType() string
}
