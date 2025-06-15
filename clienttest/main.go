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

// getConfigFromEnv reads RPS and DURATION from environment variables or uses defaults
func getConfigFromEnv() (int, int) {
	rps := DEFAULT_RPS
	duration := DEFAULT_DURATION

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

	return rps, duration
}

func main() {
	rps, testDuration := getConfigFromEnv()
	loadtest.AttackGetAccounts(rps, testDuration)
}
