package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"golang-clean-architecture/internal/config"
	"golang-clean-architecture/internal/delivery/messaging/sqs"
	"golang-clean-architecture/internal/repository"
	"golang-clean-architecture/internal/usecase"
)

func main() {
	viperConfig := config.NewViper()
	logger := config.NewLogger(viperConfig)
	logger.Info("Starting SQS worker service")

	db := config.NewDatabase(viperConfig, logger)
	validate := config.NewValidator(viperConfig)

	contactRepository := repository.NewContactRepository(db, logger)
	contactUseCase := usecase.NewContactUseCase(db, logger, validate, contactRepository)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	contactHandler := sqs.NewContactHandler(logger, contactUseCase)
	worker := sqs.NewWorker(logger, viperConfig, contactHandler)
	go worker.Start(ctx)

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGTERM)
	<-terminateSignals

	logger.Info("Stopping SQS worker service")
	cancel()
	worker.Wait()
}
