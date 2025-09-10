# Database Testing Comparison Demo

This project demonstrates the performance and setup differences between using ad-hoc containerized databases versus pre-existing databases for testing in Go.

> This project is brought to you by [groundcover](https://www.groundcover.com).

## Features

- ✅ Simple user table migration
- ✅ Randomized table names for test isolation
- ✅ Automatic cleanup after tests
- ✅ Command-line flag to switch between database modes
- ✅ Performance benchmarking
- ✅ Both containerized (testcontainers) and existing database support

## Project Structure

```
.
├── pkg/database/
│   ├── interface.go      # Database interface and common functions
│   ├── containerized.go  # Testcontainers implementation
│   └── existing.go       # Pre-existing database implementation
├── migrations/
│   └── 001_create_users.sql  # SQL migration file
├── main_test.go          # Test suite
├── Makefile             # Build and test commands
└── README.md            # This file
```

## Quick Start

### 1. Initialize Dependencies

```bash
go mod tidy
```

### 2. Run Tests with Containerized Database (Default)

```bash
make test-containerized
# or simply
go test -v ./...
```

### 3. Run Tests with Pre-existing Database

The tests are configured to use a TimescaleDB instance by default:

```bash
make test-existing
# or
go test -v -existing-db ./...
```

The default connection uses:
- Host: `XXXXX:5432`
- User: `postgres`
- Password: `password`
- Database: `postgres`

### 4. Compare Both Approaches

```bash
make test-both
```

### 5. Run the Demo

For a quick demonstration of the performance difference:

```bash
./demo.sh
```

## Configuration

### Environment Variables for Pre-existing Database

When using the `-existing-db` flag, you can configure the database connection:

- `DB_HOST` (default: XXXXX)
- `DB_USER` (default: postgres)  
- `DB_PASSWORD` (default: password)
- `DB_NAME` (default: postgres)

Example:
```bash
DB_HOST=mydb.example.com DB_USER=testuser go test -v -existing-db ./...
```

## Performance Benchmarking

### Benchmark Containerized Database
```bash
make bench-containerized
```

### Benchmark Pre-existing Database
```bash
make bench-existing
```

## Key Differences Demonstrated

### Containerized Database (Testcontainers)
- **Pros:**
  - Complete isolation
  - No external dependencies
  - Consistent environment
  - Easy CI/CD integration
- **Cons:**
  - Slower startup time
  - Higher resource usage
  - Requires Docker

### Pre-existing Database
- **Pros:**
  - Faster test execution
  - Lower resource overhead
  - Shared across multiple test runs
- **Cons:**
  - Requires external setup
  - Potential for test interference
  - Environment dependencies

## Test Features

- **Randomized Table Names**: Each test creates a uniquely named table (e.g., `test_users_20240115_143022_1234`)
- **Automatic Cleanup**: Tables are dropped after each test
- **Migration Testing**: Tests the creation and population of a simple users table
- **Data Validation**: Verifies both table structure and data integrity

## Example Output

```
Starting containerized DB test at 2024-01-15T14:30:22Z
Using table name: test_users_20240115_143022_1234
Migration test completed successfully for table: test_users_20240115_143022_1234
Successfully cleaned up table: test_users_20240115_143022_1234
Containerized DB test completed in 15.234s

Starting existing DB test at 2024-01-15T14:30:45Z  
Using table name: test_users_20240115_143045_5678
Migration test completed successfully for table: test_users_20240115_143045_5678
Successfully cleaned up table: test_users_20240115_143045_5678
Existing DB test completed in 234ms
```

## Dependencies

- Go 1.21+
- Docker (for containerized tests)
- PostgreSQL driver (`github.com/lib/pq`)
- Testcontainers (`github.com/testcontainers/testcontainers-go`)
- Testify (`github.com/stretchr/testify`)

## Cleanup

- **Containerized databases** (testcontainers) are automatically cleaned up after each test
- **TimescaleDB instance** is persistent and shared - test tables are cleaned up automatically
- **Manual Docker containers** (if used): `docker stop postgres && docker rm postgres`

