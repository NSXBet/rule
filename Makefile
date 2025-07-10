# Makefile for rule project
.PHONY: help test bench test-verbose test-coverage lint format clean install-tools mod-tidy build run-tests

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
test: ## Run all tests
	@echo "Running tests..."
	go test ./...

test-verbose: ## Run tests with verbose output
	@echo "Running tests (verbose)..."
	go test -v ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-fixtures: ## Run only fixture tests  
	@echo "Running fixture tests..."
	go test ./test/...

test-round1: ## Run round 1 tests specifically
	@echo "Running round 1 tests..."
	go test ./test -run TestRulesRound1

# Benchmarking
bench: ## Run all benchmarks
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

bench-comparison: ## Run comparison benchmarks only
	@echo "Running comparison benchmarks..."
	go test -bench=BenchmarkComparison -benchmem ./...

bench-optimized: ## Run optimized engine benchmarks
	@echo "Running optimized benchmarks..."  
	go test -bench=BenchmarkOptimized -benchmem ./...

bench-datetime: ## Run datetime benchmarks
	@echo "Running datetime benchmarks..."
	go test -bench=BenchmarkDateTime -benchmem ./...

# Code quality and formatting
format: ## Format code using gofumpt, goimports, and golines
	@echo "Formatting code..."
	find . -name "*.go" -not -path "./vendor/*" | xargs gofumpt -w
	find . -name "*.go" -not -path "./vendor/*" | xargs goimports -w
	find . -name "*.go" -not -path "./vendor/*" | xargs golines -w --max-len=120

lint: ## Run golangci-lint
	@echo "Running linter..."
	golangci-lint run

lint-fix: ## Run golangci-lint with auto-fix
	@echo "Running linter with auto-fix..."
	golangci-lint run --fix

# Performance and profiling
profile: ## Run CPU profiling
	@echo "Running CPU profile..."
	go test -bench=BenchmarkOptimizedEngineSimple -cpuprofile=cpu.prof ./...
	@echo "Profile saved to cpu.prof. View with: go tool pprof cpu.prof"

profile-mem: ## Run memory profiling
	@echo "Running memory profile..."
	go test -bench=BenchmarkOptimizedEngineSimple -memprofile=mem.prof ./...
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

pre-commit: format lint test ## Run pre-commit checks (format, lint, test)
	@echo "Pre-commit checks passed!"

ci: mod-tidy build lint test bench ## Run continuous integration pipeline
	@echo "CI pipeline completed successfully!"

# Quick commands for common workflows
quick-test: ## Quick test run (no verbose output)
	@go test ./...

quick-bench: ## Quick benchmark run
	@go test -bench=. ./... | grep -E "(Benchmark|PASS|FAIL)"

# Documentation
docs: ## Generate documentation
	@echo "Generating documentation..."
	go doc -all . > API.md
	@echo "Documentation generated: API.md"

# Git hooks
install-git-hooks: ## Install git pre-commit hooks
	@echo "Installing git hooks..."
	@echo '#!/bin/sh\nmake pre-commit' > .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "Git pre-commit hook installed"