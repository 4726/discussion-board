package main

import (
	"flag"
	"fmt"
	"github.com/4726/discussion-board/services/common"
)

var log = common.NewLogger("search")

func main() {
	configPath := flag.String("config", "config.json", "config file path")
	flag.Parse()

	cfg, err := ConfigFromJSON(*configPath)
	if err != nil {
		log.Entry().Error(err)
		return
	}

	api, err := NewRestAPI(cfg.ESIndex, cfg.ESAddr)
	if err != nil {
		log.Entry().Fatal(err)
	}
	err = api.Run(fmt.Sprintf(":%v", cfg.ListenPort))
	log.Entry().Error(err)
}
