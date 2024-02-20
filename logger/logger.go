package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"
)

type LogLevel int

// 日志级别常量
const (
	DEBUG   LogLevel = 1
	INFO    LogLevel = 2
	WARNING LogLevel = 3
	ERROR   LogLevel = 4
)

// ANSI转义码颜色常量
const (
	RED    = "\033[31m"
	YELLOW = "\033[33m"
	CYAN   = "\033[36m"
	WHITE  = "\033[37m"
)

// LogFormatter 日志格式化工具
type LogFormatter struct {
}

// 日志格式化
func (lf LogFormatter) format(lv LogLevel, msg any, withColor bool) string {
	var level string
	var color string
	var timestamp time.Time

	switch lv {
	case DEBUG:
		level = "<DEBUG>"
		color = WHITE
		break
	case INFO:
		level = "<INFO>"
		color = CYAN
		break
	case WARNING:
		level = "<WARNING>"
		color = YELLOW
		break
	case ERROR:
		level = "<ERROR>"
		color = RED
		break
	}

	timestamp = time.Now()

	var result string
	if withColor {
		result = fmt.Sprintf(
			"%s[PoliteDog] %-10s %v | %s\n",
			color,
			level,
			timestamp.Format("2006-01-02 15:04:05"),
			msg,
		)
	} else {
		result = fmt.Sprintf(
			"[PoliteDog] %-10s %v | %s\n",
			level,
			timestamp.Format("2006-01-02 15:04:05"),
			msg,
		)
	}

	return result
}

// Logger 日志打印封装
type Logger struct {
	writers   []logWriter
	Level     LogLevel
	LogPath   string
	Formatter LogFormatter
}

func NewLogger() *Logger {
	return &Logger{}
}

func DefaultLogger() *Logger {
	logger := NewLogger()
	logger.Level = DEBUG
	logger.writers = append(logger.writers, logWriter{
		level:  logger.Level,
		output: os.Stdout,
	})

	return logger
}

func fileWriter(filename string) io.Writer {
	writer, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	return writer
}

type logWriter struct {
	level  LogLevel
	output io.Writer
}

func (l *Logger) SetLogPath(logPath string) {
	l.LogPath = logPath
	writer := logWriter{
		level:  l.Level,
		output: fileWriter(path.Join(l.LogPath, "app.log")),
	}
	l.writers = append(l.writers, writer)
}

func (l *Logger) print(lv LogLevel, msg any) {
	// Logger设置的级别高于传入的级别时，不打印
	if l.Level > lv {
		return
	}

	var formatStr string
	for _, writer := range l.writers {
		if writer.output == os.Stdout {
			formatStr = l.Formatter.format(lv, msg, true)
		} else {
			formatStr = l.Formatter.format(lv, msg, false)
		}

		_, _ = fmt.Fprintf(writer.output, formatStr)
	}
}

func (l *Logger) Debug(msg any) {
	l.print(DEBUG, msg)
}

func (l *Logger) Info(msg any) {
	l.print(INFO, msg)
}

func (l *Logger) Warning(msg any) {
	l.print(WARNING, msg)
}

func (l *Logger) Error(msg any) {
	l.print(ERROR, msg)
}
