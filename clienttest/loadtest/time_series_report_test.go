package loadtest

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitCSV_CreatesFileAndDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	dirPath := filepath.Join(tmpDir, "nested/dir")
	csvPath := filepath.Join(dirPath, "timeseries_test.csv")

	report := NewTimeSeriesReport("test", dirPath)
	err := report.InitCSV()
	if err != nil {
		t.Fatalf("InitCSV failed: %v", err)
	}

	// Check file exists
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		t.Errorf("CSV file was not created at %s", csvPath)
	}

	// Check directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		t.Errorf("Directory was not created at %s", dirPath)
	}

	// Clean up
	report.Close()
}

func TestInitCSV_InvalidPath(t *testing.T) {
	// Use an invalid path (simulate by using a directory as a file)
	dir := t.TempDir()
	report := NewTimeSeriesReport("test", dir)
	err := report.InitCSV() // Should fail because dir is a directory
	if err == nil {
		t.Error("Expected error when using directory as file path, got nil")
	}
}
