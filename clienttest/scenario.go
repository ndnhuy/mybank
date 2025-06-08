package main

import (
	"fmt"
	"log"

	mybankerror "com.ndnhuy.mybank/mybankerror"
)

type Scenario interface {
	// Run executes the test scenario
	Run() error
	Expectation() interface{}
}

type TransferScenario struct {
	Name     string
	FromUser *User
	ToUser   *User
	Amount   float64

	actions []action // actions to be performed in the scenario
}

type action struct {
	accountId     string
	accountName   string
	balanceChange float64 // positive for deposit, negative for withdrawal
}

func (a action) String() string {
	return fmt.Sprintf("Account %v %v", a.accountName, a.balanceChange)
}

func (s *TransferScenario) Run() error {
	if s.Name != "" {
		log.Printf("Start scenario: %s", s.Name)
	}

	// Create accounts for both users
	_, err := s.FromUser.CreateAccount()
	if err != nil && err != mybankerror.AccountAlreadyCreatedError {
		return err
	}

	_, err = s.ToUser.CreateAccount()
	if err != nil && err != mybankerror.AccountAlreadyCreatedError {
		return err
	}

	// Perform the transfer
	s.actions = append(s.actions, action{
		accountId:     s.FromUser.AccountId,
		accountName:   s.FromUser.Name,
		balanceChange: -s.Amount, // negative for withdrawal
	})
	s.actions = append(s.actions, action{
		accountId:     s.ToUser.AccountId,
		accountName:   s.ToUser.Name,
		balanceChange: s.Amount, // positive for deposit
	})

	if err := s.FromUser.TransferTo(s.ToUser, s.Amount); err != nil {
		log.Println("Transfer failed:", err)
		return err
	}
	log.Printf("Transfer of %.2f from %s to %s completed successfully", s.Amount, s.FromUser.Name, s.ToUser.Name)

	return nil
}

func (s *TransferScenario) Expectation() interface{} {
	return s.actions
}
