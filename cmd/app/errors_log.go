package main

import (
	"os"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

var DefaultLevels = []logrus.Level{
	logrus.FatalLevel,
	logrus.ErrorLevel,
	logrus.WarnLevel,
	logrus.InfoLevel,
	logrus.DebugLevel,
}

type Hook struct {
	App     *newrelic.Application
	levels  []logrus.Level
	logChan chan<- string
}

func NewHook(App *newrelic.Application, levels []logrus.Level) *Hook {
	logChan := make(chan string, 10)

	go func() {
		App.WaitForConnection(5 * time.Second)
		logrus.Info("NR successfully connected")
		for message := range logChan {
			App.RecordLog(newrelic.LogData{
				Message: message,
			})
		}
	}()

	return &Hook{App, levels, logChan}
}

func (h *Hook) Fire(entry *logrus.Entry) error {
	message, _ := entry.String()
	h.logChan <- message
	return nil
}

func (h *Hook) Levels() []logrus.Level {
	if h.levels == nil {
		return DefaultLevels
	}
	return h.levels
}

func HandlePanic() {
	if r := recover(); r != nil {
		logrus.Errorf("Panic: %+v", r)
		os.Exit(1)
	}
}
