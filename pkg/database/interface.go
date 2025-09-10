package database

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

// DBConfig holds database connection configuration
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// DatabaseManager handles database operations
type DatabaseManager interface {
	Connect() (*sql.DB, error)
	Cleanup() error
	GetConnectionString() string
	GetRandomTableName(prefix string) string
}

// GenerateRandomTableName creates a randomized table name with timestamp
func GenerateRandomTableName(prefix string) string {
	rand.Seed(time.Now().UnixNano())
	timestamp := time.Now().Format("20060102_150405")
	randomSuffix := rand.Intn(10000)
	return fmt.Sprintf("%s_%s_%d", prefix, timestamp, randomSuffix)
}

// ExecuteMigration runs the migration SQL on the given database
func ExecuteMigration(db *sql.DB, tableName string) error {
	migrationSQL := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		INSERT INTO %s (email, name) VALUES 
			('john@example.com', 'John Doe'),
			('jane@example.com', 'Jane Smith')
		ON CONFLICT (email) DO NOTHING;
	`, tableName, tableName)

	_, err := db.Exec(migrationSQL)
	return err
}

// CleanupTable drops the test table
func CleanupTable(db *sql.DB, tableName string) error {
	_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName))
	return err
}

