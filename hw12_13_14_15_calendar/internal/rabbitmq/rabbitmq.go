package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/app"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Rabbit struct {
	addr          string
	queueName     string
	handlersCount int
	conn          *amqp.Connection
	channel       *amqp.Channel
	logger        app.Logger
}

type HandleFunc func(context.Context, <-chan amqp.Delivery)

func New(config *config.Config, logger app.Logger) *Rabbit {
	return &Rabbit{
		addr:          config.AMQP.Addr,
		queueName:     config.AMQP.Name,
		handlersCount: config.AMQP.HandlersCount,
		logger:        logger,
	}
}

func NewConnection(config *config.Config, logger app.Logger) (*Rabbit, error) {
	rabbit := New(config, logger)

	var err error
	rabbit.conn, err = amqp.Dial(rabbit.addr)
	if err != nil {
		rabbit.logger.Error(fmt.Errorf("fail connect to rabbitmq %w", err))
		return nil, err
	}

	rabbit.channel, err = rabbit.conn.Channel()
	if err != nil {
		rabbit.logger.Error(fmt.Errorf("fail open channel"))
	}

	rabbit.logger.Info(fmt.Errorf("connected to rabbitmq"))

	return rabbit, nil
}

func (r *Rabbit) Close() error {
	if err := r.channel.Cancel("", true); err != nil {
		r.logger.Error(fmt.Errorf("fail cancel consumer %w", err))
	}
	if err := r.conn.Close(); err != nil {
		r.logger.Error(fmt.Errorf("fail close amqp connection %w", err))
	}

	return nil
}

func (r *Rabbit) Publish(data interface{}) error {
	encodedData, err := json.Marshal(data)
	if err != nil {
		r.logger.Error(fmt.Errorf("fail encode data"))
	}

	return r.channel.PublishWithContext(
		context.Background(),
		"",
		r.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        encodedData,
		},
	)
}

func (r *Rabbit) Consume(ctx context.Context, handleFunc HandleFunc) error {
	for {
		msgs, err := r.announceQueue()
		if err != nil {
			return err
		}

		for i := 0; i < r.handlersCount; i++ {
			go handleFunc(ctx, msgs)
		}
	}
}

func (r *Rabbit) announceQueue() (<-chan amqp.Delivery, error) {
	_, err := r.channel.QueueDeclare(r.queueName, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("queue declare failed: %w", err)
	}

	msgs, err := r.channel.Consume(
		r.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("consume failed: %w", err)
	}

	return msgs, nil
}
