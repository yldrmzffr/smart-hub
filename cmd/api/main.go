package main

import (
	"github.com/kelseyhightower/envconfig"
	"os"
	"smart-hub/config"
	"smart-hub/internal/common/logger"
)

var (
	cfg config.Config
)

func configSetup() {
	err := envconfig.Process("", &cfg)
	if err != nil {
		logger.Error("Config setup error", err)
		os.Exit(1)
	}
}

func loggerSetup() {
	err := logger.InitLogger(&cfg.Log)
	if err != nil {
		panic(err)
	}
}

func main() {

	// Setup configuration
	configSetup()

	// Initialize logger
	loggerSetup()

	logger.Info("Service started", "service", cfg.Service.Name, "port", cfg.Service.Port)

}
