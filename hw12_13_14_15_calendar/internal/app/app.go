package app

import (
	"context"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/storage"
	"time"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(err error)
	Error(err error)
}

type Storage interface {
	Connect(ctx context.Context) error
	Close() error
	Create(s storage.Event) error
	Delete(id string) error
	Edit(id string, e storage.Event) error
	SelectForTheDay(date time.Time) (map[string]storage.Event, error)
	SelectForTheWeek(date time.Time) (map[string]storage.Event, error)
	SelectForTheMonth(date time.Time) (map[string]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(event storage.Event) {
	if err := a.storage.Create(event); err != nil {
		a.logger.Error(err)
	}
}

func (a *App) DeleteEvent(id string) {
	if err := a.storage.Delete(id); err != nil {
		a.logger.Error(err)
	}
}

func (a *App) EditEvent(id string, e storage.Event) {
	if err := a.storage.Edit(id, e); err != nil {
		a.logger.Error(err)
	}
}

func (a *App) SelectForTheDay(date time.Time) map[string]storage.Event {
	events, err := a.storage.SelectForTheDay(date)
	if err != nil {
		a.logger.Error(err)
	}

	return events
}
func (a *App) SelectForTheWeek(date time.Time) map[string]storage.Event {
	events, err := a.storage.SelectForTheWeek(date)
	if err != nil {
		a.logger.Error(err)
	}

	return events
}

func (a *App) SelectForTheMonth(date time.Time) map[string]storage.Event {
	events, err := a.storage.SelectForTheMonth(date)
	if err != nil {
		a.logger.Error(err)
	}

	return events
}
