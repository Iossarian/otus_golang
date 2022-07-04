package memorystorage

import (
	"errors"
	memoryStorage "github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestStorage(t *testing.T) {
	storage := New()
	dentistEvent := initDentistEvent()
	eventID := dentistEvent.ID.String()
	err := storage.Create(dentistEvent)
	require.NoError(t, err)

	t.Run("create", func(t *testing.T) {
		require.Len(t, storage.events, 1)

		deliveryEvent := initDeliveryEvent()
		err := storage.Create(deliveryEvent)

		require.True(t, errors.Is(err, memoryStorage.ErrDateAlreadyBusy))
	})

	t.Run("edit", func(t *testing.T) {
		dentistEvent.Title = "Plans to visit my dentist"

		eventsForDayBeforeEdit, err := storage.SelectForTheDay(dentistEvent.StartDate.Add(-1 * time.Hour))
		require.NoError(t, err)
		require.Len(t, eventsForDayBeforeEdit, 1)

		firstDayEvent := eventsForDayBeforeEdit[eventID]
		require.NotEqual(t, firstDayEvent, dentistEvent)

		err = storage.Edit(dentistEvent.ID.String(), dentistEvent)
		require.NoError(t, err)

		eventsForDayAfterEdit, err := storage.SelectForTheDay(dentistEvent.StartDate.Add(-1 * time.Hour))
		require.NoError(t, err)

		updatedEvent := eventsForDayAfterEdit[eventID]
		require.Equal(t, dentistEvent, updatedEvent)
	})

	t.Run("select", func(t *testing.T) {
		emptyForTheDayEvents, err := storage.SelectForTheDay(dentistEvent.StartDate.AddDate(0, 0, -2))
		require.NoError(t, err)
		require.Len(t, emptyForTheDayEvents, 0)

		forTheDayEvents, err := storage.SelectForTheDay(dentistEvent.StartDate.Add(-1 * time.Hour))
		require.NoError(t, err)
		require.Len(t, forTheDayEvents, 1)
		require.Equal(t, dentistEvent, forTheDayEvents[eventID])

		emptyForTheWeekEvents, err := storage.SelectForTheWeek(dentistEvent.StartDate.AddDate(0, 0, -8))
		require.NoError(t, err)
		require.Len(t, emptyForTheWeekEvents, 0)

		forTheWeekEvents, err := storage.SelectForTheWeek(dentistEvent.StartDate.AddDate(0, 0, -1))
		require.NoError(t, err)
		require.Len(t, forTheWeekEvents, 1)
		require.Equal(t, dentistEvent, forTheWeekEvents[eventID])

		emptyForTheMonthEvents, err := storage.SelectForTheMonth(dentistEvent.StartDate.AddDate(0, -1, -1))
		require.NoError(t, err)
		require.Len(t, emptyForTheMonthEvents, 0)

		forTheMonthEvents, err := storage.SelectForTheMonth(dentistEvent.StartDate.AddDate(0, 0, -1))
		require.NoError(t, err)
		require.Len(t, forTheMonthEvents, 1)
		require.Equal(t, dentistEvent, forTheMonthEvents[eventID])
	})

	t.Run("delete", func(t *testing.T) {
		err = storage.Delete(dentistEvent.ID.String())
		require.NoError(t, err)

		eventsForTheDay, err := storage.SelectForTheDay(dentistEvent.StartDate.Add(-1 * time.Hour))
		require.NoError(t, err)
		require.Len(t, eventsForTheDay, 0)

		err = storage.Delete(uuid.New().String())
		require.True(t, errors.Is(err, memoryStorage.ErrEventNotFound))
	})
}

func initDeliveryEvent() memoryStorage.Event {
	return memoryStorage.Event{
		ID:          uuid.New(),
		Title:       "Courier delivery",
		Description: "",
		UserID:      1,
		StartDate:   time.Now().Add(time.Hour),
		EndDate:     time.Now().Add(time.Hour),
	}
}

func initDentistEvent() memoryStorage.Event {
	return memoryStorage.Event{
		ID:          uuid.New(),
		Title:       "Dentist visit",
		Description: "My Dentist visiting",
		UserID:      1,
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(time.Hour * time.Duration(2)),
	}
}
