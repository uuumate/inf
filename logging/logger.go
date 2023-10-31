package logging

import (
	"time"

	"github.com/uuumate/inf/rolling"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogLevel int8

const (
	LogLevelDebug = LogLevel(zapcore.DebugLevel)
	LogLevelInfo  = LogLevel(zapcore.InfoLevel)
	LogLevelError = LogLevel(zapcore.ErrorLevel)
)

type RollingFormat string

const (
	RollingFormatHour     = RollingFormat(rolling.HourlyRolling)
	RollingFormatDay      = RollingFormat(rolling.DailyRolling)
	RollingFormatDayMonth = RollingFormat(rolling.MonthlyRolling)
)

type LogConfig struct {
	LogPath  string
	LogLevel LogLevel
	Rolling  RollingFormat
}

var defaultEncoderConfig = zapcore.EncoderConfig{
	MessageKey:     "message",
	TimeKey:        "time",
	CallerKey:      "caller",
	LevelKey:       "level",
	LineEnding:     zapcore.DefaultLineEnding,
	NameKey:        "Logger",
	StacktraceKey:  "stack",
	EncodeCaller:   zapcore.ShortCallerEncoder,
	EncodeLevel:    zapcore.CapitalColorLevelEncoder,
	EncodeTime:     MilliSecondTimeEncoder,
	EncodeDuration: zapcore.StringDurationEncoder,
	EncodeName:     zapcore.FullNameEncoder,
}

func MilliSecondTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func (c *LogConfig) withDefault() {
	if len(c.LogPath) == 0 {
		c.LogPath = "logs"
	}

	if len(c.Rolling) == 0 {
		c.Rolling = RollingFormatDay
	}

	if c.LogLevel < -1 || c.LogLevel > 1 {
		c.LogLevel = LogLevelInfo
	}
}

type Logger struct {
	*zap.SugaredLogger
	rollingFiles []*rolling.File
	loglevel     zap.AtomicLevel
}

var defaultLogger *Logger

func InitLogger(cfg *LogConfig) {
	cfg.withDefault()

	defaultLogger = &Logger{}
	level := zap.NewAtomicLevelAt(zap.DebugLevel)
	debugRollingFile := rolling.NewRollingFile(cfg.LogPath, "debug")
	infoRollingFile := rolling.NewRollingFile(cfg.LogPath, "info")
	errorRollingFile := rolling.NewRollingFile(cfg.LogPath, "error")
	debugSyncWriter := zapcore.AddSync(debugRollingFile)
	infoSyncWriter := zapcore.AddSync(infoRollingFile)
	errorSyncWriter := zapcore.AddSync(errorRollingFile)

	debugLogEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if cfg.LogLevel > LogLevel(zapcore.DebugLevel) {
			return false
		}
		return level.Enabled(lvl)
	})
	errorLogEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})
	infoLogEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return cfg.LogLevel <= LogLevel(zapcore.InfoLevel) && zapcore.InfoLevel == lvl
	})

	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewConsoleEncoder(defaultEncoderConfig), debugSyncWriter, debugLogEnabler),
		zapcore.NewCore(zapcore.NewConsoleEncoder(defaultEncoderConfig), infoSyncWriter, infoLogEnabler),
		zapcore.NewCore(zapcore.NewConsoleEncoder(defaultEncoderConfig), errorSyncWriter, errorLogEnabler),
	)

	defaultLogger.rollingFiles = append(defaultLogger.rollingFiles, debugRollingFile, infoRollingFile, errorRollingFile)
	defaultLogger.SugaredLogger = zap.New(core).WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
	defaultLogger.SugaredLogger.Named("ydnLogger")

	go defaultLogger.syncTimer()
}

func (l *Logger) Sync() {
	for i := 0; i < len(l.rollingFiles); i++ {
		_ = l.rollingFiles[i].Sync()
	}
}

func (l *Logger) syncTimer() {
	ticker := time.NewTicker(time.Millisecond * 1000)
	for {
		select {
		case <-ticker.C:
			l.Sync()
		}
	}
}
