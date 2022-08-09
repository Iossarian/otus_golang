package internalhttp

import (
	"context"
	"fmt"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"time"
)

type Server struct {
	app    Application
	logger Logger
	config Config
	server *http.Server
}

type Logger interface {
	Info(err error)
	Error(err error)
	LogRequest(r *http.Request, statusCode int, requestDuration time.Duration)
}

type Application interface {
	CreateEvent(event storage.Event)
	EditEvent(id string, e storage.Event)
	DeleteEvent(id string)
	SelectForTheDay(date time.Time) map[string]storage.Event
	SelectForTheWeek(date time.Time) map[string]storage.Event
	SelectForTheMonth(date time.Time) map[string]storage.Event
}

type Config interface {
	GetAddr() string
}

type loggingMiddleware struct {
	logger Logger
}

func NewServer(logger Logger, app Application, config Config) *Server {
	router := mux.NewRouter()
	httpServer := &http.Server{
		Addr:    config.GetAddr(),
		Handler: router,
	}

	srv := &Server{
		logger: logger,
		app:    app,
		config: config,
		server: httpServer,
	}

	router.HandleFunc("/", srv.Hello).Methods("GET")
	lm := loggingMiddleware{logger: logger}
	router.Use(lm.Middleware)

	return srv
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Error(fmt.Errorf("listen: %s", err))
		os.Exit(1)
	}
	<-ctx.Done()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info(fmt.Errorf("calendar is shutting down"))
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error(fmt.Errorf("shutdown: %s", err))
	}
	<-ctx.Done()

	return nil
}

func (s *Server) Hello(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if _, err := w.Write([]byte("hello-world")); err != nil {
		s.logger.Error(err)
	}
}
