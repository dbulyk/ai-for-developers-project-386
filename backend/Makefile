.PHONY: run build test lint fmt vet check tidy clean tools smoke

run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

test:
	go test ./...

lint:
	./bin/golangci-lint run ./...

fmt:
	gofmt -w .

vet:
	go vet ./...

check: fmt vet lint test

tidy:
	go mod tidy

clean:
	rm -rf bin/

tools:
	GOBIN=$(CURDIR)/bin go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

smoke:
	bash scripts/smoke.sh
