package logger

import (
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any)
}

type LoggerImpl struct {
	*slog.Logger
}

func (l *LoggerImpl) Fatal(msg string, args ...any) {
	l.Error(msg, args...)
	os.Exit(1)
}

var String = slog.String
var Int = slog.Int
var Int64 = slog.Int64
var Any = slog.Any

var levels = map[string]slog.Level{
	"error": slog.LevelError,
	"warn":  slog.LevelWarn,
	"info":  slog.LevelInfo,
	"debug": slog.LevelDebug,
}

func New(build string, version string) Logger {
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.source", false)
	viper.SetDefault("log.time", true)

	level, ok := levels[viper.GetString("log.level")]
	if !ok {
		level = levels["info"]
	}

	var handler slog.Handler
	isTimeLogged := viper.GetBool("log.time")
	opts := &slog.HandlerOptions{
		AddSource: viper.GetBool("log.source"),
		Level:     level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && !isTimeLogged {
				return slog.Attr{}
			}

			return a
		},
	}

	switch viper.GetString("log.style") {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return &LoggerImpl{
		slog.New(handler).With(slog.String("version", version), slog.String("build", build)),
	}
}
