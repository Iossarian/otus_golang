package storage

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

const DateLayout = "2006-01-02"
const DayDuration = "day"
const WeekDuration = "week"
const MonthDuration = "month"

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
}

func NewEvent() *Event {
	event := new(Event)
	event.ID = uuid.New()

	return event
}
