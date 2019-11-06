package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Username, Password, DBName, Addr string
}

func ConfigFromJSON(file string) (Config, error) {
	c := Config{}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return c, err
	}
	err = json.Unmarshal(data, &c)
	return c, err
}
