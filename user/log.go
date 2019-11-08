package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strings"
	"github.com/bshuster-repo/logrus-logstash-hook"
	"net"
)

var (
	appFields = logrus.Fields{"service": "user"}
	log       = logrus.New()
)

func init() {
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

	file, err := os.OpenFile("logs/user.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		standardLoggingEntry().Error(err)
		log.SetOutput(os.Stderr)
	}	

	conn, err := net.Dial("tcp", "localhost:8911")
	if err != nil {
		standardLoggingEntry().Error(err)
	} else {
		hook := logrustash.New(conn, logrustash.DefaultFormatter(logrus.Fields{}))
		log.Hooks.Add(hook)
	}
}

func standardLoggingEntry() *logrus.Entry {
	return log.WithFields(appFields)
}

func standardRequestLoggingEntry(ctx *gin.Context) *logrus.Entry {
	return standardLoggingEntry().
		WithField("from", ctx.ClientIP()).
		WithField("statusCode", ctx.Writer.Status())
}
