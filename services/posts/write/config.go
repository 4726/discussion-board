package main

import (
	"github.com/4726/discussion-board/services/common"
)

type Config struct {
	Username, Password, DBName, Addr string
	common.DefaultConfig `mapstructure:",squash"`
}

func ConfigFromFile(file string) (Config, error) {
	c := Config{}
	err := common.LoadConfig(file, "postswrite", &c)
	return c, err
}
