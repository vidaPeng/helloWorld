package pkg

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.SugaredLogger

func InitLogger() {
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:       "time",
			LevelKey:      "level",
			NameKey:       "logger",
			CallerKey:     "caller",
			MessageKey:    "msg",
			StacktraceKey: "stacktrace",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   zapcore.CapitalColorLevelEncoder,
			EncodeTime:    zapcore.ISO8601TimeEncoder,
			EncodeCaller:  zapcore.ShortCallerEncoder,
		},
		DisableStacktrace: true, // 不需要打印对战信息
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
	}

	logger, err := cfg.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	log = logger.Sugar()
}

func addTrace(ctx context.Context, fields []zap.Field) []zap.Field {
	sc := trace.SpanFromContext(ctx).SpanContext()
	if sc.IsValid() {
		fields = append(fields,
			zap.String("trace_id", sc.TraceID().String()),
		)
	}
	return fields
}

func InfoTrace(ctx context.Context, msg string, fields ...zap.Field) {
	log.Desugar().With(addTrace(ctx, fields)...).Info(msg)
}

func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}
