package main

import (
	"flag"
)

func main() {
	configPath := flag.String("config", "config.json", "config file path")

	flag.Parse()

	cfg, err := ConfigFromJSON(*configPath)
	if err != nil {
		standardLoggingEntry().Fatal(err)
	}

	api, err := NewRestAPI(cfg)
	if err != nil {
		standardLoggingEntry().Fatal(err)
	}

	err = api.Run(":14000")

	standardLoggingEntry().Fatal(err)
}
