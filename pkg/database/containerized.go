package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// ContainerizedDB manages a containerized PostgreSQL database
type ContainerizedDB struct {
	container testcontainers.Container
	config    DBConfig
}

// NewContainerizedDB creates a new containerized database instance
func NewContainerizedDB() *ContainerizedDB {
	return &ContainerizedDB{}
}

// Connect starts a PostgreSQL container and returns a database connection
func (c *ContainerizedDB) Connect() (*sql.DB, error) {
	ctx := context.Background()

	// Create PostgreSQL container
	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Minute)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	c.container = postgresContainer

	// Get connection details
	host, err := postgresContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, fmt.Errorf("failed to get container port: %w", err)
	}

	c.config = DBConfig{
		Host:     host,
		Port:     port.Int(),
		User:     "testuser",
		Password: "testpass",
		DBName:   "testdb",
		SSLMode:  "disable",
	}

	// Create database connection
	db, err := sql.Open("postgres", c.GetConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// Cleanup terminates the container
func (c *ContainerizedDB) Cleanup() error {
	if c.container != nil {
		ctx := context.Background()
		return c.container.Terminate(ctx)
	}
	return nil
}

// GetConnectionString returns the PostgreSQL connection string
func (c *ContainerizedDB) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.config.Host, c.config.Port, c.config.User, c.config.Password, c.config.DBName, c.config.SSLMode)
}

// GetRandomTableName generates a random table name
func (c *ContainerizedDB) GetRandomTableName(prefix string) string {
	return GenerateRandomTableName(prefix)
}

