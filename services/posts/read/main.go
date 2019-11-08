package main

import (
	"flag"
	"fmt"
)

func main() {
	configPath := flag.String("config", "config.json", "config file path")
	flag.Parse()

	cfg, err := ConfigFromJSON(*configPath)
	if err != nil {
		standardLoggingEntry().Error(err)
		return
	}

	api, err := NewRestAPI(cfg)
	if err != nil {
		standardLoggingEntry().Error(err)
	}

	standardLoggingEntry().Infof("starting server on port: %v", cfg.ListenPort)

	err = api.Run(fmt.Sprintf(":%v", cfg.ListenPort))
	standardLoggingEntry().Error(err)
}
