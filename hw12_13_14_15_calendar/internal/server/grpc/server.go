package grpc

import (
	context "context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"net"
	"net/http"
	"os"
	"time"
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
	CreateEvent(event storage.Event) error
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
	encodedEvent, _ := json.Marshal(in.Event)
	event := storage.NewEvent()
	event.UserID = in.Event.GetUserID()
	event.Title = in.Event.GetTitle()
	event.Description = in.Event.GetDescription()
	event.StartDate, _ = time.Parse(time.RFC3339, in.Event.GetStartDate())
	event.EndDate, _ = time.Parse(time.RFC3339, in.Event.GetEndDate())
	if err := json.Unmarshal(encodedEvent, event); err != nil {
		return nil, err
	}

	if err := s.app.CreateEvent(*event); err != nil {
		s.logger.Error(fmt.Errorf("fail create event %w ", err))
		return nil, err
	}

	return &pb.CreateResponse{}, nil
}

func (s *CalendarService) Edit(_ context.Context, in *pb.EditRequest) (*pb.EditResponse, error) {
	encodedEvent, _ := json.Marshal(in.Event)
	var event storage.Event
	if err := json.Unmarshal(encodedEvent, &event); err != nil {
		return nil, err
	}

	if err := s.app.EditEvent(event.ID.String(), event); err != nil {
		s.logger.Error(fmt.Errorf("event edition fail: %w", err))
		return nil, err
	}

	return &pb.EditResponse{}, nil
}

func (s *CalendarService) Delete(_ context.Context, in *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	if err := s.app.DeleteEvent(in.ID); err != nil {
		s.logger.Error(fmt.Errorf("event deleting fail: %w", err))
		return nil, err
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
			StartDate:   event.StartDate.Format(storage.DateLayout),
			EndDate:     event.EndDate.Format(storage.DateLayout),
		}
	}
	fmt.Println("result is ", result)
	return &pb.ListResponse{List: result}, nil
}
