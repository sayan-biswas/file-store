package logger

import (
	"os"

	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger
var LogLevel zap.AtomicLevel

func init() {

	var encoderConfig zapcore.EncoderConfig
	var option []zap.Option
	var core zapcore.Core

	encoderConfig = zapcore.EncoderConfig{
		LevelKey:       "level",
		MessageKey:     "message",
		TimeKey:        "time",
		NameKey:        "name",
		CallerKey:      zapcore.OmitKey,
		FunctionKey:    zapcore.OmitKey,
		StacktraceKey:  zapcore.OmitKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006/01/02 - 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	level := zap.NewAtomicLevel()
	core = zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(colorable.NewColorableStdout()),
		level,
	)

	option = []zap.Option{
		zap.ErrorOutput(os.Stderr),
	}

	Log = zap.New(core, option...)
}

func SetLevel(level string) error {
	var logLevel zapcore.Level
	if err := logLevel.Set(level); err != nil {
		return err
	}
	LogLevel.SetLevel(logLevel)
	return nil
}
