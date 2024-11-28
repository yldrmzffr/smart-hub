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
	pbHealth "smart-hub/gen/proto/health/v1"
	pbFeature "smart-hub/gen/proto/smart_feature/v1"
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

type App struct {
	cfg        config.Config
	grpcServer *grpc.Server
	db         database.Database
}

func NewApp() *App {
	return &App{
		grpcServer: grpc.NewServer(),
	}
}

func (a *App) configSetup() error {
	return envconfig.Process("", &a.cfg)
}

func (a *App) loggerSetup() error {
	return logger.InitLogger(&a.cfg.Log)
}

func (a *App) databaseSetup(ctx context.Context) error {
	// Migrate database
	if err := migrations.RunMigrations(a.cfg.Database.GetURL()); err != nil {
		return fmt.Errorf("database migrations error: %w", err)
	}

	// Connect to database
	db, err := database.NewPostgresDatabase(ctx, &database.PostgreConfig{DSN: a.cfg.Database.GetDSN()})
	if err != nil {
		return fmt.Errorf("database connection error: %w", err)
	}
	a.db = db
	return nil
}

func (a *App) smartFeatureSetup() {
	smartFeatureRepo := postgres.NewPGSmartFeatureRepository(a.db)
	smartFeatureService := service.NewSmartFeatureService(smartFeatureRepo)
	smartFeatureMapper := mapper.NewSmartFeatureMapper()
	smartFeatureHandler := handler.NewSmartFeatureHandler(smartFeatureService, smartFeatureMapper)
	pbFeature.RegisterSmartFeatureServiceServer(a.grpcServer, smartFeatureHandler)
}

func (a *App) smartModelSetup() {
	smartModelRepo := postgres.NewPGSmartModelRepository(a.db)
	smartModelService := service.NewSmartModelService(smartModelRepo)
	smartModelMapper := mapper.NewSmartModelMapper()
	smartModelHandler := handler.NewSmartModelHandler(smartModelService, smartModelMapper)
	pbModel.RegisterSmartModelServiceServer(a.grpcServer, smartModelHandler)
}

func (a *App) healthSetup() {
	healthHandler := handler.NewHealthHandler(a.db)
	pbHealth.RegisterHealthServer(a.grpcServer, healthHandler)
}

func (a *App) shutdown() {
	logger.Info("Shutting down server...")
	a.grpcServer.GracefulStop()
	if a.db != nil {
		a.db.Close()
	}
	logger.Info("Server stopped")
}

func main() {
	ctx := context.Background()
	app := NewApp()
	defer app.shutdown()

	// Initialize base components
	if err := app.configSetup(); err != nil {
		logger.Error("Config setup error", err)
		os.Exit(1)
	}

	if err := app.loggerSetup(); err != nil {
		logger.Error("Logger setup error", err)
		os.Exit(1)
	}

	if err := app.databaseSetup(ctx); err != nil {
		logger.Error("Database setup error", err)
		os.Exit(1)
	}

	// Initialize modules
	app.healthSetup()
	app.smartModelSetup()
	app.smartFeatureSetup()

	// Start server
	address := fmt.Sprintf(":%s", app.cfg.Service.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error("Failed to listen", err)
		os.Exit(1)
	}

	logger.Info(fmt.Sprintf("Starting GRPC server on %s", address))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.grpcServer.Serve(listener); err != nil {
			logger.Error("Failed to serve", err)
			os.Exit(1)
		}
	}()

	<-stop
}
