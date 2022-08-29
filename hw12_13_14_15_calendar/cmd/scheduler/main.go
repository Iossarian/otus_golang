package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/app"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/config"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/logger"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/storage/factory"
	_ "github.com/lib/pq"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "./../../configs/scheduler_config.yml", "Path to configuration file")
}

func main() {
	flag.Parse()
	appConfig := config.New(configPath)
	logFile, _ := os.OpenFile(appConfig.Logger.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	defer logFile.Close()
	applogger := logger.New(appConfig, logFile)

	storageSource, err := factory.GetStorage(appConfig)
	if err != nil {
		applogger.Error(err)
	}
	err = storageSource.Connect(context.Background())
	if err != nil {
		applogger.Error(fmt.Errorf("fail open storage: %w", err))
	}
	defer func(storageSource app.Storage) {
		err := storageSource.Close()
		if err != nil {
			applogger.Error(err)
		}
	}(storageSource)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	queue, _ := rabbitmq.NewConnection(appConfig, applogger)
	appScheduler := scheduler.New(appConfig, storageSource, applogger, queue)
	go func() {
		applogger.Info(fmt.Errorf("%s", "sheduler is running..."))
		if err := appScheduler.Run(ctx); err != nil {
			applogger.Info(fmt.Errorf("listen: %w", err))
		}
	}()
	<-ctx.Done()

	if err := appScheduler.Shutdown(); err != nil {
		applogger.Error(fmt.Errorf("fail stop scheduler %w", err))
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer func() {
		cancel()
	}()
}
