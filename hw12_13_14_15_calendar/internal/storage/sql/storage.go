package sqlstorage

import (
	"context"
	"fmt"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/config"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

type Storage struct {
	dsn string
	db  *sqlx.DB
	ctx context.Context
}

func New(c config.Config) *Storage {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable sslmode=disable", c.DBUser, c.DBPassword, c.DBTable)

	return &Storage{
		dsn: dsn,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.Open("postgres", s.dsn)
	if err != nil {
		return err
	}

	s.db = db
	s.ctx = ctx

	return nil
}

func (s *Storage) Close() error {
	if err := s.db.Close(); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Create(e storage.Event) error {
	query := `
				INSERT INTO events 
					(id, user_id, title, description, start_date, end_date, notification_date)
				VALUES
					($1, $2, $3, $4, $5, $6, $7)
				;
	`

	_, err := s.db.ExecContext(
		s.ctx,
		query,
		e.ID,
		e.UserID,
		e.Title,
		e.Description,
		e.StartDate,
		e.EndDate,
		e.NotifyDate,
	)

	return err
}

func (s *Storage) Delete(id string) error {
	query := `
				DELETE
				FROM
					events
				WHERE
					id = $1
	`

	_, err := s.db.ExecContext(s.ctx, query, id)

	return err
}

func (s *Storage) Edit(eventID string, e storage.Event) error {
	query := `
				UPDATE 
					events 
				SET
					user_id = $2,
					title = $3,
					description = $4, 
					start_date = $5, 
					end_date = $6,
					notification_date = $7
				WHERE 
					id = $1

	`
	_, err := s.db.ExecContext(
		s.ctx,
		query,
		eventID,
		e.UserID,
		e.Title,
		e.Description,
		e.StartDate,
		e.EndDate,
		e.NotifyDate,
	)

	return err
}

func (s *Storage) List(date time.Time, duration string) (map[string]storage.Event, error) {
	switch duration {
	case storage.DayDuration:
		return s.SelectBetween(date, date.AddDate(0, 0, 1))
	case storage.WeekDuration:
		return s.SelectBetween(date, date.AddDate(0, 0, 7))
	case storage.MonthDuration:
		return s.SelectBetween(date, date.AddDate(0, 1, 0))
	default:
		return s.SelectBetween(date, date.AddDate(0, 0, 1))
	}
}

func (s *Storage) SelectBetween(startDate time.Time, endDate time.Time) (map[string]storage.Event, error) {
	sql := `
		SELECT
		 id,
		 title,
		 description,
		 start_date,
		 end_date,
		 user_id,
		 notification_date
		FROM
		  events
		WHERE
		  start_date BETWEEN $1 AND $2
		;
	`
	rows, err := s.db.QueryxContext(s.ctx, sql, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make(map[string]storage.Event, 0)
	for rows.Next() {
		var event storage.Event
		err := rows.StructScan(&event)
		if err != nil {
			return nil, err
		}
		events[event.ID.String()] = event
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Storage) GetByNotificationPeriod(startDate, endDate time.Time) (map[string]storage.Event, error) {
	sql := `
		SELECT
		 id,
		 title,
		 description,
		 start_date,
		 end_date,
		 user_id,
		 notification_date
		FROM
		  events
		WHERE
		  is_notified = 0
		AND
		  notification_date BETWEEN $1 AND $2
		;
	`
	rows, err := s.db.QueryxContext(s.ctx, sql, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make(map[string]storage.Event, 0)
	for rows.Next() {
		var event storage.Event
		err := rows.StructScan(&event)
		if err != nil {
			return nil, err
		}
		events[event.ID.String()] = event
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Storage) ChangeNotifyStatus(eventID string) error {
	query := `
				UPDATE 
					events 
				SET
					is_notified = 1
				WHERE 
					id = $1

	`
	_, err := s.db.ExecContext(
		s.ctx,
		query,
		eventID,
	)

	return err
}

func (s *Storage) DeleteOldNotifiedEvents() error {
	query := `
				DELETE FROM 
					events 
				WHERE 
					is_notified = 1
				AND 
				    end_date <= $1

	`
	fmt.Println(time.Now().AddDate(-1, 0, 0))
	_, err := s.db.ExecContext(
		s.ctx,
		query,
		time.Now().AddDate(-1, 0, 0),
	)

	return err
}
