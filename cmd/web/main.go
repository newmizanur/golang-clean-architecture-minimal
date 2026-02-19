package main

import (
	"fmt"
	"golang-clean-architecture/internal/config"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	client := config.NewDatabase(viperConfig, log)
	defer client.Close()
	validate := config.NewValidator(viperConfig)
	app := config.NewEcho(viperConfig)

	config.Bootstrap(&config.BootstrapConfig{
		Client:   client,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   viperConfig,
	})

	webPort := viperConfig.GetInt("web.port")
	err := app.Start(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
