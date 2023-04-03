package utils

import (
	"github.com/sirupsen/logrus"
	"os"
)

func InitLogRus() {
	config, err := LoadConfig("./config")
	if err != nil {
		logrus.Infof("Cannot load config: %s", err)
	}

	formatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}
	logrus.SetFormatter(formatter)

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	level, err := logrus.ParseLevel(config.GenData.LogLevel)
	if err != nil {
		logrus.Errorf("Error parsing: %s", err)
	}

	logrus.SetLevel(level)
}
