package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/gorilla/mux"
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
	LogHTTPRequest(r *http.Request, statusCode int, requestDuration time.Duration)
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

type loggingMiddleware struct {
	logger Logger
}

func NewServer(logger Logger, app Application, config Config) *Server {
	router := mux.NewRouter()
	httpServer := &http.Server{
		Addr:              config.GetHTTPAddr(),
		Handler:           router,
		ReadHeaderTimeout: time.Second * 10,
	}

	srv := &Server{
		logger: logger,
		app:    app,
		config: config,
		server: httpServer,
	}

	router.HandleFunc("/create", srv.Create).Methods("POST")
	router.HandleFunc("/edit", srv.Edit).Methods("PUT")
	router.HandleFunc("/delete", srv.Delete).Methods("DELETE")
	router.HandleFunc("/list", srv.List).Methods("GET")
	lm := loggingMiddleware{logger: logger}
	router.Use(lm.Middleware)

	return srv
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Error(fmt.Errorf("listen: %w", err))
		os.Exit(1)
	}
	<-ctx.Done()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info(fmt.Errorf("calendar is shutting down"))
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error(fmt.Errorf("shutdown: %w", err))
	}
	<-ctx.Done()

	return nil
}

func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	event := storage.NewEvent()
	if err := decoder.Decode(&event); err != nil {
		s.logger.Error(fmt.Errorf("event decoding fail: %w", err))
	}

	id, err := s.app.CreateEvent(*event)
	if err != nil {
		s.logger.Error(fmt.Errorf("event creation fail: %w", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) Edit(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var event storage.Event
	if err := decoder.Decode(&event); err != nil {
		s.logger.Error(fmt.Errorf("event decoding fail: %w", err))
	}

	if err := s.app.EditEvent(event.ID.String(), event); err != nil {
		s.logger.Error(fmt.Errorf("event edition fail: %w", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.logger.Error(fmt.Errorf("fail parse form %w", err))
	}

	id := r.Form.Get("id")
	if err := s.app.DeleteEvent(id); err != nil {
		s.logger.Error(fmt.Errorf("event deleting fail: %w", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) List(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.logger.Error(fmt.Errorf("fail parse form %w", err))
	}

	date := r.Form.Get("date")
	duration := r.Form.Get("duration")
	parsedDate, _ := time.Parse(storage.DateLayout, date)

	list := s.app.List(parsedDate, duration)
	if err := json.NewEncoder(w).Encode(list); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
