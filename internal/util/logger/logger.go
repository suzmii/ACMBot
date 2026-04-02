package logger

import (
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/suzmii/ACMBot/config"

	"github.com/sirupsen/logrus"
)

const pkgPath = "internal/util/logger"

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"

	TimeColor = Blue    // Blue
	FuncColor = Magenta // Magenta

)

func init() {
	logrus.SetFormatter(&formatter{})
}

func Init() {
	cfg := config.LoadConfig().Logger

	if cfg.AlterToken != "" {
		logrus.AddHook(&AlertHook{SendKey: cfg.AlterToken})

	}
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		logrus.Fatal("invalid logger level")
	}
	logrus.SetLevel(level)
}

type formatter struct{}

func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b strings.Builder

	// 日志级别
	b.WriteString(levelColor(entry.Level))
	b.WriteString(shotLevel(entry.Level))
	b.WriteString(Reset)

	// 时间
	b.WriteString(TimeColor)
	b.WriteByte('[')
	b.WriteString(time.Now().Format("2006-01-02 15:04:05.000"))
	b.WriteString("]")
	b.WriteString(Reset)
	b.WriteByte(' ')

	// caller file and func
	funcName, file, line := findCaller()
	b.WriteString(FuncColor)
	b.WriteString(file)
	b.WriteByte(':')
	b.WriteString(strconv.Itoa(line))
	b.WriteByte(':')
	b.WriteString(funcName)
	b.WriteString("():\n")
	b.WriteString(Reset)

	// 日志内容
	b.WriteString(levelColor(entry.Level))
	b.WriteString(entry.Message)
	b.WriteString(Reset)
	b.WriteByte('\n')

	return []byte(b.String()), nil
}
func findCaller() (funcName, file string, line int) {
	for skip := 4; skip < 15; skip++ {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			continue
		}
		if !isInternal(file) {
			funcName = runtime.FuncForPC(pc).Name()
			return funcName, file, line
		}
	}
	return "unknown", "unknown", 0
}

func isInternal(file string) bool {
	return strings.Contains(file, "sirupsen/logrus") || strings.Contains(file, pkgPath)
}

func levelColor(level logrus.Level) string {
	switch level {
	case logrus.TraceLevel:
		return Magenta
	case logrus.DebugLevel:
		return Green
	case logrus.InfoLevel:
		return Cyan
	case logrus.WarnLevel:
		return Yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return Red
	default:
		return Reset
	}
}

func shotLevel(level logrus.Level) string {
	switch level {
	case logrus.TraceLevel:
		return "[TRACE]"
	case logrus.DebugLevel:
		return "[DEBUG]"
	case logrus.InfoLevel:
		return "[INFO] "
	case logrus.WarnLevel:
		return "[WARN] "
	case logrus.ErrorLevel:
		return "[ERROR]"
	case logrus.FatalLevel:
		return "[FATAL]"
	case logrus.PanicLevel:
		return "[PANIC]"
	default:
		return "[DEFAULT]"
	}
}
