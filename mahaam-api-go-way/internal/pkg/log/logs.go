package logs

import (
	"fmt"
	"mahaam-api/internal/pkg/configs"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *zap.SugaredLogger

type CreateLogFunc = func(trafficId uuid.UUID, logType, message string, nodeIP string)

var createLogFunc CreateLogFunc

func Init(createLogFn CreateLogFunc) {
	var zapCfg zap.Config
	if configs.EnvName == "local" {
		zapCfg = zap.NewDevelopmentConfig()
	} else {
		zapCfg = zap.NewProductionConfig()
	}

	// Set up log rotation using lumberjack with config values
	logWriter := &lumberjack.Logger{
		Filename:   configs.LogFile,
		MaxSize:    configs.LogFileSizeLimit, // megabytes
		MaxBackups: configs.LogFileCountLimit,
		MaxAge:     28, // days (could be made configurable)
		Compress:   true,
	}

	// Create encoder config using template from config if available
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

	logger := zap.New(core)
	Log = logger.Sugar()
	createLogFunc = createLogFn
}

func Info(trafficId uuid.UUID, args ...any) {
	msg := getMessage(args...)
	logMsg := msg
	if trafficId != uuid.Nil {
		logMsg = fmt.Sprintf("TrafficId: %s, %s", trafficId, msg)
	}
	Log.Info(logMsg)
	if createLogFunc != nil {
		createLogFunc(trafficId, "Info", msg, configs.NodeIP)
	}
}

func Error(trafficId uuid.UUID, args ...any) {
	msg := getMessage(args...)
	logMsg := msg
	if trafficId != uuid.Nil {
		logMsg = fmt.Sprintf("TrafficId: %s, %s", trafficId, msg)
	}
	Log.Error(logMsg)
	if createLogFunc != nil {
		createLogFunc(trafficId, "Error", msg, configs.NodeIP)
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
