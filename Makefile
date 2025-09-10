.PHONY: test test-containerized test-existing test-both bench bench-containerized bench-existing help

# Default target
help:
	@echo "Available targets:"
	@echo "  test-containerized  - Run tests using containerized database (default)"
	@echo "  test-existing      - Run tests using pre-existing database"
	@echo "  test-both          - Run tests with both database types"
	@echo "  bench-containerized - Benchmark using containerized database"
	@echo "  bench-existing     - Benchmark using pre-existing database"
	@echo ""
	@echo "Environment variables for existing DB (with defaults):"
	@echo "  DB_HOST=XXXXX"
	@echo "  DB_USER=postgres"
	@echo "  DB_PASSWORD=password"
	@echo "  DB_NAME=postgres"
	@echo ""
	@echo "Note: The existing DB tests connect to a TimescaleDB instance by default."
	@echo "Override with environment variables if you need different credentials."

# Test with containerized database (default behavior)
test-containerized:
	@echo "Running tests with containerized database..."
	go test -v ./...

# Test with pre-existing database
test-existing:
	@echo "Running tests with pre-existing database..."
	go test -v -existing-db ./...

# Run tests with both database types for comparison
test-both:
	@echo "Running tests with containerized database..."
	go test -v ./...
	@echo ""
	@echo "Running tests with pre-existing database..."
	go test -v -existing-db ./...

# Benchmark with containerized database
bench-containerized:
	@echo "Running benchmark with containerized database..."
	go test -bench=BenchmarkContainerizedDB -benchmem ./...

# Benchmark with pre-existing database
bench-existing:
	@echo "Running benchmark with pre-existing database..."
	go test -bench=BenchmarkExistingDB -benchmem -existing-db ./...

# Initialize Go modules
init:
	go mod tidy
	go mod download

