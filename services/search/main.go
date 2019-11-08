package main

import "flag"
import "fmt"

func main() {
	configPath := flag.String("config", "config.json", "config file path")
	flag.Parse()

	cfg, err := ConfigFromJSON(*configPath)
	if err != nil {
		standardLoggingEntry().Error(err)
		return
	}

	api, err := NewRestAPI(cfg.ESIndex, cfg.ESAddr)
	if err != nil {
		standardLoggingEntry().Fatal(err)
	}
	err = api.Run(fmt.Sprintf(":%v", cfg.ListenPort))
	standardLoggingEntry().Error(err)
}
