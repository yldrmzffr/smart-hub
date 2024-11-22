package main

import (
	"context"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"os"
	"smart-hub/config"
	"smart-hub/internal/common/database"
	"smart-hub/internal/common/database/migrations"
	"smart-hub/internal/common/logger"
)

var (
	cfg config.Config
)

func configSetup() {
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(fmt.Errorf("config setup error: %w", err))
	}
}

func loggerSetup() {
	err := logger.InitLogger(&cfg.Log)
	if err != nil {
		panic(fmt.Errorf("logger setup error: %w", err))
	}
}

func migrateDatabase(cfg *config.DatabaseConfig) {
	dbUrl := cfg.GetURL()
	if err := migrations.RunMigrations(dbUrl); err != nil {
		panic(fmt.Errorf("database migrations error: %w", err))
	}
}

func panicError() {
	if r := recover(); r != nil {
		logger.Error("Panic error", "error", r)
		os.Exit(1)
	}
}

func main() {
	ctx := context.Background()

	defer panicError()

	// Setup configuration
	configSetup()

	// Initialize logger
	loggerSetup()

	// Migrate database
	migrateDatabase(&cfg.Database)

	// Initialize database connection
	db, err := database.NewPostgresDatabase(ctx, &database.PostgreConfig{DSN: cfg.Database.GetDSN()})
	if err != nil {
		logger.Error("Database Connection Error", err)
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("Service started", "service", cfg.Service.Name, "port", cfg.Service.Port)

}
