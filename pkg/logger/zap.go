package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newZapLogger(level,timeFromat string) *zap.Logger{
	globallevel:=parseLevel(level)
	highPriorty:=zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l>=zapcore.ErrorLevel
	})
	lowPriority:=zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l>=globallevel && l <zapcore.ErrorLevel
	})
	consoleInfos:=zapcore.Lock(os.Stdout)
	consoleErrors:=zapcore.Lock(os.Stderr)

	encoderCfg:=zap.NewDevelopmentEncoderConfig()
	if len(timeFromat)>0{
		customTimeFormat = timeFromat
		encoderCfg.EncodeTime=customTimeEncoder
	}else{
		encoderCfg.EncodeTime=zapcore.ISO8601TimeEncoder
	}

	consoleEncode:=zapcore.NewJSONEncoder(encoderCfg)
	core:=zapcore.NewTee(
		zapcore.NewCore(consoleEncode,consoleErrors,highPriorty),
		zapcore.NewCore(consoleEncode,consoleInfos,lowPriority),
	)
	logger:=zap.New(core)
	return logger
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(customTimeFormat))
}

func parseLevel(level string) zapcore.Level {
	switch level {
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelError:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
func GetZapLogger(l Logger) *zap.Logger {
	if l == nil {
		return newZapLogger(LevelInfo, time.RFC3339)
	}

	switch v := l.(type) {
	case *loggerImpl:
		return v.zap
	default:
		l.Info("logger.WithFields: invalid logger type, creating a new zap logger", String("level", LevelInfo), String("time_format", time.RFC3339))
		return newZapLogger(LevelInfo, time.RFC3339)
	}
}
