package logging

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-ini/ini"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggerBase struct {
	Logger      *zap.Logger
	WriteToAll  bool
	ServiceName string
}

var all *zap.Logger

func Setup() {
	all = NewLogger("", zapcore.InfoLevel, 100, 3, 7, true, "All")
}

func getRunMode() string {
	cfg, err := ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini':'%v'", err)
	}
	runMode := cfg.Section("").Key("RunMode").MustString("debug")
	if runMode == "" {
		runMode = "release"
	}
	return runMode
}

func getLogPathAndExt() (path, ext string) {
	cfg, err := ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini':'%v'", err)
	}
	path = cfg.Section("app").Key("LogSavePath").MustString("./log/")
	ext = cfg.Section("app").Key("LogFileExt").MustString("log")
	return
}

func NewLogger(filePath string, level zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool, serviceName string) *zap.Logger {
	path, ext := getLogPathAndExt()
	var realFilePath string = filePath
	if filePath == "" {
		realFilePath = path + strings.ToLower(serviceName) + "." + ext
	}
	core := newCore(realFilePath, level, maxSize, maxBackups, maxAge, compress, serviceName)
	if serviceName == "All" {
		return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	} else {
		return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.Fields(zap.String("service name", serviceName)))
	}
}

func newCore(filePath string, level zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool, serviceName string) zapcore.Core {

	serviceFileLogger := getLumberjack(filePath, maxSize, maxBackups, maxAge, compress)

	atomicLevel := zap.NewAtomicLevel()
	if getRunMode() == "debug" {
		atomicLevel.SetLevel(zapcore.DebugLevel)
	} else {
		atomicLevel.SetLevel(level)
	}

	consoleEncoderConfig := getEncoderConfig()
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	fileEncoderConfig := getEncoderConfig()
	fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)
	fileEncoder := zapcore.NewJSONEncoder(fileEncoderConfig)

	if serviceName == "All" {
		return zapcore.NewTee(
			zapcore.NewCore(fileEncoder, zapcore.AddSync(&serviceFileLogger), atomicLevel),
		)
	} else {
		return zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), atomicLevel),
			zapcore.NewCore(fileEncoder, zapcore.AddSync(&serviceFileLogger), atomicLevel),
		)
	}
}

func customizeTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func getLumberjack(filePath string, maxSize int, maxBackups int, maxAge int, compress bool) lumberjack.Logger {
	return lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
	}
}

func getEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     customizeTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
}

func (log LoggerBase) Debug(msg string, fields ...zapcore.Field) {
	defer log.Logger.Sync()
	defer all.Sync()
	log.Logger.Debug(msg, fields...)
	if log.WriteToAll {
		fields = append([]zapcore.Field{zap.String("service name", log.ServiceName)}, fields...)
		all.Debug(msg, fields...)
	}
}

func (log LoggerBase) Info(msg string, fields ...zapcore.Field) {
	defer log.Logger.Sync()
	defer all.Sync()
	log.Logger.Info(msg, fields...)
	if log.WriteToAll {
		fields = append([]zapcore.Field{zap.String("service name", log.ServiceName)}, fields...)
		all.Info(msg, fields...)
	}
}

func (log LoggerBase) Warn(msg string, fields ...zapcore.Field) {
	defer log.Logger.Sync()
	defer all.Sync()
	log.Logger.Warn(msg, fields...)
	if log.WriteToAll {
		fields = append([]zapcore.Field{zap.String("service name", log.ServiceName)}, fields...)
		all.Warn(msg, fields...)
	}
}
func (log LoggerBase) Error(msg string, fields ...zapcore.Field) {
	defer log.Logger.Sync()
	defer all.Sync()
	log.Logger.Error(msg, fields...)
	if log.WriteToAll {
		fields = append([]zapcore.Field{zap.String("service name", log.ServiceName)}, fields...)
		all.Error(msg, fields...)
	}
}

func (log LoggerBase) Fatal(msg string, fields ...zapcore.Field) {
	defer log.Logger.Sync()
	defer all.Sync()
	if log.WriteToAll {
		tempFields := append([]zapcore.Field{zap.String("service name", log.ServiceName), zap.String("FATAL", "Exit")}, fields...)
		all.DPanic(msg, tempFields...)
	}
	log.Logger.Fatal(msg, fields...)
}

func GetField(key, value string) zapcore.Field {
	return zap.String(key, value)
}
