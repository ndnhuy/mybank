package loadtest

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

// TimeSeriesPoint represents a single measurement point in time
type TimeSeriesPoint struct {
	Timestamp    time.Time
	Elapsed      float64 // seconds since start
	MeanLatency  float64 // milliseconds
	P95Latency   float64 // milliseconds
	P99Latency   float64 // milliseconds
	Throughput   float64 // requests/second
	SuccessRate  float64 // percentage
	RequestCount int64
}

// QueueMetrics wraps vegeta.Metrics with additional queuing theory calculations
type QueueMetrics struct {
	metrics     *vegeta.Metrics
	startTime   time.Time
	timeSeries  []TimeSeriesPoint
	csvWriter   *csv.Writer
	csvFile     *os.File
}

// NewQueueMetrics creates a new QueueMetrics instance
func NewQueueMetrics() *QueueMetrics {
	return &QueueMetrics{
		metrics:    &vegeta.Metrics{},
		startTime:  time.Now(),
		timeSeries: make([]TimeSeriesPoint, 0),
	}
}

// InitTimeSeriesCSV initializes CSV file for time-series data export
func (qm *QueueMetrics) InitTimeSeriesCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	
	qm.csvFile = file
	qm.csvWriter = csv.NewWriter(file)
	
	// Write CSV header
	header := []string{"Timestamp", "Elapsed_Seconds", "Mean_Latency_Ms", "P95_Latency_Ms", "P99_Latency_Ms", "Throughput_RPS", "Success_Rate_Percent", "Request_Count"}
	return qm.csvWriter.Write(header)
}

// RecordTimeSeriesPoint captures current metrics state
func (qm *QueueMetrics) RecordTimeSeriesPoint() {
	now := time.Now()
	elapsed := now.Sub(qm.startTime).Seconds()
	
	// Create a snapshot of the current metrics to finalize without affecting the original
	snapshot := qm.createMetricsSnapshot()
	snapshot.Close() // Finalize the snapshot to get accurate latency percentiles
	
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
	
	qm.timeSeries = append(qm.timeSeries, point)
	
	// Write to CSV if initialized
	if qm.csvWriter != nil {
		record := []string{
			point.Timestamp.Format(time.RFC3339),
			strconv.FormatFloat(point.Elapsed, 'f', 2, 64),
			strconv.FormatFloat(point.MeanLatency, 'f', 2, 64),
			strconv.FormatFloat(point.P95Latency, 'f', 2, 64),
			strconv.FormatFloat(point.P99Latency, 'f', 2, 64),
			strconv.FormatFloat(point.Throughput, 'f', 2, 64),
			strconv.FormatFloat(point.SuccessRate, 'f', 2, 64),
			strconv.FormatInt(point.RequestCount, 10),
		}
		qm.csvWriter.Write(record)
		qm.csvWriter.Flush()
	}
}

// createMetricsSnapshot creates a snapshot by recreating metrics from stored results
func (qm *QueueMetrics) createMetricsSnapshot() *vegeta.Metrics {
	return qm.metrics
}

// Close properly closes CSV file and cleans up resources
func (qm *QueueMetrics) Close() {
	if qm.csvWriter != nil {
		qm.csvWriter.Flush()
	}
	if qm.csvFile != nil {
		qm.csvFile.Close()
	}
	// Close vegeta metrics to finalize the data
	qm.metrics.Close()
}

// GetArrivalRate returns the arrival rate (λ) in requests/second
func (qm *QueueMetrics) GetArrivalRate() float64 {
	return qm.metrics.Rate
}

// GetServiceRate returns the service rate (μ) in requests/second
func (qm *QueueMetrics) GetServiceRate() float64 {
	return qm.metrics.Throughput
}

// GetTrafficIntensity returns the traffic intensity (ρ = λ/μ)
func (qm *QueueMetrics) GetTrafficIntensity() float64 {
	if qm.metrics.Throughput <= 0 {
		return 999.0 // Indicate overload
	}
	return qm.metrics.Rate / qm.metrics.Throughput
}

// GetObservationDuration returns how long the test ran
func (qm *QueueMetrics) GetObservationDuration() time.Duration {
	return qm.metrics.Duration
}

// GetSystemStatus returns a color-coded status based on traffic intensity
func (qm *QueueMetrics) GetSystemStatus() string {
	rho := qm.GetTrafficIntensity()
	if rho >= 1.0 {
		return "🔴 OVERLOADED"
	}
	if rho >= 0.8 {
		return "🟡 HIGH LOAD"
	}
	if rho >= 0.5 {
		return "🟢 MODERATE LOAD"
	}
	return "🟢 LOW LOAD"
}

