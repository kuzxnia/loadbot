package logger

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/kuzxnia/loadbot/lbot/pkg/config"
)

type Logger struct {
	config     *config.Config
	log        *log.Logger
	outputFile *os.File
}

var std = NewLogger(nil, os.Stderr, "", log.LstdFlags)

func Default() *Logger {
	return std
}

func NewLogger(config *config.Config, out io.Writer, prefix string, flag int) *Logger {
	return &Logger{
		config: config,
		log:    log.New(out, prefix, flag),
	}
}

func (l *Logger) SetConfig(config *config.Config) {
	l.config = config

	if filePath := l.config.DebugFile; filePath != "" {
		if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
			l.outputFile, err = os.Create(filePath)
			if err != nil {
				panic("Cannot create file with given path " + filePath + " error " + err.Error())
			}
		} else if err == nil {
			l.outputFile, err = os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0o644)
			if err != nil {
				panic("Cannot open file with given path " + filePath + " error " + err.Error())
			}
		}
		l.log.SetOutput(l.outputFile)
	}
}

func (l *Logger) Info(msg ...any) {
	l.log.Println("INFO", msg)
}

func (l *Logger) Debug(msg ...any) {
	if l.config != nil && l.config.Debug {
		l.log.Println("DEBUG", msg)
	}
}

func (l *Logger) Error(msg ...any) {
	if l.config != nil && l.config.Debug {
		l.log.Println("ERROR", msg)
	}
}

func (l *Logger) CloseOutputFile() {
	if l.outputFile != nil {
		if err := l.outputFile.Close(); err != nil {
			panic("Cannot close log file: " + l.config.DebugFile)
		}
	}
}
