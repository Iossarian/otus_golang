package storage

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	DateLayout    = "2006-01-02"
	DayDuration   = "day"
	WeekDuration  = "week"
	MonthDuration = "month"
)

var (
	ErrDateAlreadyBusy = errors.New("date already busy")
	ErrEventNotFound   = errors.New("event not exists")
)

type EventDate struct {
	time.Time
}

type Event struct {
	ID          uuid.UUID `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	UserID      int64     `db:"user_id"`
	StartDate   time.Time `db:"start_date"`
	EndDate     time.Time `db:"end_date"`
	NotifyDate  time.Time `db:"notification_date"`
	IsNotified  int8      `db:"is_notified"`
}

func NewEvent() *Event {
	event := new(Event)
	event.ID = uuid.New()
	event.IsNotified = 0

	return event
}
