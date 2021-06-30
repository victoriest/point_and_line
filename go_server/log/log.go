package log

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"strings"
	"time"
)

var _debugEnable = true
var _infoEnable = true
var _warningEnable = true
var _errorEnable = true

var (
	Global *zap.SugaredLogger
)

func InitZapLogger(debug bool, profile string, logDir string, levelString string, maxDays int) *zap.SugaredLogger {
	fileName := fmt.Sprintf("%s/%s.log", logDir, profile)
	warningFileName := fmt.Sprintf("%s/%s_warning.log", logDir, profile)
	errFileName := fmt.Sprintf("%s/%s_error.log", logDir, profile)

	if maxDays < 1 {
		// 默认30天
		maxDays = 30
	}
	level := parseLevel(levelString)
	_debugEnable = level < zapcore.InfoLevel
	_infoEnable = level < zapcore.WarnLevel
	_warningEnable = level < zapcore.ErrorLevel
	_errorEnable = level < zapcore.DPanicLevel

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "T",
		LevelKey:      "L",
		NameKey:       "N",
		CallerKey:     "C",
		MessageKey:    "M",
		StacktraceKey: "S",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel: func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(levelCapitalString(l))
		},
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	cores := make([]zapcore.Core, 0)
	if debug {
		cores = append(cores, zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level))
	}
	cores = append(cores, zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(timeDivisionWriter(fileName, maxDays)),
		level))

	if level < zapcore.ErrorLevel {
		cores = append(cores, zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(timeDivisionWriter(warningFileName, maxDays)),
			zapcore.WarnLevel))
	}
	if level < zapcore.DPanicLevel {
		cores = append(cores, zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(timeDivisionWriter(errFileName, maxDays)),
			zapcore.ErrorLevel))
	}

	//syncers = append(syncers, zapcore.AddSync(&lumberjack.Logger{
	//	Filename:   fileName,
	//	MaxSize:    500, // megabytes
	//	MaxBackups: 3,
	//	MaxAge:     28, // days
	//	LocalTime:  true,
	//}))

	// 最后创建具体的Logger
	core := zapcore.NewTee(cores...)
	caller := zap.AddCaller()
	callerSkipOpt := zap.AddCallerSkip(1)
	// From a zap.Core, it's easy to construct a Logger.
	logger := zap.New(core, caller, callerSkipOpt, zap.AddStacktrace(zap.ErrorLevel))
	Global = logger.Sugar()
	return Global
}

func Debug(msg ...interface{}) {
	Global.Debug(msg)
}

func Info(msg ...interface{}) {
	Global.Info(msg)
}

func Warn(msg ...interface{}) {
	Global.Warn(msg)
}

func Error(msg ...interface{}) {
	Global.Error(msg)
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func timeDivisionWriter(filename string, maxDays int) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YYmmdd
	// 保存30天内的日志，每天分割一次日志
	hook, err := rotatelogs.New(
		filename+".%Y%m%d", // 没有使用go风格反人类的format格式
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*time.Duration(maxDays)),
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	if err != nil {
		panic(err)
	}
	return hook
}

// 兼容string级别设置
func parseLevel(level string) zapcore.Level {
	level = strings.ToLower(level)
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// CapitalString returns an all-caps ASCII representation of the log level.
func levelCapitalString(l zapcore.Level) string {
	// Printing levels in all-caps is common enough that we should export this
	// functionality.
	switch l {
	case zapcore.DebugLevel:
		return "[DEBUG]"
	case zapcore.InfoLevel:
		return "[INFO]"
	case zapcore.WarnLevel:
		return "[WARN]"
	case zapcore.ErrorLevel:
		return "[ERROR]"
	case zapcore.DPanicLevel:
		return "[DPANIC]"
	case zapcore.PanicLevel:
		return "[PANIC]"
	case zapcore.FatalLevel:
		return "[FATAL]"
	default:
		return fmt.Sprintf("[LEVEL(%d)]", l)
	}
}
