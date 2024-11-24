package main

import (
	"context"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"smart-hub/config"
	pbModel "smart-hub/gen/proto/smart_model/v1"
	"smart-hub/internal/application/service"
	"smart-hub/internal/common/database"
	"smart-hub/internal/common/database/migrations"
	"smart-hub/internal/common/logger"
	"smart-hub/internal/infrastructure/database/postgres"
	"smart-hub/internal/presentation/grpc/handler"
	"smart-hub/internal/presentation/grpc/mapper"
	"syscall"
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

	// Initialize repositories
	smartModelRepo := postgres.NewPGSmartModelRepository(db)

	// Initialize services
	smartModelService := service.NewSmartModelService(smartModelRepo)

	// Initialize mappers
	smartModelMapper := mapper.NewSmartModelMapper()

	// Initialize handlers
	smartModelHandler := handler.NewSmartModelHandler(smartModelService, smartModelMapper)

	// Initialize GRPC server
	grpcServer := grpc.NewServer()

	pbModel.RegisterSmartModelServiceServer(grpcServer, smartModelHandler)

	address := fmt.Sprintf(":%s", cfg.Service.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error("Failed to listen", err)
		os.Exit(1)
	}

	logger.Info(fmt.Sprintf("Starting GRPC server on %s", address))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			logger.Error("Failed to serve", err)
			os.Exit(1)
		}
	}()

	<-stop
	logger.Info("Shutting down server...")

	grpcServer.GracefulStop()
	logger.Info("Server stopped")

}
