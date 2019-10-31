package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strings"
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
	log.SetOutput(os.Stderr)
}
