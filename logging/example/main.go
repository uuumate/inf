package main

import (
	"github.com/uuumate/inf/logging"
	"go.uber.org/zap"

	"time"
)

func main() {
	logging.InitLogger(&logging.LogConfig{
		LogLevel: logging.LogLevelDebug,
	})

	logging.Debug("this is a debug log")

	logging.Debugf("this is a Debugf log, type: %s", "debugf")

	for i := 0; i < 10; i++ {
		logging.Debugw("this is a Debugw log", zap.String("type", "debugw"))
		time.Sleep(time.Second)
	}

	logging.Info("this is a info log")

	logging.Info("this is a info log")

	logging.Infof("this is a info log, type: %s", "infof")

	logging.Infow("this is a infow log", zap.String("type", "infow"))

	loggingTest()

	logging.Sync()

}

func loggingTest() {
	logger := logging.For("logger", "func", "loggingTest")
	logger.Infof("HAHHAHAH")
}
