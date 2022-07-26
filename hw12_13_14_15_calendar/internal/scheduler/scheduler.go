package scheduler

import (
	"context"
	"fmt"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/app"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/config"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/storage"
	"time"
)

type Queue interface {
	Publish(data interface{}) error
	Close() error
}

type Scheduler struct {
	timeout time.Duration
	storage app.Storage
	logger  app.Logger
	queue   Queue
}

func New(config config.Config, storage app.Storage, logger app.Logger, queue Queue) *Scheduler {
	return &Scheduler{
		timeout: config.EventScanTimeout,
		storage: storage,
		logger:  logger,
		queue:   queue,
	}
}

func (s *Scheduler) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.timeout)
	defer ticker.Stop()

	for {
		go func() {
			s.notify()
			err := s.storage.DeleteOldNotifiedEvents()
			if err != nil {
				s.logger.Error(fmt.Errorf("fail delete old events %w", err))
			}
		}()

		select {
		case <-ctx.Done():
			s.logger.Info(fmt.Errorf("stop scheduler"))

			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func (s *Scheduler) notify() {
	startDate := time.Now().Add(-s.timeout + time.Second)
	endDate := time.Now()

	events, err := s.storage.GetByNotificationPeriod(startDate, endDate)
	if err != nil {
		s.logger.Error(fmt.Errorf("fail get event for notify %w", err))
	}

	s.logger.Info(fmt.Errorf("found %d events to notify", len(events)))

	newEvent := storage.NewEvent()
	newEvent.Title = "New event"
	err = s.queue.Publish(newEvent)
	if err != nil {
		println(err)
	}

	for _, event := range events {
		if err := s.queue.Publish(event); err != nil {
			s.logger.Error(fmt.Errorf("fail publish event message"))
		}
	}
}

func (s *Scheduler) Shutdown() error {
	s.logger.Info(fmt.Errorf("scheduler is shutting down"))

	return s.queue.Close()
}
