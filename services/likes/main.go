package main

import (
	"flag"
	"fmt"
	"github.com/4726/discussion-board/services/common"
)

var log = common.NewLogger("likes")

func main() {
	configPath := flag.String("config", "config.json", "config file path")

	flag.Parse()

	cfg, err := ConfigFromJSON(*configPath)
	if err != nil {
		log.Entry().Fatal(err)
	}

	api, err := NewGRPCApi(cfg)
	if err != nil {
		log.Entry().Fatal(err)
	}

	err = api.Run(fmt.Sprintf(":%v", cfg.ListenPort))

	log.Entry().Fatal(err)
}
