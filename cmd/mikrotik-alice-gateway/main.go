package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joeshaw/envdecode"
	_ "github.com/joho/godotenv/autoload"
	"github.com/oklog/run"
	zerolog "github.com/rs/zerolog/log"

	"mikrotik-alice-gateway/internal/device_provider"
	"mikrotik-alice-gateway/internal/device_provider/mikrotik"
	"mikrotik-alice-gateway/internal/device_provider/wrap_logger"
	"mikrotik-alice-gateway/internal/log"
	"mikrotik-alice-gateway/internal/notifier/alice"
	"mikrotik-alice-gateway/internal/services"
	"mikrotik-alice-gateway/internal/services/checker"
	"mikrotik-alice-gateway/internal/services/rest"
	"mikrotik-alice-gateway/internal/storage/sql"
)

type config struct {
	Logger   log.Config
	Rest     rest.Config
	Storage  sql.Config
	Checker  checker.Config
	Notifier alice.Config
}

const signalChLen = 10

func main() {
	var cfg config
	if err := envdecode.StrictDecode(&cfg); err != nil {
		zerolog.Fatal().Err(err).Msg("Cannot decode config envs")
	}

	logger, err := log.New(cfg.Logger)
	if err != nil {
		zerolog.Fatal().Err(err).Msg("Cannot init logger")
	}

	ctx, cancel := context.WithCancel(logger.WithContext(context.Background()))

	g := &run.Group{}
	{
		stop := make(chan os.Signal, signalChLen)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		g.Add(func() error {
			<-stop
			return nil
		}, func(error) {
			signal.Stop(stop)
			cancel()
			close(stop)
		})
	}

	orderRunner := services.OrderRunner{}

	storage := sql.New(cfg.Storage)
	if err := storage.Connect(ctx); err != nil {
		logger.Panic().Err(err).Msg("Failed connect to db")
	}
	defer func() {
		if err := storage.Disconnect(ctx); err != nil {
			logger.Error().Err(err).Msg("Failed disconnect from db")
		}
	}()
	notifier := alice.New(cfg.Notifier)
	checkerInstance := checker.New(cfg.Checker, storage, func(routerID, address, username, password string) device_provider.DeviceProvider {
		return wrap_logger.New(mikrotik.New(mikrotik.Config{Password: password, Address: address, Username: username}), routerID, storage)
	}, notifier)
	if err := orderRunner.SetupService(ctx, checkerInstance, "checker", g); err != nil {
		logger.Fatal().Err(err).Msg("Failed setup checker service")
	}
	restService, err := rest.New(ctx, cfg.Rest, logger.With().Str("role", "rest").Logger(), storage, checkerInstance)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed create rest service")
	}
	if err := orderRunner.SetupService(ctx, restService, "rest", g); err != nil {
		logger.Fatal().Err(err).Msg("Failed setup rest service")
	}

	logger.Info().Msg("Running the service...")
	if err := g.Run(); err != nil {
		logger.Fatal().Err(err).Msg("The service has been stopped with error")
	}
	logger.Info().Msg("The service is stopped")
}
