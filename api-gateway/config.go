package main

import (
	"github.com/4726/discussion-board/services/common"
)

type Config struct {
	ListenPort        int
	TLSCert, TLSKey   string
	LikesService      Service
	MediaService      Service
	PostsReadService  Service
	PostsWriteService Service
	SearchService     Service
	UserService       Service
}

type Service struct {
	Addr, TLSCert, TLSServerName string
}

func ConfigFromFile(file string) (Config, error) {
	c := Config{}
	err := common.LoadConfig(file, "gateway", &c)
	return c, err
}
