package internalhttp

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/storage"
	"net/http"
)

type Server struct {
	app    Application
	logger Logger
	config Config
}

/*func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/all" {
		s.logger.Info(r.RemoteAddr + "[" + time.Now().UTC().String() + "]" + r.Method + r.RequestURI + strconv.Itoa(int(r.ContentLength)) + r.UserAgent())
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(s.app.GetEvents())
	}
}*/

func (s *Server) getAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(s.app.GetEvents())
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Application interface {
	CreateEvent(ctx context.Context, title, date string) error
	GetEvent(id string) *storage.Event
	DeleteEvent(id string)
	GetEvents() map[string]*storage.Event
	EditEvent(id string, e *storage.Event)
}

type Config interface {
	GetAddr() string
}

func NewServer(logger Logger, app Application, config Config) *Server {
	return &Server{
		logger: logger,
		app:    app,
		config: config,
	}
}

func (s *Server) Start(ctx context.Context) error {
	/*server := &http.Server{
		Addr:    s.config.GetAddr(),
		Handler: s,
	}*/

	http.HandleFunc("/all", loggingMiddleware(s.getAll, s))
	http.ListenAndServe(s.config.GetAddr(), nil)
	<-ctx.Done()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {

	return errors.New("fail")
}
