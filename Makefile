# Makefile for building the Docker image

.PHONY: build

build:
	docker build -t mybank .

.PHONY: clean

clean:
	docker rm -f mybank 2>/dev/null || true

.PHONY: run

run: clean
	docker run -d -p 8080:8080 --name mybank mybank

.PHONY: test

test:
	go test -count=1 ./...
