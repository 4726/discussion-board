package common

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func LoadConfig(filePath, appName string, obj interface{}) error {
	ext := filepath.Ext(filePath)[1:]
	viper.SetConfigType(ext)
	viper.SetEnvPrefix("dboard" + "_" + appName)
	for _, v := range os.Environ() {
		tokens := strings.Split(v, "=")
		if !strings.HasPrefix(tokens[0], "DBOARD"+"_"+strings.ToUpper(appName)) {
			continue
		}
		key := strings.TrimPrefix(tokens[0], "DBOARD"+"_"+strings.ToUpper(appName)+"_")
		viper.BindEnv(key)
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		// return err
	}

	if err := viper.ReadConfig(bytes.NewBuffer(data)); err != nil {
		// return err
	}

	return viper.Unmarshal(obj)
}
