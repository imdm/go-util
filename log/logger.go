package log

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"time"

	"github.com/imdm/go-util/metrics"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	metricsLogCount = metrics.NewCounter("log_count", []string{"level"})
)

// Logger wraps zap.SugaredLogger with two output cores: a rotate file and color console.
// It also registers a custom metrics collection hook.
type Logger struct {
	*zap.SugaredLogger
	level      zapcore.Level
	fileLog    bool
	path       string
	consoleLog bool
}

// Option configures a logger.
type Option func(*Logger)

// WithLevel returns a Option that sets the logger's level field.
func WithLevel(lvl string) Option {
	return func(l *Logger) {
		l.level = zapLevel(lvl)
	}
}

// WithFileLog returns a Option that sets the logger's fileLog field.
// If fileLog equals to true, the logger will output log to file.
func WithFileLog(b bool) Option {
	return func(l *Logger) {
		l.fileLog = b
	}
}

// WithFilePath returns a Option that sets the logger's path field.
// path is file path of fileLog. It only takes effect when fileLog equals to true.
func WithFilePath(p string) Option {
	return func(l *Logger) {
		l.path = p
	}
}

// WithConsoleLog returns a Option that sets the logger's consoleLog field.
// If consoleLog equals to true, the logger will output log to console.
func WithConsoleLog(b bool) Option {
	return func(l *Logger) {
		l.consoleLog = b
	}
}

func zapLevel(l string) zapcore.Level {
	switch l {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}

// New returns a Logger instance.
func New(options ...Option) *Logger {
	l := &Logger{
		level:      zap.DebugLevel,
		consoleLog: true,
	}
	for _, o := range options {
		o(l)
	}
	var cores []zapcore.Core
	if l.consoleLog {
		cores = append(cores, l.consoleCore())
	}
	if l.fileLog {
		cores = append(cores, l.fileCore())
	}
	core := zapcore.RegisterHooks(zapcore.NewTee(cores...), func(e zapcore.Entry) error {
		return metricsLogCount.Add(1, map[string]string{"level": e.Level.String()})
	})
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	l.SugaredLogger = zapLogger.Sugar()
	return l
}

func (l *Logger) consoleCore() zapcore.Core {
	conf := zap.NewProductionEncoderConfig()
	conf.EncodeLevel = zapcore.CapitalColorLevelEncoder
	conf.EncodeTime = cstTimeEncoder

	lvlEnabler := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= l.level
	})
	return zapcore.NewCore(zapcore.NewJSONEncoder(conf), zapcore.Lock(os.Stderr), lvlEnabler)
}

func (l *Logger) fileCore() zapcore.Core {
	path, err := filepath.Abs(l.path)
	if err != nil {
		panic(err)
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		panic(err)
	}
	accessWriter := getWriter(path, "access.log")
	errorWriter := getWriter(path, "error.log")

	conf := zap.NewProductionEncoderConfig()
	conf.EncodeTime = cstTimeEncoder
	encoder := zapcore.NewJSONEncoder(conf)

	accessEnabler := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= l.level && lev < zapcore.ErrorLevel
	})
	errorEnabler := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= l.level && lev >= zapcore.ErrorLevel
	})
	accessCore := zapcore.NewCore(encoder, accessWriter, accessEnabler)
	errorCore := zapcore.NewCore(encoder, errorWriter, errorEnabler)
	return zapcore.NewTee(accessCore, errorCore)
}

func getWriter(path, fileName string) zapcore.WriteSyncer {
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(path, fileName),
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     30, // days
	})
}

// convert time zone to cst(shanghai)
func cstTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	//loc, err := time.LoadLocation("Asia/Shanghai")
	//if err != nil {
	//	fmt.Printf("load location failed, error: %s", err.Error())
	//	return
	//}
	//enc.AppendString(t.In(loc).Format("2006-01-02 15:04:05"))
	enc.AppendString(t.In(time.Local).Format("2006-01-02 15:04:05"))
}
