package main

import (
	"github.com/4726/discussion-board/services/common"
)

type Config struct {
	ESIndex, ESAddr string
	common.DefaultConfig `mapstructure:",squash"`
}

func ConfigFromFile(file string) (Config, error) {
	c := Config{}
	err := common.LoadConfig(file, "search", &c)
	return c, err
}
