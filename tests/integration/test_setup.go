package integration

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// projectRoot returns the absolute path to the project root
func projectRoot() string {
	// Start with the current working directory
	_, filename, _, _ := runtime.Caller(0)
	// Go up from tests/integration to the project root
	dir := filepath.Dir(filename)
	return filepath.Join(dir, "../..")
}

type TestDatabase struct {
	Container testcontainers.Container
	Pool      *pgxpool.Pool
	ConnStr   string
}

func SetupTestDatabase(ctx context.Context) (*TestDatabase, error) {
	// PostgreSQL container configuration
	req := testcontainers.ContainerRequest{
		Image:        "postgres:17-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	// Start container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	// Get host and port
	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, fmt.Errorf("failed to get container port: %w", err)
	}

	// Connection string
	connStr := fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable", host, port.Port())

	// Connect to database with retries
	var pool *pgxpool.Pool
	var connectErr error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		pool, connectErr = pgxpool.New(ctx, connStr)
		if connectErr == nil {
			// Try a simple query (ping) to ensure connectivity is truly established
			if pingErr := pool.Ping(ctx); pingErr == nil {
				log.Printf("Successfully connected to test database on attempt %d", i+1)
				break // Success!
			} else {
				// Ping failed, close the potentially problematic pool and prepare to retry
				pool.Close()
				connectErr = fmt.Errorf("ping failed after connect: %w", pingErr)
			}
		}

		// If we are here, connection or ping failed. Log and wait.
		log.Printf("Failed to connect/ping test database (attempt %d/%d), retrying in 1 second... Error: %v", i+1, maxRetries, connectErr)
		if i < maxRetries-1 {
			time.Sleep(1 * time.Second)
		}
	}
	// Check if connection ultimately failed after retries
	if connectErr != nil {
		if termErr := container.Terminate(ctx); termErr != nil {
			log.Printf("Failed to terminate container after connection failure: %v", termErr)
		}
		return nil, fmt.Errorf("failed to connect to database after %d retries: %w", maxRetries, connectErr)
	}

	// Initialize test database with migrations
	if err := applyMigrations(ctx, pool); err != nil {
		pool.Close()
		if termErr := container.Terminate(ctx); termErr != nil {
			log.Printf("Failed to terminate container after migration failure: %v", termErr)
		}
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return &TestDatabase{
		Container: container,
		Pool:      pool,
		ConnStr:   connStr,
	}, nil
}

func applyMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	// Path to migrations relative to the project root
	migrationsPath := filepath.Join(projectRoot(), "internal", "infrastructure", "persistence", "database", "migrations")
	log.Printf("Looking for migrations in: %s", migrationsPath)

	// Get all migration files
	var upMigrations []string
	err := filepath.Walk(migrationsPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".up.sql") {
			upMigrations = append(upMigrations, path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to read migrations: %w", err)
	}

	if len(upMigrations) == 0 {
		return fmt.Errorf("no migration files found in %s", migrationsPath)
	}

	// Sort migrations by filename (which should include a timestamp)
	sort.Strings(upMigrations)
	log.Printf("Found %d migration files", len(upMigrations))

	// Apply each migration
	for _, migrationPath := range upMigrations {
		migrationSQL, err := os.ReadFile(migrationPath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", migrationPath, err)
		}

		_, err = pool.Exec(ctx, string(migrationSQL))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", migrationPath, err)
		}

		log.Printf("Applied migration: %s", migrationPath)
	}

	return nil
}

func (db *TestDatabase) Close(ctx context.Context) {
	if db.Pool != nil {
		db.Pool.Close()
	}
	if db.Container != nil {
		if err := db.Container.Terminate(ctx); err != nil {
			log.Printf("Failed to terminate container: %v", err)
		}
	}
}

// TestMain sets up the test environment once for all tests
func TestMain(m *testing.M) {
	// Run the tests
	code := m.Run()
	os.Exit(code)
}
