# Makefile for MyBank Load Testing Tool

# Default values
RPS ?= 10
DURATION ?= 30
ATTACK_TYPE ?= transfers

export RPS
export DURATION
export ATTACK_TYPE

.PHONY: run default custom clean baseline stress breaking visualize all-tests

default: run

run:
	@echo "Running load test attack=$${ATTACK_TYPE} with $${RPS} RPS for $${DURATION} seconds..."
	RPS=$${RPS} DURATION=$${DURATION} ATTACK_TYPE=$${ATTACK_TYPE} go run main.go

custom:
	@echo "Usage: make run RPS=50 DURATION=60"

clean:
	@echo "Cleaning up temporary files..."
	@rm -rf loadtest-main

tiny:
	@echo "Running baseline test..."
	RPS=10 DURATION=30 NUM_SOURCE_CUSTOMERS=1 NUM_DEST_CUSTOMERS=1 SCENARIO=tiny go run main.go > reports/tiny.txt

baseline:
	@echo "Running baseline test..."
	RPS=10 DURATION=60 SCENARIO=baseline go run main.go > reports/baseline.txt

stress:
	@echo "Running stress test..."
	RPS=100 DURATION=120 SCENARIO=stress go run main.go > reports/stress.txt

breaking:
	@echo "Running breaking test..."
	RPS=300 DURATION=120 SCENARIO=breaking go run main.go > reports/breaking.txt

visualize:
	@echo "Generating time-series visualizations..."
	@python3 visualize.py --reports-dir reports --output-dir reports
	@echo "Visualization complete! Check reports/ for generated plots."

all-tests:
	@echo "Running all load test scenarios..."
	@make clean
	@make baseline
	@make stress  
	@make breaking
	@make visualize
	@echo "All tests completed with visualizations!"
