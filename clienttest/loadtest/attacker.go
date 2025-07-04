package loadtest

import (
	"fmt"
	"time"

	"com.ndnhuy.mybank/utils"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type Attacker struct {
	targeter                 vegeta.Targeter  // Target URL for the load test
	rate                     vegeta.Rate      // Rate of requests per second
	duration                 time.Duration    // Duration of the load test in seconds
	attacker                 *vegeta.Attacker // Vegeta attacker instance
	metrics                  *vegeta.Metrics  // Pointer to metrics for accumulating results
	customerTransferTargeter *CustomerTransferTargeter
	timeSeriesCallback       func()           // Callback function to record time-series data
}

func NewAttacker(targetURL string, method string, rps, durationInSeconds int, metrics *vegeta.Metrics) *Attacker {
	return &Attacker{
		targeter: vegeta.NewStaticTargeter(vegeta.Target{
			Method: method,
			URL:    targetURL,
		}),
		rate:     vegeta.Rate{Freq: rps, Per: time.Second},
		duration: time.Duration(durationInSeconds) * time.Second,
		attacker: vegeta.NewAttacker(),
		metrics:  metrics,
	}
}

// SetTimeSeriesCallback sets the callback function for recording time-series data
func (a *Attacker) SetTimeSeriesCallback(callback func()) {
	a.timeSeriesCallback = callback
}

func (a *Attacker) Attack() {
	requestCount := 0
	lastRecordTime := time.Now()
	recordInterval := 5 * time.Second // Record metrics every 5 seconds
	
	for res := range a.attacker.Attack(a.targeter, a.rate, a.duration, "Load Test") {
		a.metrics.Add(res)
		requestCount++

		// Print progress every 10 requests
		if requestCount%10 == 0 {
			fmt.Printf(".")
		}

		// Record time-series data at regular intervals
		if time.Since(lastRecordTime) >= recordInterval && a.timeSeriesCallback != nil {
			a.timeSeriesCallback()
			lastRecordTime = time.Now()
		}

		if res.Error == "" && res.Code == 200 {
			reqId := res.Headers[utils.X_REQUEST_ID][0]
			a.customerTransferTargeter.GetSuccessCallbackByRequestId(reqId)()
		}
	}
	
	// Record final data point
	if a.timeSeriesCallback != nil {
		a.timeSeriesCallback()
	}
}

func (a *Attacker) Duration() time.Duration {
	return a.duration
}
