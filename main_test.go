package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"ignite-mesh-testing/pkg/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	useExistingDB = flag.Bool("existing-db", false, "Use pre-existing database instead of containerized one")
)

func TestMain(m *testing.M) {
	flag.Parse()
	code := m.Run()
	os.Exit(code)
}

func TestUserMigrationWithContainerizedDB(t *testing.T) {
	if *useExistingDB {
		t.Skip("Skipping containerized DB test when using existing DB")
	}

	start := time.Now()
	fmt.Printf("Starting containerized DB test at %s\n", start.Format(time.RFC3339))

	// Setup containerized database
	dbManager := database.NewContainerizedDB()
	defer func() {
		if err := dbManager.Cleanup(); err != nil {
			t.Logf("Failed to cleanup containerized DB: %v", err)
		}
	}()

	db, err := dbManager.Connect()
	require.NoError(t, err, "Failed to connect to containerized database")
	defer db.Close()

	// Run the migration test
	runMigrationTest(t, db, dbManager)

	duration := time.Since(start)
	fmt.Printf("Containerized DB test completed in %v\n", duration)
}

func TestUserMigrationWithExistingDB(t *testing.T) {
	if !*useExistingDB {
		t.Skip("Skipping existing DB test when using containerized DB")
	}

	start := time.Now()
	fmt.Printf("Starting existing DB test at %s\n", start.Format(time.RFC3339))

	// Setup existing database connection
	config := database.DBConfig{
		Host:     getEnvOrDefault("DB_HOST", "XXXXX"),
		Port:     5432,
		User:     getEnvOrDefault("DB_USER", "postgres"),
		Password: getEnvOrDefault("DB_PASSWORD", "password"),
		DBName:   getEnvOrDefault("DB_NAME", "postgres"),
		SSLMode:  "disable",
	}

	dbManager := database.NewExistingDB(config)
	db, err := dbManager.Connect()
	require.NoError(t, err, "Failed to connect to existing database")
	defer db.Close()

	// Run the migration test
	runMigrationTest(t, db, dbManager)

	duration := time.Since(start)
	fmt.Printf("Existing DB test completed in %v\n", duration)
}

func runMigrationTest(t *testing.T, db *sql.DB, dbManager database.DatabaseManager) {
	// Generate random table name
	tableName := dbManager.GetRandomTableName("test_users")
	fmt.Printf("Using table name: %s\n", tableName)

	// Ensure cleanup happens
	defer func() {
		if err := database.CleanupTable(db, tableName); err != nil {
			t.Logf("Failed to cleanup table %s: %v", tableName, err)
		} else {
			fmt.Printf("Successfully cleaned up table: %s\n", tableName)
		}
	}()

	// Run migration
	err := database.ExecuteMigration(db, tableName)
	require.NoError(t, err, "Failed to execute migration")

	// Verify table was created and data was inserted
	var count int
	err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)).Scan(&count)
	require.NoError(t, err, "Failed to query table")
	assert.Equal(t, 2, count, "Expected 2 users in the table")

	// Verify specific data
	var email, name string
	err = db.QueryRow(fmt.Sprintf("SELECT email, name FROM %s WHERE email = $1", tableName), "john@example.com").Scan(&email, &name)
	require.NoError(t, err, "Failed to query specific user")
	assert.Equal(t, "john@example.com", email)
	assert.Equal(t, "John Doe", name)

	// Test inserting additional data
	_, err = db.Exec(fmt.Sprintf("INSERT INTO %s (email, name) VALUES ($1, $2)", tableName), "test@example.com", "Test User")
	require.NoError(t, err, "Failed to insert additional user")

	// Verify total count
	err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)).Scan(&count)
	require.NoError(t, err, "Failed to query table after insert")
	assert.Equal(t, 3, count, "Expected 3 users in the table after insert")

	fmt.Printf("Migration test completed successfully for table: %s\n", tableName)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Benchmark functions to demonstrate performance differences
func BenchmarkContainerizedDB(b *testing.B) {
	if *useExistingDB {
		b.Skip("Skipping containerized DB benchmark when using existing DB")
	}

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		dbManager := database.NewContainerizedDB()
		db, err := dbManager.Connect()
		if err != nil {
			b.Fatalf("Failed to connect to containerized database: %v", err)
		}

		tableName := dbManager.GetRandomTableName("bench_users")

		b.StartTimer()

		// Run migration
		err = database.ExecuteMigration(db, tableName)
		if err != nil {
			b.Fatalf("Failed to execute migration: %v", err)
		}

		b.StopTimer()

		// Cleanup
		database.CleanupTable(db, tableName)
		db.Close()
		dbManager.Cleanup()
	}
}

func BenchmarkExistingDB(b *testing.B) {
	if !*useExistingDB {
		b.Skip("Skipping existing DB benchmark when using containerized DB")
	}

	config := database.DBConfig{
		Host:     getEnvOrDefault("DB_HOST", "XXXXX"),
		Port:     5432,
		User:     getEnvOrDefault("DB_USER", "postgres"),
		Password: getEnvOrDefault("DB_PASSWORD", "password"),
		DBName:   getEnvOrDefault("DB_NAME", "postgres"),
		SSLMode:  "disable",
	}

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		dbManager := database.NewExistingDB(config)
		db, err := dbManager.Connect()
		if err != nil {
			b.Fatalf("Failed to connect to existing database: %v", err)
		}

		tableName := dbManager.GetRandomTableName("bench_users")

		b.StartTimer()

		// Run migration
		err = database.ExecuteMigration(db, tableName)
		if err != nil {
			b.Fatalf("Failed to execute migration: %v", err)
		}

		b.StopTimer()

		// Cleanup
		database.CleanupTable(db, tableName)
		db.Close()
	}
}
