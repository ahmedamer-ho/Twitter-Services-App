package observability

import (
	"context"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *otelzap.Logger

// InitLogger initializes a global zap logger wrapped by otelzap.
// This allows automatic extraction of trace_id and span_id from the context.
func InitLogger() error {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	
	zapLogger, err := config.Build()
	if err != nil {
		return err
	}

	// Wrap zap logger with otelzap to trace log entries
	Log = otelzap.New(zapLogger,
		otelzap.WithMinLevel(zapcore.DebugLevel),
	)

	return nil
}

// Logger returns the otelzap global logger
func Logger(ctx context.Context) otelzap.LoggerWithCtx {
	if Log == nil {
		// Fallback to no-op context logger if not initialized
		return otelzap.New(zap.NewNop()).Ctx(ctx)
	}
	return Log.Ctx(ctx)
}
