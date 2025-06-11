package loadtest

import (
	"fmt"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type Attacker struct {
	targeter vegeta.Targeter  // Target URL for the load test
	rate     vegeta.Rate      // Rate of requests per second
	duration time.Duration    // Duration of the load test in seconds
	attacker *vegeta.Attacker // Vegeta attacker instance
	metrics  *vegeta.Metrics  // Pointer to metrics for accumulating results
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

func (a *Attacker) Attack() {
	requestCount := 0
	for res := range a.attacker.Attack(a.targeter, a.rate, a.duration, "Load Test") {
		a.metrics.Add(res)
		requestCount++

		// Print progress every 10 requests
		if requestCount%10 == 0 {
			fmt.Printf(".")
		}

		// Log any non-200 responses
		if res.Code != 200 {
			fmt.Printf("\n⚠️  Non-200 response: %d %s\n", res.Code, res.Error)
		}
	}
}

func (a *Attacker) Duration() time.Duration {
	return a.duration
}
