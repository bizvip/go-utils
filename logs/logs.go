/******************************************************************************
 * Copyright (c) Archer++ 2024.                                               *
 ******************************************************************************/

package logs

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type ConsoleHook struct{}

func (hook *ConsoleHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *ConsoleHook) Fire(entry *logrus.Entry) error {
	fmt.Println(entry.Message)
	return nil
}

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{})
	// log.SetLevel(logrus.InfoLevel)

	logDir := "runtime/logs/"
	logFile := time.Now().Format("2006-01-02") + ".log"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		_ = os.MkdirAll(logDir, 0755)
	}

	file, err := os.OpenFile(logDir+logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	log.AddHook(&ConsoleHook{})
}

func Logger() *logrus.Logger {
	return log
}
