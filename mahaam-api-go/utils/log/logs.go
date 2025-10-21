package logs

import (
	"fmt"
	"mahaam-api/utils/conf"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger interface {
	Info(trafficId uuid.UUID, args ...any)
	Error(trafficId uuid.UUID, args ...any)
}

type logger struct {
	log               *zap.SugaredLogger
	createTrafficFunc func(trafficId uuid.UUID, logType, message string)
}

func NewLogger(cfg *conf.Conf, createTrafficFunc func(trafficId uuid.UUID, logType, message string)) Logger {
	var zapCfg zap.Config
	if cfg.EnvName == "local" {
		zapCfg = zap.NewDevelopmentConfig()
	} else {
		zapCfg = zap.NewProductionConfig()
	}

	// log rotation
	logWriter := &lumberjack.Logger{
		Filename:   cfg.LogFile,
		MaxSize:    cfg.LogFileSizeLimit, // megabytes
		MaxBackups: cfg.LogFileCountLimit,
		MaxAge:     28, // days
		Compress:   true,
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:          "asctime",
		LevelKey:         "levelname",
		MessageKey:       "message",
		EncodeTime:       zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		EncodeLevel:      zapcore.CapitalLevelEncoder,
		LineEnding:       "\n",
		ConsoleSeparator: " ",
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig), // use console encoder for plain text
		zapcore.AddSync(logWriter),
		zapCfg.Level,
	)

	l := zap.New(core)
	return &logger{log: l.Sugar(), createTrafficFunc: createTrafficFunc}
}

func (l *logger) Info(trafficId uuid.UUID, args ...any) {
	msg := getMessage(args...)
	logMsg := msg
	if trafficId != uuid.Nil {
		logMsg = fmt.Sprintf("TrafficId: %s, %s", trafficId, msg)
	}
	l.log.Info(logMsg)
	if trafficId != uuid.Nil {
		l.createTrafficFunc(trafficId, "Info", msg)
	}
}

func (l *logger) Error(trafficId uuid.UUID, args ...any) {
	msg := getMessage(args...)
	logMsg := msg
	if trafficId != uuid.Nil {
		logMsg = fmt.Sprintf("TrafficId: %s, %s", trafficId, msg)
	}
	l.log.Error(logMsg)
	if trafficId != uuid.Nil {
		l.createTrafficFunc(trafficId, "Error", msg)
	}
}

func getMessage(args ...any) string {
	if len(args) == 0 {
		return ""
	}
	if len(args) == 1 {
		return fmt.Sprint(args[0])
	}
	template := args[0].(string)
	fmtArgs := args[1:]

	if len(fmtArgs) == 0 {
		return template
	}

	if template != "" {
		return fmt.Sprintf(template, fmtArgs...)
	}

	if len(fmtArgs) == 1 {
		if str, ok := fmtArgs[0].(string); ok {
			return str
		}
	}
	return fmt.Sprint(fmtArgs...)
}
