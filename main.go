package main

import (
	"clamp-core/config"
	"clamp-core/handlers"
	"clamp-core/listeners"
	"clamp-core/migrations"
	"clamp-core/models"
	"clamp-core/repository"
	"clamp-core/services"

	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Loading config...")
	err := config.Load()
	if err != nil {
		log.Fatalf("Loading config failed: %s", err)
	}

	logLevel, err := log.ParseLevel(config.ENV.LogLevel)
	if err != nil {
		log.Fatalf("Parsing log level failed: %s", err)
	}

	log.SetLevel(logLevel)

	log.Info("Initializing DB...")
	err = repository.InitDB()
	if err != nil {
		log.Fatalf("Initialzing DB failed: %s", err)
	}

	log.Info("Pinging DB...")
	err = repository.GetDB().Ping()
	if err != nil {
		log.Fatalf("DB ping failed: %s", err)
	}

	log.Info("Running DB migrations...")
	var cliArgs models.CLIArguments = os.Args[1:]
	os.Setenv("PORT", config.ENV.PORT)
	migrations.Migrate()

	if cliArgs.Parse().Find("migrate-only", "no") == "yes" {
		os.Exit(0)
	}

	log.Info("Initializing service request workers...")
	err = services.InitServiceRequestWorkers()
	if err != nil {
		log.Fatalf("Initializinng service request workers failed: %s", err)
	}

	log.Info("Initializing resume workers...")
	err = services.InitResumeWorkers()
	if err != nil {
		log.Fatalf("Initializinng resume workers failed: %s", err)
	}

	if config.ENV.EnableRabbitMQIntegration {
		listeners.AMQPStepResponseListener.Listen()
	}
	if config.ENV.EnableKafkaIntegration {
		listeners.KafkaStepResponseListener.Listen()
	}

	handlers.LoadHTTPRoutes()
	log.Info("Calling listener")
}
