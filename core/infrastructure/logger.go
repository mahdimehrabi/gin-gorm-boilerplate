package infrastructure

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger structure
type Logger struct {
	Zap *zap.SugaredLogger
}

// NewLogger sets up logger
func NewLogger(env Env) Logger {

	config := zap.NewDevelopmentConfig()

	if env.Environment == "development" {
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, _ := config.Build(zap.Hooks(func(entry zapcore.Entry) error {
		if entry.Level == zapcore.ErrorLevel {
			defer sentry.Flush(2 * time.Second)
			sentry.CaptureMessage(fmt.Sprintf("%s, Line No: %d :: %s", entry.Caller.File, entry.Caller.Line, entry.Message, entry.Stack))
		}
		return nil
	}))

	sugar := logger.Sugar()

	return Logger{
		Zap: sugar,
	}
}
