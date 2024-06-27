package logging

import (
	"fmt"
	"path"
	"runtime"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

// Logger представляет структуру журнала.
type Logger struct {
	*logrus.Entry
	colored bool // Флаг включения цветного вывода. Зависит от значения LogFormat.
}

// LogFormat формат вывода журнала.
type LogFormat string

const (
	// LogFormatText формат обычного текста для записи в файл.
	LogFormatText LogFormat = "text"
	// LogFormatJSON формат JSON.
	LogFormatJSON LogFormat = "json"
)

// NewLogger создает и возвращает Logger.
func NewLogger(level int, format LogFormat, hooks ...logrus.Hook) *Logger {
	l := logrus.New()
	l.SetReportCaller(false)
	for _, hook := range hooks {
		l.AddHook(hook)
	}

	var coloredOutput bool
	switch format {
	case LogFormatText:
		l.Formatter = &logrus.TextFormatter{
			TimestampFormat: "02-01-2006 15:04:05",
			FullTimestamp:   true,
			DisableColors:   true,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := path.Base(f.File)
				return "", fmt.Sprintf("%s:%d", filename, f.Line)
			},
		}
	case LogFormatJSON:
		l.Formatter = &logrus.JSONFormatter{
			TimestampFormat: "02-01-2006 15:04:05",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		}
	default:
		l.Formatter = &logrus.TextFormatter{
			TimestampFormat: "02-01-2006 15:04:05",
			FullTimestamp:   true,
			DisableColors:   false,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := path.Base(f.File)
				return "", fmt.Sprintf("%s:%d", filename, f.Line)
			},
		}
		coloredOutput = true
	}
	l.SetLevel(logrus.Level(level))

	return &Logger{
		Entry:   logrus.NewEntry(l),
		colored: coloredOutput,
	}
}

// LoggerColor представляет типы цветов журнала.
type LoggerColor = color.Attribute

// Возможные цвета текста.
var (
	LoggerColorRed    = color.FgRed
	LoggerColorYellow = color.FgYellow
	LoggerColorGreen  = color.FgGreen
)

// ColorizeString применяет цвет к переданной строке.
func (l Logger) ColorizeString(s string, c LoggerColor) string {
	if !l.colored {
		return s
	}
	return color.New(c).Sprint(s)
}
