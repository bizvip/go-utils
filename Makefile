.PHONY: all vet test clean help

# Default target  
all: test

# Default target with vet
all-with-vet: vet test

# Run go vet (with automatic formatting)
vet:
	@echo "Running go fmt..."
	@go fmt ./...
	@echo "✅ Formatting completed successfully"
	@echo "Running go vet..."
	@go vet ./...
	@echo "✅ Vetting completed successfully"

# Run all tests from entire project (with automatic formatting)
test:
	@echo "Running go fmt..."
	@go fmt ./...
	@echo "✅ Formatting completed successfully"
	@echo "Running tests..."
	@go test -v ./...
	@echo "✅ All tests completed successfully"

# Run tests with coverage (with automatic formatting)
test-coverage:
	@echo "Running go fmt..."
	@go fmt ./...
	@echo "✅ Formatting completed successfully"
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Tests with coverage completed successfully"
	@echo "Coverage report generated: coverage.html"

# Run benchmarks (with automatic formatting)
bench:
	@echo "Running go fmt..."
	@go fmt ./...
	@echo "✅ Formatting completed successfully"
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...
	@echo "✅ Benchmarks completed successfully"

# Clean generated files
clean:
	@echo "Cleaning up..."
	@rm -f coverage.out coverage.html
	@echo "✅ Cleanup completed"

# Tidy go modules
tidy:
	@echo "Tidying go modules..."
	@go mod tidy
	@echo "✅ Go mod tidy completed successfully"

# Run full CI pipeline
ci: tidy vet test
	@echo "✅ CI pipeline completed successfully"

# Show help
help:
	@echo "Available targets:"
	@echo "  all           - Run tests (with auto-formatting) [default]"
	@echo "  all-with-vet  - Run vet + tests (with auto-formatting)"
	@echo "  vet           - Run go vet (with auto-formatting)"
	@echo "  test          - Run all unit tests from entire project (with auto-formatting)"
	@echo "  test-coverage - Run tests with coverage report (with auto-formatting)"
	@echo "  bench         - Run benchmarks (with auto-formatting)"
	@echo "  tidy          - Tidy go modules"
	@echo "  ci            - Run full CI pipeline (tidy, vet, test)"
	@echo "  clean         - Remove generated files"
	@echo "  help          - Show this help message"
	@echo ""
	@echo "Note: All test-related targets automatically run 'go fmt' first"