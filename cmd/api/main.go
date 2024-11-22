package main

import (
	"github.com/kelseyhightower/envconfig"
	"os"
	"smart-hub/config"
)

var (
	cfg config.Config
)

func configSetup() {
	err := envconfig.Process("", &cfg)
	if err != nil {
		os.Exit(1)
	}
}

func main() {

	// Setup configuration
	configSetup()
}
