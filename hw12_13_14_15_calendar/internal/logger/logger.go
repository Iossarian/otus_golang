package logger

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	level      string
	infoLogger *log.Logger
	filePath   string
}

func New(level string) *Logger {
	return &Logger{
		level:    level,
		filePath: "logs",
	}
}

func (l *Logger) Info(msg string) {
	logFile, err := os.OpenFile(l.filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logFile.Close()
	wrt := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(wrt)

	log.Println(msg)
}

func (l *Logger) Error(msg string) {
	log.Print(msg)
}

// TODO
