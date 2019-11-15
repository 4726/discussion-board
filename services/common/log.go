package common

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var log = logrus.New()

type Log struct {
	entry *logrus.Entry
}

func NewLogger(serviceName string) *Log {
	entry := log.WithFields(logrus.Fields{
		"service": serviceName,
	})
	log.SetFormatter(&logrus.JSONFormatter{})

	// file, err := os.OpenFile(fmt.Sprintf("logs/%s.log", serviceName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err == nil {
	// 	log.SetOutput(file)
	// } else {
	// 	entry.Error(err)
		log.SetOutput(os.Stderr)
	// }

	return &Log{entry}
}

func (l *Log) Entry() *logrus.Entry {
	return l.entry
}

func (l *Log) RequestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		finish := time.Now()

		statusCode := c.Writer.Status()

		logMessage := ""
		i, ok := c.Get("log info")
		if ok {
			switch v := i.(type) {
			case string:
				logMessage = v
			case error:
				logMessage = v.Error()
			default:
			}
		}

		e := l.entry.WithFields(logrus.Fields{
			"StatusCode": statusCode,
			"ClientIP":  c.ClientIP(),
			"Method": c.Request.Method,
			"Path": c.Request.URL.Path,
			"Message": logMessage,
			"Latency": finish.Sub(start),
		})

		switch v := statusCode; {
		case v < 400:
			e.Info()
		case v < 500:
			e.Warn()
		default:
			e.Fatal()
		}
	}
}