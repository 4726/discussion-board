package main

import (
	"github.com/4726/discussion-board/services/common"
)

type Config struct {
	ListenPort   int
	ServiceAddrs Addrs
}

type Addrs struct {
	Likes, Media, PostsRead, PostsWrite, Search, User string
}

func ConfigFromFile(file string) (Config, error) {
	c := Config{}
	err := common.LoadConfig(file, "gateway", &c)
	return c, err
}
