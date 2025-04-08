package logging

import (
	"go.uber.org/zap/zapcore"
)

type LogLevel string

const (
	InfoLevel  LogLevel = "info"
	DebugLevel LogLevel = "debug"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
)

func (l LogLevel) String() string {
	return map[LogLevel]string{
		InfoLevel:  "info",
		DebugLevel: "debug",
		WarnLevel:  "warn",
		ErrorLevel: "error",
	}[l]
}

func (l LogLevel) ToZapLevel() zapcore.Level {
	switch l {
	case InfoLevel:
		return zapcore.InfoLevel
	case DebugLevel:
		return zapcore.DebugLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

type ReleaseLevel string

const (
	DevelopmentLevel ReleaseLevel = "development"
	ProductionLevel  ReleaseLevel = "production"
)

func (l ReleaseLevel) String() string {
	return map[ReleaseLevel]string{
		DevelopmentLevel: "development",
		ProductionLevel:  "production",
	}[l]
}
