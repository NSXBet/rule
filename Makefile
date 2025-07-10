# Makefile for rule project
.PHONY: help test bench test-verbose test-coverage lint format clean install-tools mod-tidy build run-tests fuzz

# Default target
.DEFAULT_GOAL := help

help: ## Show this help message
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# Dependencies and tools
install-tools: ## Install required development tools
	@echo "Installing development tools..."
	go install golang.org/x/tools/cmd/goimports@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/segmentio/golines@latest
	@echo "Tools installed successfully"

# Module management
mod-tidy: ## Clean up go.mod and go.sum
	@echo "Tidying modules..."
	go mod tidy
	go mod download

# Building
build: ## Build the project
	@echo "Building..."
	go build ./...

# Testing
test: ## Run tests with verbose output
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./... && go tool cover -func=coverage.out
	@go tool cover -func=coverage.out

race:
	@echo "Running tests with coverage and race condition detection..."
	@go test -race -v -coverprofile=coverage.out ./... && go tool cover -func=coverage.out

examples: ## Run all example tests
	@echo "Running example tests..."
	go test -run Example -v

fuzz: ## Run fuzz tests for comprehensive edge case detection
	@echo "Running fuzz tests..."
	@echo "Testing rule execution parsing..."
	go test -fuzz=FuzzRuleExecution -fuzztime=10s
	@echo "Testing string operations..."
	go test -fuzz=FuzzStringOperations -fuzztime=5s
	@echo "Testing numeric operations..."
	go test -fuzz=FuzzNumericOperations -fuzztime=5s
	@echo "Testing datetime operations..."
	go test -fuzz=FuzzDateTimeOperations -fuzztime=5s
	@echo "Testing property access..."
	go test -fuzz=FuzzPropertyAccess -fuzztime=5s
	@echo "Testing array operations..."
	go test -fuzz=FuzzArrayOperations -fuzztime=5s
	@echo "Testing boolean operations..."
	go test -fuzz=FuzzBooleanOperations -fuzztime=3s
	@echo "Testing complex rules..."
	go test -fuzz=FuzzComplexRules -fuzztime=5s
	@echo "Testing mixed type comparisons..."
	go test -fuzz=FuzzMixedTypeComparisons -fuzztime=5s
	@echo "Fuzz testing completed successfully!"

# Benchmarking
bench: ## Run all benchmarks
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Code quality and formatting
format: ## Format code using gofumpt, goimports, and golines
	@echo "Formatting code..."
	@find . -name "*.go" -not -path "./vendor/*" | xargs gofumpt -w
	@find . -name "*.go" -not -path "./vendor/*" | xargs goimports -w
	@find . -name "*.go" -not -path "./vendor/*" | xargs golines -w --max-len=120

lint: ## Run golangci-lint
	@echo "Running linter..."
	@golangci-lint run --fix

# Performance and profiling
profile: ## Run CPU profiling
	@echo "Running CPU profile..."
	@go test -bench=BenchmarkOptimizedEngineSimple -cpuprofile=cpu.prof ./...
	@echo "Profile saved to cpu.prof. View with: go tool pprof cpu.prof"

profile-mem: ## Run memory profiling
	@echo "Running memory profile..."
	@go test -bench=BenchmarkOptimizedEngineSimple -memprofile=mem.prof ./...
	@echo "Profile saved to mem.prof. View with: go tool pprof mem.prof"

# Cleanup
clean: ## Clean build artifacts and temp files
	@echo "Cleaning up..."
	go clean
	rm -f coverage.out coverage.html
	rm -f cpu.prof mem.prof
	rm -f *.test

# Development workflow
dev-setup: install-tools mod-tidy ## Setup development environment
	@echo "Development environment setup complete"
