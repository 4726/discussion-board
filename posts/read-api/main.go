package main

import (
	"flag"
	"fmt"
)

func main() {
	configPath := flag.String("config", "config.json", "config file path")
	port := flag.Int("port", 14000, "listen port")
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

	standardLoggingEntry().Infof("starting server on port: %v", *port)

	err = api.Run(fmt.Sprintf(":%v", *port))
	standardLoggingEntry().Error(err)
}
