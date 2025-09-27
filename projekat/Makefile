.PHONY: build run test clean install-deps lint

# Build the binary
build:
	go build -o bin/mini-zanzibar cmd/server/main.go

# Run the server
run:
	go run cmd/server/main.go

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -cover ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
install-deps:
	go mod download
	go mod tidy

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Run all checks
check: fmt vet lint test

# Docker build
docker-build:
	docker build -t mini-zanzibar .

# Docker run
docker-run:
	docker run -p 8080:8080 mini-zanzibar

# Development server with auto-reload (requires air)
dev:
	air