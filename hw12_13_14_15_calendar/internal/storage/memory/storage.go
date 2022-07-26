package memorystorage

import (
	"context"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/storage"
	"sync"
	"time"
)

type Storage struct {
	events map[string]storage.Event
	mu     sync.RWMutex
}

var events = map[string]storage.Event{}

func New() *Storage {
	events = make(map[string]storage.Event)
	str := &Storage{events: events}

	return str
}

func (s *Storage) Create(e storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, event := range s.events {
		if event.StartDate.Before(e.StartDate) && event.EndDate.After(e.StartDate) && event.UserID == e.UserID {
			return storage.ErrDateAlreadyBusy
		}
	}

	s.events[e.ID.String()] = e

	return nil
}

func (s *Storage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := events[id]; !ok {
		return storage.ErrEventNotFound
	}

	delete(s.events, id)

	return nil
}

func (s *Storage) Edit(id string, e storage.Event) error {
	s.mu.RLock()
	if _, ok := events[id]; !ok {
		return storage.ErrEventNotFound
	}
	s.mu.RUnlock()

	s.mu.Lock()
	events[id] = e
	s.mu.Unlock()

	return nil
}

func (s *Storage) List(date time.Time, duration string) (map[string]storage.Event, error) {
	switch duration {
	case storage.DayDuration:
		return s.list(date, date.AddDate(0, 0, 1))
	case storage.WeekDuration:
		return s.list(date, date.AddDate(0, 0, 7))
	case storage.MonthDuration:
		return s.list(date, date.AddDate(0, 1, 0))
	default:
		return s.list(date, date.AddDate(0, 0, 1))
	}
}

func (s *Storage) list(startDate, endDate time.Time) (map[string]storage.Event, error) {
	list := make(map[string]storage.Event, 0)

	for id, event := range s.events {
		if event.StartDate.After(startDate) && event.StartDate.Before(endDate) {
			list[id] = event
		}
	}

	return list, nil
}

func (s *Storage) GetByNotificationPeriod(startDate, endDate time.Time) (map[string]storage.Event, error) {
	events, _ := s.list(startDate, endDate)

	resultMap := make(map[string]storage.Event, 0)
	for id, event := range events {
		if event.IsNotified == 0 {
			resultMap[id] = event
		}
	}

	return resultMap, nil
}

func (s *Storage) ChangeNotifyStatus(eventID string) error {
	for _, event := range s.events {
		if event.ID.String() == eventID {
			event.IsNotified = 1
		}
	}

	return nil
}

func (s *Storage) Connect(_ context.Context) error {
	return nil
}

func (s *Storage) Close() error {
	return nil
}

func (s *Storage) DeleteOldNotifiedEvents() error {
	for id, event := range s.events {
		if event.EndDate.After(time.Now().AddDate(-1, 0, 0)) {
			delete(s.events, id)
		}
	}

	return nil
}
