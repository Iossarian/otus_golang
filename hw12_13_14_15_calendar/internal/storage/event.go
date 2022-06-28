package storage

import (
	"fmt"
	"github.com/google/uuid"
)

type EventNotFoundErr struct {
	Id string
}

func (e EventNotFoundErr) Error() string {
	return fmt.Sprintf("event %s not found", e.Id)
}

type EventAlreadyExistErr struct {
	Id string
}

func (e EventAlreadyExistErr) Error() string {
	return fmt.Sprintf("event with id %s not found", e.Id)
}

type Event struct {
	ID                                                       uuid.UUID
	Title, Date, Duration, Description, NotifyBefore, UserId string
}
