package grpc

import (
	context "context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	app    Application
	logger Logger
	config Config
	server *grpc.Server
}

type Logger interface {
	Info(err error)
	Error(err error)
	LogHTTPRequest(r *http.Request, statusCode int, requestDuration time.Duration)
	LogGRPCRequest(code codes.Code, method, address string, requestDuration time.Duration)
}

type Application interface {
	CreateEvent(event storage.Event) (id string, err error)
	EditEvent(id string, e storage.Event) error
	DeleteEvent(id string) error
	List(date time.Time, duration string) map[string]storage.Event
}

type Config interface {
	GetHTTPAddr() string
	GetGRPCAddr() string
}

type CalendarService struct {
	app    Application
	logger Logger
	pb.UnimplementedCalendarServer
}

func NewServer(logger Logger, app Application, config Config) *Server {
	chainInterceptor := grpc.ChainUnaryInterceptor(
		LoggingInterceptor(logger),
	)
	grpcServer := grpc.NewServer(chainInterceptor)

	service := &CalendarService{
		app:    app,
		logger: logger,
	}
	service.app = app
	pb.RegisterCalendarServer(grpcServer, service)

	srv := &Server{
		logger: logger,
		app:    app,
		config: config,
		server: grpcServer,
	}

	return srv
}

func (s *Server) Start(ctx context.Context) error {
	lsn, err := net.Listen("tcp", s.config.GetGRPCAddr())
	if err != nil {
		s.logger.Error(fmt.Errorf("fail start gprc server: %w", err))
	}

	reflection.Register(s.server)
	if err := s.server.Serve(lsn); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error(fmt.Errorf("listen: %w", err))
		os.Exit(1)
	}
	<-ctx.Done()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info(fmt.Errorf("calendar is shutting down"))
	s.server.GracefulStop()
	<-ctx.Done()

	return nil
}

func (s *CalendarService) Create(_ context.Context, in *pb.CreateRequest) (*pb.CreateResponse, error) {
	id, err := s.app.CreateEvent(convertToEvent(in.GetEvent()))
	if err != nil {
		s.logger.Error(fmt.Errorf("fail create event %w ", err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.CreateResponse{ID: id}, nil
}

func (s *CalendarService) Edit(_ context.Context, in *pb.EditRequest) (*pb.EditResponse, error) {
	pbEvent := in.GetEvent()
	eventUUID, _ := uuid.FromBytes([]byte(pbEvent.GetID()))
	event := storage.Event{
		ID:          eventUUID,
		Title:       pbEvent.GetTitle(),
		Description: pbEvent.GetDescription(),
		StartDate:   pbEvent.GetStartDate().AsTime(),
		EndDate:     pbEvent.GetEndDate().AsTime(),
		NotifyDate:  pbEvent.GetNotificationDate().AsTime(),
		UserID:      pbEvent.GetUserID(),
	}

	if err := s.app.EditEvent(event.ID.String(), event); err != nil {
		s.logger.Error(fmt.Errorf("event edition fail: %w", err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.EditResponse{Event: convertToPbEvent(event)}, nil
}

func (s *CalendarService) Delete(_ context.Context, in *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	if err := s.app.DeleteEvent(in.ID); err != nil {
		s.logger.Error(fmt.Errorf("event deleting fail: %w", err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.DeleteResponse{}, nil
}

func (s *CalendarService) List(_ context.Context, in *pb.ListRequest) (*pb.ListResponse, error) {
	parsedDate, _ := time.Parse(storage.DateLayout, in.Date)
	list := s.app.List(parsedDate, in.Duration)

	result := make(map[string]*pb.Event)
	for id, event := range list {
		result[id] = &pb.Event{
			ID:          event.ID.String(),
			UserID:      event.UserID,
			Title:       event.Title,
			Description: event.Description,
			StartDate:   timestamppb.New(event.StartDate),
			EndDate:     timestamppb.New(event.EndDate),
		}
	}

	return &pb.ListResponse{List: result}, nil
}

func convertToPbEvent(e storage.Event) *pb.Event {
	return &pb.Event{
		ID:               e.ID.String(),
		Title:            e.Title,
		Description:      e.Description,
		StartDate:        timestamppb.New(e.StartDate),
		EndDate:          timestamppb.New(e.EndDate),
		NotificationDate: timestamppb.New(e.NotifyDate),
		UserID:           e.UserID,
	}
}

func convertToEvent(e *pb.Event) storage.Event {
	event := storage.NewEvent()
	event.UserID = e.GetUserID()
	event.Title = e.GetTitle()
	event.Description = e.GetDescription()
	event.StartDate = e.GetStartDate().AsTime()
	event.EndDate = e.GetEndDate().AsTime()
	event.NotifyDate = e.GetNotificationDate().AsTime()

	return *event
}
