package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/rs/zerolog"
)

type Logger struct {
	logger *zerolog.Logger
	once   sync.Once
}

func (logger *Logger) Name() string { return "Initializing cb-logger-lib!" }

func (logger *Logger) Init() error {
	logger.once.Do(func() {
		var (
			writer  io.Writer
			level   zerolog.Level
			service string
			env     string
		)

		env = strings.ToLower(os.Getenv("LOG_LEVEL"))

		switch env {
		case "host", "dev":
			writer = os.Stdout
			level = zerolog.DebugLevel
		case "prod", "test":
			// TODO: отправлять логи в Loki вместо подавления вывода
			writer = io.Discard
			level = zerolog.InfoLevel
		default:
			writer = os.Stdout
			level = zerolog.InfoLevel
		}
		service = "cb-users-auth"

		cw := zerolog.ConsoleWriter{Out: writer}
		z := zerolog.New(cw).With().
			Timestamp().
			Str("service", service).
			Logger().
			Level(level)
		logger.logger = &z
	})
	return nil
}

func Debug(args ...interface{}) {
	if AppLogger.logger != nil {
		AppLogger.logger.Debug().Msg(fmt.Sprint(args...))
	}
}

func Debugf(format string, args ...interface{}) {
	if AppLogger.logger != nil {
		AppLogger.logger.Debug().Msg(fmt.Sprintf(format, args...))
	}
}

func Info(args ...interface{}) {
	if AppLogger.logger != nil {
		AppLogger.logger.Info().Msg(fmt.Sprint(args...))
	}
}

func Infof(format string, args ...interface{}) {
	if AppLogger.logger != nil {
		AppLogger.logger.Info().Msg(fmt.Sprintf(format, args...))
	}
}

func Warn(args ...interface{}) {
	if AppLogger.logger != nil {
		AppLogger.logger.Warn().Msg(fmt.Sprint(args...))
	}
}

func Warnf(format string, args ...interface{}) {
	if AppLogger.logger != nil {
		AppLogger.logger.Warn().Msg(fmt.Sprintf(format, args...))
	}
}

func Error(args ...interface{}) {
	if AppLogger.logger != nil {
		AppLogger.logger.Error().Msg(fmt.Sprint(args...))
	}
}

func Errorf(format string, args ...interface{}) {
	if AppLogger.logger != nil {
		AppLogger.logger.Error().Msg(fmt.Sprintf(format, args...))
	}
}

// Fatal логирует и завершает процесс (через zerolog Fatal -> os.Exit(1))
func Fatal(args ...interface{}) {
	if AppLogger.logger != nil {
		AppLogger.logger.Fatal().Msg(fmt.Sprint(args...))
	}
}

func Fatalf(format string, args ...interface{}) {
	if AppLogger.logger != nil {
		AppLogger.logger.Fatal().Msg(fmt.Sprintf(format, args...))
	}
}

var AppLogger Logger
