package main

import (
	"github.com/4726/discussion-board/services/common"
)

type Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
	common.DefaultConfig `mapstructure:",squash"`
}

func ConfigFromFile(file string) (Config, error) {
	c := Config{}
	err := common.LoadConfig(file, "media", &c)
	return c, err
}
