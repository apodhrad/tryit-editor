package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	DEFAULT_LOG_LEVEL string = "INFO"
	DEFAULT_LOG_FILE  string = "/tmp/tryit-editor.log"
)

func init() {
	setOutputs()
	setFormatter()
	setLogLevel(os.Getenv("LOG_LEVEL"))
}

func setLogLevel(logLevel string) {
	logLevel = strings.ToLower(logLevel)
	if logLevel == "" {
		logLevel = DEFAULT_LOG_LEVEL
	}
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		panic(err)
	}
	logrus.SetLevel(level)
}

func setFormatter() {
	logrus.SetFormatter(&myFormatter{logrus.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05",
		ForceColors:            true,
		DisableLevelTruncation: true,
	}})
}

func setOutputs() {
	file, err := os.OpenFile(DEFAULT_LOG_FILE, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, file)
	logrus.SetOutput(mw)
}

type myFormatter struct {
	logrus.TextFormatter
}

func (f *myFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// this whole mess of dealing with ansi color codes is required if you want the colored output otherwise you will lose colors in the log levels
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = 31 // gray
	case logrus.WarnLevel:
		levelColor = 33 // yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = 31 // red
	default:
		levelColor = 36 // blue
	}
	levelString := strings.ToUpper(entry.Level.String())[0:4]
	return []byte(fmt.Sprintf("[%s] [\x1b[%dm%s\x1b[0m] %s\n", entry.Time.Format(f.TimestampFormat), levelColor, levelString, entry.Message)), nil
}
