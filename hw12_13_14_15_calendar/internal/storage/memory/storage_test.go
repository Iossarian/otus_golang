package memorystorage

import (
	"errors"
	"testing"
	"time"

	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	memoryStorage := New()
	dentistEvent := initDentistEvent()
	eventID := dentistEvent.ID.String()
	_, err := memoryStorage.Create(dentistEvent)
	require.NoError(t, err)

	t.Run("create", func(t *testing.T) {
		require.Len(t, memoryStorage.events, 1)

		deliveryEvent := initDeliveryEvent()
		_, err := memoryStorage.Create(deliveryEvent)

		require.True(t, errors.Is(err, storage.ErrDateAlreadyBusy))
	})

	t.Run("edit", func(t *testing.T) {
		dentistEvent.Title = "Plans to visit my dentist"

		eventsForDayBeforeEdit, err := memoryStorage.List(dentistEvent.StartDate.Add(-1*time.Hour), storage.DayDuration)
		require.NoError(t, err)
		require.Len(t, eventsForDayBeforeEdit, 1)

		firstDayEvent := eventsForDayBeforeEdit[eventID]
		require.NotEqual(t, firstDayEvent, dentistEvent)

		err = memoryStorage.Edit(dentistEvent.ID.String(), dentistEvent)
		require.NoError(t, err)

		eventsForDayAfterEdit, err := memoryStorage.List(dentistEvent.StartDate.Add(-1*time.Hour), storage.DayDuration)
		require.NoError(t, err)

		updatedEvent := eventsForDayAfterEdit[eventID]
		require.Equal(t, dentistEvent, updatedEvent)
	})

	t.Run("select", func(t *testing.T) {
		emptyForTheDayEvents, err := memoryStorage.List(dentistEvent.StartDate.AddDate(0, 0, -2), storage.DayDuration)
		require.NoError(t, err)
		require.Len(t, emptyForTheDayEvents, 0)

		forTheDayEvents, err := memoryStorage.List(dentistEvent.StartDate.Add(-1*time.Hour), storage.DayDuration)
		require.NoError(t, err)
		require.Len(t, forTheDayEvents, 1)
		require.Equal(t, dentistEvent, forTheDayEvents[eventID])

		emptyForTheWeekEvents, err := memoryStorage.List(dentistEvent.StartDate.AddDate(0, 0, -8), storage.WeekDuration)
		require.NoError(t, err)
		require.Len(t, emptyForTheWeekEvents, 0)

		forTheWeekEvents, err := memoryStorage.List(dentistEvent.StartDate.AddDate(0, 0, -1), storage.WeekDuration)
		require.NoError(t, err)
		require.Len(t, forTheWeekEvents, 1)
		require.Equal(t, dentistEvent, forTheWeekEvents[eventID])

		emptyForTheMonthEvents, err := memoryStorage.List(dentistEvent.StartDate.AddDate(0, -1, -1), storage.MonthDuration)
		require.NoError(t, err)
		require.Len(t, emptyForTheMonthEvents, 0)

		forTheMonthEvents, err := memoryStorage.List(dentistEvent.StartDate.AddDate(0, 0, -1), storage.MonthDuration)
		require.NoError(t, err)
		require.Len(t, forTheMonthEvents, 1)
		require.Equal(t, dentistEvent, forTheMonthEvents[eventID])
	})

	t.Run("delete", func(t *testing.T) {
		err = memoryStorage.Delete(dentistEvent.ID.String())
		require.NoError(t, err)

		eventsForTheDay, err := memoryStorage.List(dentistEvent.StartDate.Add(-1*time.Hour), storage.DayDuration)
		require.NoError(t, err)
		require.Len(t, eventsForTheDay, 0)

		err = memoryStorage.Delete(uuid.New().String())
		require.True(t, errors.Is(err, storage.ErrEventNotFound))
	})
}

func initDeliveryEvent() storage.Event {
	event := storage.NewEvent()
	event.Title = "Courier delivery"
	event.Description = ""
	event.UserID = 1
	event.StartDate = time.Now().Add(time.Hour)
	event.EndDate = time.Now().Add(time.Hour)

	return *event
}

func initDentistEvent() storage.Event {
	event := storage.NewEvent()
	event.Title = "Dentist visit"
	event.Description = "My Dentist visiting"
	event.UserID = 1
	event.StartDate = time.Now()
	event.EndDate = time.Now().Add(time.Hour * time.Duration(2))

	return *event
}
