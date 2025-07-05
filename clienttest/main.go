package main

import (
	"os"
	"strconv"

	"com.ndnhuy.mybank/loadtest"
)

const (
	DEFAULT_RPS      = 10 // Default requests per second
	DEFAULT_DURATION = 30 // Default duration in seconds
)

// getConfigFromEnv reads RPS, DURATION, NUM_SOURCE_CUSTOMERS, NUM_DEST_CUSTOMERS from environment variables or uses defaults
func getConfigFromEnv() (int, int, int, int, string) {
	rps := DEFAULT_RPS
	duration := DEFAULT_DURATION
	numSourceCustomers := 100
	numDestCustomers := 100
	scenario := "default"

	if envRPS := os.Getenv("RPS"); envRPS != "" {
		if parsed, err := strconv.Atoi(envRPS); err == nil && parsed > 0 {
			rps = parsed
		}
	}

	if envDuration := os.Getenv("DURATION"); envDuration != "" {
		if parsed, err := strconv.Atoi(envDuration); err == nil && parsed > 0 {
			duration = parsed
		}
	}

	if envNumSource := os.Getenv("NUM_SOURCE_CUSTOMERS"); envNumSource != "" {
		if parsed, err := strconv.Atoi(envNumSource); err == nil && parsed > 0 {
			numSourceCustomers = parsed
		}
	}

	if envNumDest := os.Getenv("NUM_DEST_CUSTOMERS"); envNumDest != "" {
		if parsed, err := strconv.Atoi(envNumDest); err == nil && parsed > 0 {
			numDestCustomers = parsed
		}
	}

	if envScenario := os.Getenv("SCENARIO"); envScenario != "" {
		scenario = envScenario
	}

	return rps, duration, numSourceCustomers, numDestCustomers, scenario
}

func main() {
	rps, testDuration, numSourceCustomers, numDestCustomers, scenario := getConfigFromEnv()

	loadtest.AttackTransfers(rps, testDuration, numSourceCustomers, numDestCustomers, scenario)
}
