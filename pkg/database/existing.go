package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// ExistingDB manages connection to a pre-existing PostgreSQL database
type ExistingDB struct {
	config DBConfig
}

// NewExistingDB creates a new existing database instance
func NewExistingDB(config DBConfig) *ExistingDB {
	return &ExistingDB{
		config: config,
	}
}

// Connect creates a connection to the existing database
func (e *ExistingDB) Connect() (*sql.DB, error) {
	// Create database connection
	db, err := sql.Open("postgres", e.GetConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// Cleanup does nothing for existing databases (we don't want to terminate them)
func (e *ExistingDB) Cleanup() error {
	return nil
}

// GetConnectionString returns the PostgreSQL connection string
func (e *ExistingDB) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		e.config.Host, e.config.Port, e.config.User, e.config.Password, e.config.DBName, e.config.SSLMode)
}

// GetRandomTableName generates a random table name
func (e *ExistingDB) GetRandomTableName(prefix string) string {
	return GenerateRandomTableName(prefix)
}

