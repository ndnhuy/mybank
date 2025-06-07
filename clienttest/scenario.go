package main

import "log"

type Scenario interface {
	// Run executes the test scenario
	Run() error
	Expectation() interface{}
}

type TransferScenario struct {
	FromUser *User
	ToUser   *User
	Amount   float64

	actions []action // actions to be performed in the scenario
}

type action struct {
	accountId     string
	balanceChange float64 // positive for deposit, negative for withdrawal
}

func (s *TransferScenario) Run() error {
	// Create accounts for both users
	fromAccount, err := s.FromUser.CreateAccount()
	if err != nil {
		return err
	}

	toAccount, err := s.ToUser.CreateAccount()
	if err != nil {
		return err
	}

	// Perform the transfer
	s.actions = append(s.actions, action{
		accountId:     fromAccount.ID,
		balanceChange: -s.Amount, // negative for withdrawal
	})
	s.actions = append(s.actions, action{
		accountId:     toAccount.ID,
		balanceChange: s.Amount, // positive for deposit
	})

	if err := s.FromUser.TransferTo(fromAccount.ID, toAccount.ID, s.Amount); err != nil {
		log.Println("Transfer failed:", err)
		return err
	}
	log.Printf("Transfer of %.2f from %s to %s completed successfully", s.Amount, fromAccount.ID, toAccount.ID)

	return nil
}

func (s *TransferScenario) Expectation() interface{} {
	return s.actions
}
