package app

import (
	"context"
	"time"

	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(err error)
	Error(err error)
}

type Event interface{}

type Storage interface {
	Connect(ctx context.Context) error
	Close() error
	Create(s storage.Event) (id string, err error)
	Delete(id string) error
	Edit(id string, e storage.Event) error
	List(date time.Time, duration string) (map[string]storage.Event, error)
	GetByNotificationPeriod(startDate, endDate time.Time) (map[string]storage.Event, error)
	ChangeNotifyStatus(eventID string) error
	DeleteOldNotifiedEvents() error
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(event storage.Event) (id string, err error) {
	return a.storage.Create(event)
}

func (a *App) DeleteEvent(id string) error {
	return a.storage.Delete(id)
}

func (a *App) EditEvent(id string, e storage.Event) error {
	return a.storage.Edit(id, e)
}

func (a *App) List(date time.Time, duration string) map[string]storage.Event {
	events, err := a.storage.List(date, duration)
	if err != nil {
		a.logger.Error(err)
	}

	return events
}
