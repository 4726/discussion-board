package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/4726/discussion-board/services/common"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-lib/metrics"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

var log = common.NewLogger("api-gateway")
var closer io.Closer

func init() {
	cfg := jaegercfg.Configuration{
		ServiceName: "gateway",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	tracer, c, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		log.Entry().Error("could not setup opentracing")
	}
	closer = c
	//needs to set as global for grpc middleware to work
	opentracing.SetGlobalTracer(tracer)
}

func main() {
	configPath := flag.String("config", "config.json", "config file path")

	flag.Parse()

	cfg, err := ConfigFromFile(*configPath)
	if err != nil {
		log.Entry().Fatal(err)
	}

	api, err := NewRestAPI(cfg)
	if err != nil {
		log.Entry().Fatal(err)
	}

	err = api.Run(fmt.Sprintf(":%v", cfg.ListenPort))

	log.Entry().Fatal(err)
	defer closer.Close()
}
