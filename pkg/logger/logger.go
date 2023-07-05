package logger 

import (
	"io"
	"log"
	"os"

	"github.com/kuzxnia/mongoload/pkg/config"
)

type Logger struct {
	config *config.Config
	log    *log.Logger
}

var std = NewLogger(nil, os.Stderr, "", log.LstdFlags)

func Default() *Logger {
	return std
}

func NewLogger(config *config.Config, out io.Writer, prefix string, flag int) *Logger {
	// config debug and above to file of stderr if not set and debug enabled

	return &Logger{
		config: config,
		log:    log.New(out, prefix, flag),
	}
}

func (l *Logger) SetConfig(config *config.Config) {
	l.config = config
}

func (l *Logger) Info(msg ...any) {
	// msg = append([]string{"INFO"}, msg...)
	//
	l.log.Println("INFO", msg)
}

func (l *Logger) Debug(msg ...any) {
	// msg = append([]string{"DEBUG"}, msg...)

	if l.config != nil && l.config.Debug {
		l.log.Println("DEBUG", msg)
	}
}

func (l *Logger) Error(msg ...any) {
	// msg = append([]string{"Error"}, msg...)

	if l.config != nil && l.config.Debug {
		l.log.Println("ERROR", msg)
	}
}

// func (l *Logger) printlnWithPrefix(prefix string, msg ...string) {
//   l.log.Println(prefix, msg)
// }
