package main

import (
	"context"
	"flag"
	memorystorage "github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/storage/memory"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/app"
	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/server/http"
)

var configFile string

const memoryStorage = "memory"
const sqlStorage = "sql"

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()
	if flag.Arg(0) == "version" {
		printVersion()
		return
	}
	config := NewConfig()
	logg := logger.New(config.Logger.Level)

	storageSource := getStorage(config)
	calendar := app.New(logg, storageSource)

	server := internalhttp.NewServer(logg, calendar, config)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

func getStorage(c Config) app.Storage {
	switch c.StorageSource {
	case memoryStorage:
		return memorystorage.New()
	/*case sqlStorage:
	return sqlstorage.New()*/
	default:
		return memorystorage.New()
	}
}
