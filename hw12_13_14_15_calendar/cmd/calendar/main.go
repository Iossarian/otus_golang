package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/app"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/config"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/logger"
	internalHttp "github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/server/http"
	memoryStorage "github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/storage/memory"
	sqlStorage "github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/storage/sql"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

var configPath string

const memoryStorageName = "memory"
const sqlStorageName = "sql"

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
	logg := logger.New(appConfig)

	storageSource := getStorage(appConfig)
	err := storageSource.Connect(context.Background())
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

	server := internalHttp.NewServer(logg, calendar, appConfig)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		logg.Info(fmt.Errorf("%s", "calendar is running..."))
		if err := server.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logg.Info(fmt.Errorf("listen: %w", err))
		}
	}()

	<-ctx.Done()

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer func() {
		cancel()
	}()

	if err := server.Stop(ctx); err != nil {
		logg.Error(err)
	}
}

func getStorage(c config.Config) app.Storage {
	switch c.StorageSource {
	case memoryStorageName:
		return memoryStorage.New()
	case sqlStorageName:
		return sqlStorage.New(c)
	default:
		return sqlStorage.New(c)
	}
}
