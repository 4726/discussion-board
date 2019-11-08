package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strings"
)

var log = logrus.New()

type Log struct {
	entry *logrus.Entry
}

func NewLogger(serviceName string) *Log {
	entry := log.WithFields(logrus.Fields{
		"service": serviceName,
	})

	prettyfier := func(f *runtime.Frame) (string, string) {
		var shortFileName, shortFunctionName string
		tokens := strings.Split(f.File, "/github.com/")
		if len(tokens) > 1 {
			shortFileName = fmt.Sprintf("github.com/%s:%v", tokens[1], f.Line)
		}
		tokens = strings.Split(f.Function, "/")
		if len(tokens) > 1 {
			shortFunctionName = tokens[len(tokens)-1]
		}
		return shortFunctionName, shortFileName
	}

	log.SetFormatter(&logrus.JSONFormatter{
		CallerPrettyfier: prettyfier,
	})
	log.SetReportCaller(true)

	file, err := os.OpenFile(fmt.Sprintf("logs/%s.log", serviceName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		entry.Error(err)
		log.SetOutput(os.Stderr)
	}

	return &Log{entry}
}

func (l *Log) Entry() *logrus.Entry {
	return l.entry
}

func (l *Log) HTTPRequestEntry(ctx *gin.Context) *logrus.Entry {
	return l.entry.
		WithField("from", ctx.ClientIP()).
		WithField("statusCode", ctx.Writer.Status())
}
