package main

import (
	"github.com/4726/discussion-board/services/common"
)

type Config struct {
	Username, Password, DBName, Addr string
	ListenPort                       int
	TLSCert, TLSKey, TLSServerName   string
	IPWhitelist                      []string
}

func ConfigFromFile(file string) (Config, error) {
	c := Config{}
	err := common.LoadConfig(file, "likes", &c)
	return c, err
}
