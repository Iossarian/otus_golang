package logger

import (
	"fmt"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/config"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Logger struct {
	filePath string
}

func New(c config.Config) *Logger {
	return &Logger{
		filePath: c.LoggingFile,
	}
}

func (l *Logger) LogRequest(r *http.Request, statusCode int, requestDuration time.Duration) {
	l.Info(fmt.Errorf(
		"%s %s %s %s %s %d %s %s",
		r.RemoteAddr,
		time.Now().Format(time.RFC1123Z),
		r.Method,
		r.RequestURI,
		r.Proto,
		statusCode,
		requestDuration,
		r.UserAgent(),
	))
}

func (l *Logger) Info(error error) {
	logFile, err := os.OpenFile(l.filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			panic("can not close log file")
		}
	}(logFile)

	wrt := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(wrt)

	log.Println("Info: " + error.Error())
}

func (l *Logger) Error(error error) {
	logFile, err := os.OpenFile(l.filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			panic("can not close log file")
		}
	}(logFile)

	wrt := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(wrt)

	log.Println("Error: " + error.Error())
}
