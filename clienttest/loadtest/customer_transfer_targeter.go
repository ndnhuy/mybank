package loadtest

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"sync"

	"com.ndnhuy.mybank/domain"
	"com.ndnhuy.mybank/utils"
	"github.com/google/uuid"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

// ExtendedTarget embeds vegeta.Target and adds OnSuccess callback
// OnSuccess is called when the attack receives a 200 response
// Only used for business logic hooks
type ExtendedTarget struct {
	vegeta.Target
	OnSuccess func()
}

// CustomerTransferTargeter creates transfer requests using customer behaviors
// successCallbacks is now a field, concurrent safe
type CustomerTransferTargeter struct {
	sourceCustomers []*domain.Customer
	destCustomers   []*domain.Customer

	targeter vegeta.Targeter // Base targeter for generating requests

	successCallbacks struct {
		sync.RWMutex
		m map[string]func()
	}
}

// NewCustomerTransferTargeter creates a new customer-based transfer targeter
func NewCustomerTransferTargeter(sourceCustomers, destCustomers []*domain.Customer) *CustomerTransferTargeter {
	tt := &CustomerTransferTargeter{
		sourceCustomers: sourceCustomers,
		destCustomers:   destCustomers,
	}
	tt.successCallbacks.m = make(map[string]func())

	tt.targeter = func(t *vegeta.Target) error {
		extendedTarget := tt.generateTarget()
		*t = extendedTarget.Target
		tt.AddSuccessCallbackForRequestId(extendedTarget.Header[utils.X_REQUEST_ID][0], extendedTarget.OnSuccess)
		return nil
	}
	return tt
}

// generateTarget generates a random transfer request using customer account IDs
func (tt *CustomerTransferTargeter) generateTarget() ExtendedTarget {
	// Pick random source and destination customers
	fromIdx := rand.Intn(len(tt.sourceCustomers))
	toIdx := rand.Intn(len(tt.destCustomers))
	fromCustomer := tt.sourceCustomers[fromIdx]
	toCustomer := tt.destCustomers[toIdx]

	transferReq := domain.TransferRequest{
		FromAccountID: fromCustomer.GetAccountID(),
		ToAccountID:   toCustomer.GetAccountID(),
		Amount:        1,
	}

	body, _ := json.Marshal(transferReq)

	// fromCustomer.RecordTransfer(toCustomer, transferReq.Amount) // Record transfer for source customer

	return ExtendedTarget{
		Target: vegeta.Target{
			Method: "POST",
			URL:    utils.BASE_URL + "/accounts/transfer",
			Header: http.Header{
				"Content-Type":     []string{"application/json"},
				utils.X_REQUEST_ID: []string{uuid.NewString()},
			},
			Body: body,
		},
		OnSuccess: func() {
			fromCustomer.RecordTransfer(toCustomer, transferReq.Amount) // Record transfer for source customer
		},
	}
}

// GetSuccessCallbackByRequestId returns the success callback for a given request ID, or nil if not found
func (tt *CustomerTransferTargeter) GetSuccessCallbackByRequestId(requestId string) func() {
	tt.successCallbacks.RLock()
	defer tt.successCallbacks.RUnlock()
	return tt.successCallbacks.m[requestId]
}

// AddSuccessCallbackForRequestId adds a success callback for a given request ID
func (tt *CustomerTransferTargeter) AddSuccessCallbackForRequestId(requestId string, cb func()) {
	tt.successCallbacks.Lock()
	defer tt.successCallbacks.Unlock()
	tt.successCallbacks.m[requestId] = cb
}
