package common

import (
	"bytes"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/viper"
)

func LoadConfig(filePath, appName string, obj interface{}) error {
	ext := filepath.Ext(filePath)[1:]
	viper.SetConfigType(ext)
	viper.AutomaticEnv()
	viper.SetEnvPrefix("dboard" + "_" + appName)

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err := viper.ReadConfig(bytes.NewBuffer(data)); err != nil {
		return err
	}

	return viper.Unmarshal(obj)
}
