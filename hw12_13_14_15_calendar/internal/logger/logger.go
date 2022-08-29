package logger

import (
	"net/http"
	"os"
	"time"

	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/config"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
)

type Logger struct {
	logger *logrus.Logger
}

func New(c *config.Config, logFile *os.File) *Logger {
	logger := logrus.New()
	logLevel, _ := logrus.ParseLevel(c.Logger.Level)
	logger.SetLevel(logLevel)
	logger.SetOutput(logFile)
	logger.SetOutput(os.Stdout)

	return &Logger{logger: logger}
}

func (l *Logger) LogHTTPRequest(r *http.Request, statusCode int, requestDuration time.Duration) {
	l.logger.Infof(
		"%s %s %s %s %s %d %s %s",
		r.RemoteAddr,
		time.Now().Format(time.RFC1123Z),
		r.Method,
		r.RequestURI,
		r.Proto,
		statusCode,
		requestDuration,
		r.UserAgent(),
	)
}

func (l *Logger) LogGRPCRequest(code codes.Code, method, address string, requestDuration time.Duration) {
	l.logger.Infof(
		"%s %s %s %s %s",
		code,
		time.Now().Format(time.RFC1123Z),
		method,
		address,
		requestDuration,
	)
}

func (l *Logger) Info(err error) {
	l.logger.Info(err.Error())
}

func (l *Logger) Error(err error) {
	l.logger.Error(err.Error())
}
