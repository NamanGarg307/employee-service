package zaplogger

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/jainabhishek5986/employee-records/pkg/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger        *zap.Logger
	loggerChannel = make(chan *logItem, global.FiveThousand)
)

type Options struct {
	LogFileName string
	MaxSize     int
	MaxAge      int
	MaxBackups  int
	Level       string
}

func InitLogger(fileName string) error {
	options := Options{
		LogFileName: fileName,
		MaxSize:     global.Hundred,
		MaxAge:      global.Ten,
		MaxBackups:  global.One,
		Level:       "Info",
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	logWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   options.LogFileName,
		MaxSize:    options.MaxSize,
		MaxAge:     options.MaxAge,
		MaxBackups: options.MaxBackups,
		LocalTime:  false,
		Compress:   false,
	})

	// string to level
	var l zapcore.Level
	if err := l.Set(options.Level); err != nil {
		return err
	}

	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), logWriter), zap.NewAtomicLevelAt(l))
	logger = zap.New(core, zap.AddStacktrace(zap.ErrorLevel))

	go func() {
		for log := range loggerChannel {
			logByLevel(log.level, log.template, log.fields...)
		}
	}()

	Info(context.Background(), "Started Logger Instance")
	return nil
}

type logItem struct {
	template string
	level    zapcore.Level
	fields   []zapcore.Field
}

func Error(ctx context.Context, template string, fields ...zapcore.Field) {
	_, src, line, ok := runtime.Caller(1)
	if ok {
		fields = append(fields, zapcore.Field{Key: "caller", Type: zapcore.StringType, String: fmt.Sprintf("%s:%d", src, line)})
	}
	loggerChannel <- &logItem{level: zapcore.ErrorLevel, template: template, fields: fields}
}

func Info(ctx context.Context, template string, fields ...zapcore.Field) {
	loggerChannel <- &logItem{level: zapcore.InfoLevel, template: template, fields: fields}
}

func Debug(ctx context.Context, template string, fields ...zapcore.Field) {
	loggerChannel <- &logItem{level: zapcore.DebugLevel, template: template, fields: fields}
}

func Warn(ctx context.Context, template string, fields ...zapcore.Field) {
	loggerChannel <- &logItem{level: zapcore.WarnLevel, template: template, fields: fields}
}

func Fatal(ctx context.Context, template string, fields ...zapcore.Field) {
	loggerChannel <- &logItem{level: zapcore.FatalLevel, template: template, fields: fields}
}

func Panic(ctx context.Context, template string, fields ...zapcore.Field) {
	_, src, line, ok := runtime.Caller(1)
	if ok {
		fields = append(fields, zapcore.Field{Key: "caller", Type: zapcore.StringType, String: fmt.Sprintf("%s:%d", src, line)})
	}
	loggerChannel <- &logItem{level: zapcore.PanicLevel, template: template, fields: fields}
}

func logByLevel(level zapcore.Level, template string, fields ...zapcore.Field) {
	switch level {
	case zapcore.DebugLevel:
		logger.Debug(template, fields...)
	case zapcore.InfoLevel:
		logger.Info(template, fields...)
	case zapcore.WarnLevel:
		logger.Warn(template, fields...)
	case zapcore.ErrorLevel:
		logger.Error(template, fields...)
	case zapcore.DPanicLevel:
		logger.DPanic(template, fields...)
	case zapcore.PanicLevel:
		logger.Panic(template, fields...)
	case zapcore.FatalLevel:
		logger.Fatal(template, fields...)
	}
}