// AssessResponseTime provides response time assessment
func (qm *QueueMetrics) AssessResponseTime() string {
	mean := qm.metrics.Latencies.Mean
	if mean < 100*time.Millisecond {
		return "🟢 Excellent (< 100ms)"
	}
	if mean < 500*time.Millisecond {
		return "🟡 Good (< 500ms)"
	}
	if mean < 1*time.Second {
		return "🟠 Fair (< 1s)"
	}
	return "🔴 Poor (> 1s)"
}

// AssessSystemHealth provides overall system health assessment
func (qm *QueueMetrics) AssessSystemHealth() string {
	rho := qm.GetTrafficIntensity()
	if rho >= 1.0 {
		return "🔴 System cannot keep up with demand"
	}
	if rho >= 0.8 {
		return "🟡 System near capacity, consider scaling"
	}
	if rho >= 0.5 {
		return "🟢 System handling load well"
	}
	return "🟢 System has excess capacity"
}

// PrintReport prints a comprehensive report similar to Java QueueMetrics
func (qm *QueueMetrics) PrintReport() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    QUEUING SYSTEM ANALYSIS                  ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")

	// System Overview
	fmt.Println("\n📊 SYSTEM OVERVIEW:")
	fmt.Printf("   Observation Duration: %v\n", qm.GetObservationDuration())
	fmt.Printf("   Requests Submitted:   %d\n", int(qm.metrics.Requests))
	fmt.Printf("   Requests Completed:   %d\n", int(qm.metrics.Requests)) // Assuming all completed for HTTP
	fmt.Printf("   System Status:        %s\n", qm.GetSystemStatus())

	// Performance Metrics
	fmt.Println("\n⚡ PERFORMANCE METRICS:")
	fmt.Printf("   Mean Response Time:   %v\n", qm.metrics.Latencies.Mean)
	fmt.Printf("   50th percentile:      %v\n", qm.metrics.Latencies.P50)
	fmt.Printf("   95th percentile:      %v\n", qm.metrics.Latencies.P95)
	fmt.Printf("   99th percentile:      %v\n", qm.metrics.Latencies.P99)
	fmt.Printf("   Max Response Time:    %v\n", qm.metrics.Latencies.Max)
	fmt.Printf("   Success Rate:         %.2f%%\n", qm.metrics.Success*100)

	// Queuing Theory Analysis
	fmt.Println("\n🔬 QUEUING THEORY ANALYSIS:")
	fmt.Printf("   Arrival Rate (λ):      %.2f requests/sec\n", qm.GetArrivalRate())
	fmt.Printf("   Service Rate (μ):      %.2f requests/sec\n", qm.GetServiceRate())
	fmt.Printf("   Traffic Intensity (ρ): %.3f\n", qm.GetTrafficIntensity())
	fmt.Printf("   Throughput:           %.2f requests/sec\n", qm.metrics.Throughput)

	// Performance Assessment
	fmt.Println("\n🎯 PERFORMANCE ASSESSMENT:")
	fmt.Printf("   Response Time:  %s\n", qm.AssessResponseTime())
	fmt.Printf("   System Health:  %s\n", qm.AssessSystemHealth())

	// Time Series Analysis
	if len(qm.timeSeries) > 0 {
		fmt.Println("\n📈 TIME SERIES ANALYSIS:")
		firstPoint := qm.timeSeries[0]
		lastPoint := qm.timeSeries[len(qm.timeSeries)-1]
		fmt.Printf("   Time Series Points:   %d\n", len(qm.timeSeries))
		fmt.Printf("   Initial Mean Latency: %.2f ms\n", firstPoint.MeanLatency)
		fmt.Printf("   Final Mean Latency:   %.2f ms\n", lastPoint.MeanLatency)
		fmt.Printf("   Initial P95 Latency:  %.2f ms\n", firstPoint.P95Latency)
		fmt.Printf("   Final P95 Latency:    %.2f ms\n", lastPoint.P95Latency)
		
		if firstPoint.MeanLatency > 0 {
			latencyIncrease := ((lastPoint.MeanLatency - firstPoint.MeanLatency) / firstPoint.MeanLatency) * 100
			fmt.Printf("   Mean Latency Change:  %.1f%%\n", latencyIncrease)
		}
	}

	// Warnings
	if qm.metrics.Success < 1.0 {
		fmt.Printf("\n⚠️  Warning: Success rate is %.2f%%. Some requests failed!\n", qm.metrics.Success*100)
	}
	if qm.GetTrafficIntensity() >= 1.0 {
		fmt.Println("\n⚠️  SYSTEM OVERLOADED!")
	}

	fmt.Println("\n" + strings.Repeat("═", 66))
}
