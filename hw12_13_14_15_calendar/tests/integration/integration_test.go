package integration

import (
	"context"
	"net"
	"os"
	"testing"
	"time"

	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/storage"
	faker "github.com/bxcodec/faker/v3"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	fullDateLayout = "2006-01-02 15:04"
	dateLayout     = "2006-01-02"
)

type CalendarSuite struct {
	suite.Suite

	ctx        context.Context
	conn       *grpc.ClientConn
	grpcClient pb.CalendarClient
	db         *sqlx.DB
}

func (s *CalendarSuite) SetupSuite() {
	var err error

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}
	host := net.JoinHostPort("calendar_api", port)
	s.conn, err = grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.NoError(err)

	s.grpcClient = pb.NewCalendarClient(s.conn)

	dbAddr := "user=calendar password=calendar dbname=calendar sslmode=disable sslmode=disable host=calendar_db port=5432"
	s.db, err = sqlx.Connect("postgres", dbAddr)
	s.NoError(err)
	s.ctx = context.Background()
}

func (s *CalendarSuite) TearDownSuite() {
	s.NoError(s.db.Close())
	s.NoError(s.conn.Close())
}

func TestCalendarSuite(t *testing.T) {
	suite.Run(t, new(CalendarSuite))
}

func (s *CalendarSuite) TestCreateEvent() {
	s.Run("success", func() {
		event := createPbEvent()
		resp, err := s.grpcClient.Create(s.ctx, &pb.CreateRequest{Event: event})

		s.NoError(err)
		e, err := s.getEvent(resp.GetID())
		s.NoError(err)
		s.compareEvents(event, e)

		_, err = s.grpcClient.Delete(s.ctx, &pb.DeleteRequest{ID: resp.GetID()})
		s.NoError(err)
	})

	s.Run("create for busy date", func() {
		event := createPbEvent()
		startDate := event.StartDate
		resp, err := s.grpcClient.Create(s.ctx, &pb.CreateRequest{Event: event})

		s.NoError(err)
		newEvent := createPbEvent()
		newEvent.StartDate = startDate

		_, err = s.grpcClient.Create(s.ctx, &pb.CreateRequest{Event: event})
		s.Error(err)

		_, err = s.grpcClient.Delete(s.ctx, &pb.DeleteRequest{ID: resp.GetID()})
		s.NoError(err)
	})
}

func (s *CalendarSuite) TestEditEvent() {
	s.Run("success", func() {
		event := createPbEvent()
		resp, err := s.grpcClient.Create(s.ctx, &pb.CreateRequest{Event: event})
		s.NoError(err)

		event.ID = resp.ID
		event.Title = faker.TitleMale()
		event.Description = faker.Word()
		editedEvent, err := s.grpcClient.Edit(s.ctx, &pb.EditRequest{Event: event})
		s.NoError(err)

		s.compareEvents(editedEvent.Event, convertToEvent(event))

		_, err = s.grpcClient.Delete(s.ctx, &pb.DeleteRequest{ID: resp.GetID()})
		s.NoError(err)
	})
}

func (s *CalendarSuite) TestDeleteEvent() {
	s.Run("success", func() {
		event := createPbEvent()
		resp, err := s.grpcClient.Create(s.ctx, &pb.CreateRequest{Event: event})
		s.NoError(err)

		_, err = s.grpcClient.Delete(s.ctx, &pb.DeleteRequest{ID: resp.GetID()})
		s.NoError(err)
	})
}

func (s *CalendarSuite) TestList() {
	s.Run("for day", func() {
		event := createPbEvent()
		createResp, err := s.grpcClient.Create(s.ctx, &pb.CreateRequest{Event: event})
		s.NoError(err)

		eventDate := event.StartDate.AsTime().Add(-1 * time.Hour).Format(dateLayout)
		s.NoError(err)
		resp, _ := s.grpcClient.List(s.ctx, &pb.ListRequest{Date: eventDate, Duration: storage.DayDuration})
		s.Len(resp.GetList(), 1)

		_, err = s.grpcClient.Delete(s.ctx, &pb.DeleteRequest{ID: createResp.GetID()})
		s.NoError(err)
	})

	s.Run("for week", func() {
		event := createPbEvent()
		createResp, err := s.grpcClient.Create(s.ctx, &pb.CreateRequest{Event: event})
		s.NoError(err)

		eventDate := event.StartDate.AsTime().Add(-1 * time.Hour).Format(dateLayout)
		s.NoError(err)
		resp, _ := s.grpcClient.List(s.ctx, &pb.ListRequest{Date: eventDate, Duration: storage.WeekDuration})
		s.Len(resp.GetList(), 1)

		_, err = s.grpcClient.Delete(s.ctx, &pb.DeleteRequest{ID: createResp.GetID()})
		s.NoError(err)
	})

	s.Run("for month", func() {
		event := createPbEvent()
		createResp, err := s.grpcClient.Create(s.ctx, &pb.CreateRequest{Event: event})
		s.NoError(err)

		eventDate := event.StartDate.AsTime().Add(-1 * time.Hour).Format(dateLayout)
		s.NoError(err)
		resp, _ := s.grpcClient.List(s.ctx, &pb.ListRequest{Date: eventDate, Duration: storage.MonthDuration})
		s.Len(resp.GetList(), 1)

		_, err = s.grpcClient.Delete(s.ctx, &pb.DeleteRequest{ID: createResp.GetID()})
		s.NoError(err)
	})
}

func createPbEvent() *pb.Event {
	return &pb.Event{
		Title:            faker.TitleMale(),
		Description:      faker.Word(),
		StartDate:        timestamppb.New(time.Now().AddDate(0, 0, 1)),
		EndDate:          timestamppb.New(time.Now().AddDate(1, 0, 1).Add(time.Hour * time.Duration(2))),
		NotificationDate: timestamppb.New(time.Now().AddDate(1, 0, 1).Add(time.Hour)),
		UserID:           faker.RandomUnixTime(),
	}
}

func (s *CalendarSuite) compareEvents(pbEvent *pb.Event, e *storage.Event) {
	s.Equal(pbEvent.Description, e.Description)
	s.Equal(pbEvent.Title, e.Title)
	s.Equal(pbEvent.StartDate.AsTime().Format(fullDateLayout), e.StartDate.Format(fullDateLayout))
	s.Equal(pbEvent.EndDate.AsTime().Format(fullDateLayout), e.EndDate.Format(fullDateLayout))
}

func convertToEvent(e *pb.Event) *storage.Event {
	event := &storage.Event{}
	event.UserID = e.GetUserID()
	event.Title = e.GetTitle()
	event.Description = e.GetDescription()
	event.StartDate = e.GetStartDate().AsTime()
	event.EndDate = e.GetEndDate().AsTime()
	event.NotifyDate = e.GetNotificationDate().AsTime()

	return event
}

func (s *CalendarSuite) getEvent(id string) (*storage.Event, error) {
	event := new(storage.Event)
	err := s.db.
		QueryRowx("SELECT * FROM events WHERE id = $1", id).
		StructScan(event)

	return event, err
}
