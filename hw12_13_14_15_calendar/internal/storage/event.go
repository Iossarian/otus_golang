package storage

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	ErrDateAlreadyBusy = errors.New("date already busy")
	ErrEventNotFound   = errors.New("event not exists")
)

type Event struct {
	ID          uuid.UUID `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	UserID      int       `db:"user_id"`
	StartDate   time.Time `db:"start_date"`
	EndDate     time.Time `db:"end_date"`
}
