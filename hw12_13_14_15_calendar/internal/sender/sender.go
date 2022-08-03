package sender

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/app"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/storage"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Sender struct {
	queue   Queue
	storage app.Storage
	logger  app.Logger
}

type Queue interface {
	Consume(ctx context.Context, handleFunc rabbitmq.HandleFunc) error
	Close() error
}

func New(storage app.Storage, logger app.Logger, queue Queue) *Sender {
	return &Sender{
		storage: storage,
		logger:  logger,
		queue:   queue,
	}
}

func (s *Sender) Run(ctx context.Context) error {
	s.logger.Info(fmt.Errorf("sender "))

	return s.queue.Consume(ctx, s.eventHandle)
}

func (s *Sender) eventHandle(ctx context.Context, msgs <-chan amqp.Delivery) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-msgs:
			if !ok {
				return
			}

			s.logger.Info(fmt.Errorf("received message from a queue"))

			event := &storage.Event{}

			if err := json.Unmarshal(msg.Body, event); err != nil {
				s.logger.Error(fmt.Errorf("event unmarshal failed"))
				return
			}

			if err := msg.Ack(false); err != nil {
				s.logger.Error(fmt.Errorf("ack failed"))
			}

			if err := s.storage.ChangeNotifyStatus(event.ID.String()); err != nil {
				s.logger.Error(fmt.Errorf("notify failed"))
			}
		}
	}
}

func (s *Sender) Shutdown() error {
	s.logger.Info(fmt.Errorf("sender is shutting down"))

	return s.queue.Close()
}
