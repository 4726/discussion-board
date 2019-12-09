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

	cfg, err := ConfigFromFile(*configPath)
	if err != nil {
		log.Entry().Error(err)
		return
	}

	api, err := NewApi(cfg)
	if err != nil {
		log.Entry().Fatal(err)
	}
	err = api.Run(fmt.Sprintf(":%v", cfg.ListenPort))
	log.Entry().Error(err)
}
