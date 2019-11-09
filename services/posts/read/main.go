package main

import (
	"flag"
	"fmt"
	"github.com/4726/discussion-board/services/common"
)

var log = common.NewLogger("posts-read")

func main() {
	configPath := flag.String("config", "config.json", "config file path")
	flag.Parse()

	cfg, err := ConfigFromJSON(*configPath)
	if err != nil {
		log.Entry().Error(err)
		return
	}

	api, err := NewRestAPI(cfg)
	if err != nil {
		log.Entry().Error(err)
	}

	log.Entry().Infof("starting server on port: %v", cfg.ListenPort)

	err = api.Run(fmt.Sprintf(":%v", cfg.ListenPort))
	log.Entry().Error(err)
}