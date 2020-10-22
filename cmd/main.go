package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	service "github.com/griddis/atlant_test/cmd/service"
	"github.com/griddis/atlant_test/configs"
	"github.com/griddis/atlant_test/internal/server"
	"github.com/griddis/atlant_test/pkg/health"
	"github.com/griddis/atlant_test/pkg/repository"
	"github.com/griddis/atlant_test/tools/logging"
)

func main() {

	cfg := initConfig()
	if err := cfg.Print(); err != nil {
		fmt.Fprintf(os.Stderr, "print config: %s", err)
		os.Exit(1)
	}

	var (
		logger = logging.NewLogger(cfg.Logger.Level, cfg.Logger.TimeFormat)
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = logging.WithContext(ctx, logger)

	healthService := initHealthService(ctx, cfg)

	rctx, _ := context.WithTimeout(ctx, 10*time.Second)
	repository, err := repository.NewRepo(rctx, logger, &cfg.Database)

	if err != nil {
		logger.Error("init", "repository", "err", err)
		os.Exit(1)
	}
	defer repository.Close(rctx)

	httpClientTransport := &http.Client{}
	mainService := service.NewService(ctx, cfg, repository, httpClientTransport)

	s, err := server.NewServer(
		server.SetConfig(&cfg.Server),
		server.SetLogger(logger),
		server.SetHandler(
			map[string]http.Handler{
				"_": service.MakeHTTPHandler(ctx, mainService),
				"":  health.MakeHTTPHandler(ctx, healthService),
			}),
		server.SetGRPC(
			service.JoinGRPC(ctx, mainService),
			health.JoinGRPC(ctx, healthService),
		),
	)
	if err != nil {
		logger.Error("init", "server", "err", err)
		os.Exit(1)
	}
	defer s.Close()

	if err := s.AddHTTP(); err != nil {
		logger.Error("err", err)
		os.Exit(1)
	}

	if err = s.AddGRPC(); err != nil {
		logger.Error("err", err)
		os.Exit(1)
	}

	s.AddSignalHandler()
	s.Run()

}

func initConfig() *configs.Config {
	cfg := configs.NewConfig()
	if err := cfg.Read(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "read config: %s", err)
		os.Exit(1)
	}
	return cfg
}

func initHealthService(ctx context.Context, cfg *configs.Config) health.Service {
	healthService := health.NewService()
	healthService = health.NewLoggingService(ctx, healthService)
	return healthService
}
