# Makefile for MyBank Load Testing Tool

# Default values
RPS ?= 10
DURATION ?= 30
ATTACK_TYPE ?= transfers

export RPS
export DURATION
export ATTACK_TYPE

.PHONY: run default custom clean

default: run

run:
	@echo "Running load test attack=$${ATTACK_TYPE} with $${RPS} RPS for $${DURATION} seconds..."
	RPS=$${RPS} DURATION=$${DURATION} ATTACK_TYPE=$${ATTACK_TYPE} go run main.go

custom:
	@echo "Usage: make run RPS=50 DURATION=60"

clean:
	@echo "Cleaning up temporary files..."
	@rm -rf loadtest-main
