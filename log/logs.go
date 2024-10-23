package log

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	currentDate := time.Now().Format("2006-01-02")
	file, err := os.OpenFile("log/"+currentDate+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Logger.Fatal(err)
	}
	Logger.SetOutput(file)
}
