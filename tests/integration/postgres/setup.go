package postgres

import (
	"context"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"runtime"
	"smart-hub/internal/common/database"
	"smart-hub/internal/common/database/migrations"
	"testing"
	"time"
)

type TestConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func DefaultTestConfig() TestConfig {
	return TestConfig{
		DBHost:     getEnv("TEST_DB_HOST", "localhost"),
		DBPort:     getEnv("TEST_DB_PORT", "5433"),
		DBUser:     getEnv("TEST_DB_USER", "postgres"),
		DBPassword: getEnv("TEST_DB_PASSWORD", "postgres"),
		DBName:     getEnv("TEST_DB_NAME", "smart_hub_test"),
	}
}

func ProjectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "../../..")
	return projectRoot
}

func (c TestConfig) GetDSN() string {
	return "postgresql://" + c.DBUser + ":" + c.DBPassword + "@" + c.DBHost + ":" + c.DBPort + "/" + c.DBName + "?sslmode=disable"
}

func SetupTestDB(t *testing.T) database.Database {
	config := DefaultTestConfig()

	var db database.Database
	var err error

	for i := 0; i < 5; i++ {
		db, err = database.NewPostgresDatabase(context.Background(), &database.PostgreConfig{
			DSN: config.GetDSN(),
		})
		if err == nil {
			break
		}
		t.Logf("Waiting for database to be ready... attempt %d", i+1)
		time.Sleep(2 * time.Second)
	}
	require.NoError(t, err)

	sourceUrl := "file://" + ProjectRoot() + "/migrations"

	err = migrations.RunMigrationsWithSourceUrl(config.GetDSN(), sourceUrl)
	require.NoError(t, err)

	return db
}

func CleanupTestDB(t *testing.T, db database.Database) {
	_, err := db.GetPool().Exec(context.Background(), "TRUNCATE TABLE smart_models, smart_features CASCADE")
	require.NoError(t, err)
	db.Close()
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
