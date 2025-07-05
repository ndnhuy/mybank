package loadtest

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// TimeSeriesPoint represents a single measurement point in time
// (moved from queue_metrics.go)
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

// TimeSeriesReport manages time series points and CSV export
// (extracted from QueueMetrics)
type TimeSeriesReport struct {
	points    []TimeSeriesPoint
	csvWriter *csv.Writer
	csvFile   *os.File
}

// NewTimeSeriesReport creates a new TimeSeriesReport
func NewTimeSeriesReport() *TimeSeriesReport {
	return &TimeSeriesReport{
		points: make([]TimeSeriesPoint, 0),
	}
}

// InitCSV initializes CSV file for time-series data export
func (tsr *TimeSeriesReport) InitCSV(filepathStr string) error {
	dir := filepath.Dir(filepathStr)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	file, err := os.Create(filepathStr)
	if err != nil {
		return err
	}
	tsr.csvFile = file
	tsr.csvWriter = csv.NewWriter(file)
	header := []string{"Timestamp", "Elapsed_Seconds", "Mean_Latency_Ms", "P95_Latency_Ms", "P99_Latency_Ms", "Throughput_RPS", "Success_Rate_Percent", "Request_Count"}
	return tsr.csvWriter.Write(header)
}

// RecordPoint adds a new time series point and writes to CSV if enabled
func (tsr *TimeSeriesReport) RecordPoint(point TimeSeriesPoint) {
	tsr.points = append(tsr.points, point)
	if tsr.csvWriter != nil {
		record := []string{
			point.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
			strconv.FormatFloat(point.Elapsed, 'f', 2, 64),
			strconv.FormatFloat(point.MeanLatency, 'f', 2, 64),
			strconv.FormatFloat(point.P95Latency, 'f', 2, 64),
			strconv.FormatFloat(point.P99Latency, 'f', 2, 64),
			strconv.FormatFloat(point.Throughput, 'f', 2, 64),
			strconv.FormatFloat(point.SuccessRate, 'f', 2, 64),
			strconv.FormatInt(point.RequestCount, 10),
		}
		tsr.csvWriter.Write(record)
		tsr.csvWriter.Flush()
	}
}

// Close closes the CSV file if open
func (tsr *TimeSeriesReport) Close() {
	if tsr.csvWriter != nil {
		tsr.csvWriter.Flush()
	}
	if tsr.csvFile != nil {
		tsr.csvFile.Close()
	}
}
