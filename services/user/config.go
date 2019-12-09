package main

import (
	"github.com/4726/discussion-board/services/common"
)

type Config struct {
	Username, Password, DBName, Addr string
	ListenPort                       int
}

func ConfigFromFile(file string) (Config, error) {
	c := Config{}
	err := common.LoadConfig(file, "user", &c)
	return c, err
}
