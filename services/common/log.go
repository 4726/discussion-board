package common

import (
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const LogMsgCtxKey = "common_log_info"

var log = logrus.New()

type Log struct {
	entry *logrus.Entry
}

func NewLogger(serviceName string) *Log {
	entry := log.WithFields(logrus.Fields{
		"service": serviceName,
	})
	log.SetFormatter(&logrus.JSONFormatter{})

	log.SetOutput(os.Stdout)
	return &Log{entry}
}

func (l *Log) SetOutput(output io.Writer) {
	log.SetOutput(output)
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
		i, ok := c.Get(LogMsgCtxKey)
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
			"ClientIP":   c.ClientIP(),
			"Method":     c.Request.Method,
			"Path":       c.Request.URL.Path,
			"Message":    logMessage,
			"Latency":    finish.Sub(start),
		})

		switch v := statusCode; {
		case v < 400:
			e.Info()
		case v < 500:
			e.Warn()
		default:
			e.Error()
		}
	}
}
