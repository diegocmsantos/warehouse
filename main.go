package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/diegocmsantos/warehouse/presentation/whhttp"
	"go.uber.org/zap"
)

func main() {
	os.Exit(start())
}

func start() int {

	logEnv := getStringOrDefault("LOG_ENV", "production")
	logger, err := createLogger(logEnv)
	if err != nil {
		return 1
	}

	s, err := whhttp.New(whhttp.ServerOptions{
		Addr:        "8080",
		Host:        "localhost",
		Port:        8080,
		ReadTimeout: 5 * time.Second,
		Logger:      logger,
	})
	if err != nil {
		logger.Sugar().Errorf("could not initiate http server: %s", err)
		return 1
	}
	sigs := make(chan os.Signal, 1)
	stopped := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigs
		if err := s.Stop(); err != nil {
			logger.Error("Error stopping server", zap.Error(err))
		}
		stopped <- true
	}()

	if err := s.Start(); err != nil {
		logger.Error("Error starting server", zap.Error(err))
		return 1
	}

	<-stopped

	return 0
}

func createLogger(env string) (*zap.Logger, error) {
	switch env {
	case "production":
		return zap.NewProduction()
	case "development":
		return zap.NewDevelopment()
	default:
		return zap.NewNop(), nil
	}
}

func getStringOrDefault(name, defaultV string) string {
	v, ok := os.LookupEnv(name)
	if !ok {
		return defaultV
	}
	return v
}
