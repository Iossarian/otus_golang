package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
}

type Storage interface {
	Add(s *storage.Event) error
	Get(id string) (*storage.Event, error)
	Delete(id string)
	All() map[string]*storage.Event
	Edit(id string, e *storage.Event) error
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, title, date string) error {
	return a.storage.Add(&storage.Event{Title: title, Date: date})
}

func (a *App) GetEvent(id string) *storage.Event {
	event, err := a.storage.Get(id)
	var eventNotFoundErr *storage.EventNotFoundErr
	if errors.As(err, &eventNotFoundErr) {
		fmt.Println(err)
		return nil
	}

	return event
}

func (a *App) DeleteEvent(id string) {
	a.storage.Delete(id)
}

func (a *App) GetEvents() map[string]*storage.Event {
	return a.storage.All()
}

func (a *App) EditEvent(id string, e *storage.Event) {
	err := a.storage.Edit(id, e)
	var eventNotFoundErr *storage.EventNotFoundErr
	if errors.As(err, &eventNotFoundErr) {
		fmt.Println(err)
	}
}
