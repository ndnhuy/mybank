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
	queueMetrics             *QueueMetrics    // QueueMetrics for queuing theory calculations
	customerTransferTargeter *CustomerTransferTargeter
	timeSeriesCallback       func() // Callback function to record time-series data
	timeSeriesReport         *TimeSeriesReport
}

func NewAttacker(targetURL string, method string, rps, durationInSeconds int, metrics *QueueMetrics) *Attacker {
	return &Attacker{
		targeter: vegeta.NewStaticTargeter(vegeta.Target{
			Method: method,
			URL:    targetURL,
		}),
		rate:         vegeta.Rate{Freq: rps, Per: time.Second},
		duration:     time.Duration(durationInSeconds) * time.Second,
		attacker:     vegeta.NewAttacker(),
		queueMetrics: metrics,
	}
}

// SetTimeSeriesCallback sets the callback function for recording time-series data
func (a *Attacker) SetTimeSeriesCallback(callback func()) {
	a.timeSeriesCallback = callback
}

func (a *Attacker) SetTimeSeriesReport(report *TimeSeriesReport) {
	a.timeSeriesReport = report
}

func (a *Attacker) Attack() {
	requestCount := 0
	lastRecordTime := time.Now()
	recordInterval := 5 * time.Second // Record metrics every 5 seconds

	for res := range a.attacker.Attack(a.targeter, a.rate, a.duration, "Load Test") {
		a.queueMetrics.metrics.Add(res)
		requestCount++

		// Print progress every 10 requests
		if requestCount%10 == 0 {
			fmt.Printf(".")
		}

		if res.Error == "" && res.Code == 200 {
			reqId := res.Headers[utils.X_REQUEST_ID][0]
			a.customerTransferTargeter.GetSuccessCallbackByRequestId(reqId)()
		}

		// Record time-series data at regular intervals
		if time.Since(lastRecordTime) >= recordInterval && a.timeSeriesReport != nil {
			a.timeSeriesReport.RecordPoint(a.getCurrentTimeseriesPoint())
			lastRecordTime = time.Now()
		}
	}

	// Record final data point
	if a.timeSeriesReport != nil {
		a.timeSeriesReport.RecordPoint(a.getCurrentTimeseriesPoint())
	}
}

func (a *Attacker) Duration() time.Duration {
	return a.duration
}

func (a *Attacker) getCurrentTimeseriesPoint() TimeSeriesPoint {
	now := time.Now()
	elapsed := now.Sub(a.queueMetrics.startTime).Seconds()

	// Create a snapshot of the current metrics to finalize without affecting the original
	snapshot := a.queueMetrics.Snapshot()

	point := TimeSeriesPoint{
		Timestamp:    now,
		Elapsed:      elapsed,
		MeanLatency:  float64(snapshot.Latencies.Mean.Nanoseconds()) / 1e6, // Convert to milliseconds
		P95Latency:   float64(snapshot.Latencies.P95.Nanoseconds()) / 1e6,
		P99Latency:   float64(snapshot.Latencies.P99.Nanoseconds()) / 1e6,
		Throughput:   snapshot.Throughput,
		SuccessRate:  snapshot.Success * 100,
		RequestCount: int64(snapshot.Requests),
	}

	return point
}
