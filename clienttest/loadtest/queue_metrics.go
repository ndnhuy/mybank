package loadtest

import (
	"fmt"
	"strings"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

// QueueMetrics wraps vegeta.Metrics with additional queuing theory calculations
type QueueMetrics struct {
	*vegeta.Metrics
	startTime time.Time
}

// NewQueueMetrics creates a new QueueMetrics instance
func NewQueueMetrics() *QueueMetrics {
	return &QueueMetrics{
		Metrics:   &vegeta.Metrics{},
		startTime: time.Now(),
	}
}

// GetArrivalRate returns the arrival rate (λ) in requests/second
func (qm *QueueMetrics) GetArrivalRate() float64 {
	return qm.Rate
}

// GetServiceRate returns the service rate (μ) in requests/second
func (qm *QueueMetrics) GetServiceRate() float64 {
	return qm.Throughput
}

// GetTrafficIntensity returns the traffic intensity (ρ = λ/μ)
func (qm *QueueMetrics) GetTrafficIntensity() float64 {
	if qm.Throughput <= 0 {
		return 999.0 // Indicate overload
	}
	return qm.Rate / qm.Throughput
}

// GetObservationDuration returns how long the test ran
func (qm *QueueMetrics) GetObservationDuration() time.Duration {
	return qm.Duration
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
	mean := qm.Latencies.Mean
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
	fmt.Printf("   Requests Submitted:   %d\n", qm.Requests)
	fmt.Printf("   Requests Completed:   %d\n", qm.Requests) // Assuming all completed for HTTP
	fmt.Printf("   System Status:        %s\n", qm.GetSystemStatus())

	// Performance Metrics
	fmt.Println("\n⚡ PERFORMANCE METRICS:")
	fmt.Printf("   Mean Response Time:   %v\n", qm.Latencies.Mean)
	fmt.Printf("   50th percentile:      %v\n", qm.Latencies.P50)
	fmt.Printf("   95th percentile:      %v\n", qm.Latencies.P95)
	fmt.Printf("   99th percentile:      %v\n", qm.Latencies.P99)
	fmt.Printf("   Max Response Time:    %v\n", qm.Latencies.Max)
	fmt.Printf("   Success Rate:         %.2f%%\n", qm.Success*100)

	// Queuing Theory Analysis
	fmt.Println("\n🔬 QUEUING THEORY ANALYSIS:")
	fmt.Printf("   Arrival Rate (λ):      %.2f requests/sec\n", qm.GetArrivalRate())
	fmt.Printf("   Service Rate (μ):      %.2f requests/sec\n", qm.GetServiceRate())
	fmt.Printf("   Traffic Intensity (ρ): %.3f\n", qm.GetTrafficIntensity())
	fmt.Printf("   Throughput:           %.2f requests/sec\n", qm.Throughput)

	// Performance Assessment
	fmt.Println("\n🎯 PERFORMANCE ASSESSMENT:")
	fmt.Printf("   Response Time:  %s\n", qm.AssessResponseTime())
	fmt.Printf("   System Health:  %s\n", qm.AssessSystemHealth())

	// Warnings
	if qm.Success < 1.0 {
		fmt.Printf("\n⚠️  Warning: Success rate is %.2f%%. Some requests failed!\n", qm.Success*100)
	}
	if qm.GetTrafficIntensity() >= 1.0 {
		fmt.Println("\n⚠️  SYSTEM OVERLOADED!")
	}

	fmt.Println("\n" + strings.Repeat("═", 66))
}
