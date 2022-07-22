package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/app"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/config"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/logger"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/server/http"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/storage/factory"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "./../../configs/config.env", "Path to configuration file")
}

func main() {
	flag.Parse()
	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	appConfig := config.NewConfig(configPath)
	logFile, _ := os.OpenFile(appConfig.LoggingFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer logFile.Close()
	logg := logger.New(appConfig, logFile)

	storageSource, err := factory.GetStorage(appConfig)
	if err != nil {
		logg.Error(err)
	}
	err = storageSource.Connect(context.Background())
	if err != nil {
		logg.Error(fmt.Errorf("fail open storage: %w", err))
	}
	defer func(storageSource app.Storage) {
		err := storageSource.Close()
		if err != nil {
			logg.Error(err)
		}
	}(storageSource)

	calendar := app.New(logg, storageSource)

	httpServer := internalhttp.NewServer(logg, calendar, appConfig)
	grpcServer := grpc.NewServer(logg, calendar, appConfig)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		logg.Info(fmt.Errorf("%s", "http server is running..."))
		if err := httpServer.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logg.Info(fmt.Errorf("listen: %w", err))
		}
	}()

	go func() {
		logg.Info(fmt.Errorf("%s", "gRPC server is running..."))
		if err := grpcServer.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logg.Info(fmt.Errorf("listen: %w", err))
		}
	}()

	<-ctx.Done()

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer func() {
		cancel()
	}()

	if err := httpServer.Stop(ctx); err != nil {
		logg.Error(err)
	}
	if err := grpcServer.Stop(ctx); err != nil {
		logg.Error(err)
	}
}
