package config

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type Config struct {
	appId              int64
	privateKeyFileName string
}

var (
	config *Config
)

func (c *Config) GetAppId() int64 {
	return c.appId
}

func (c *Config) GetPrivateKeyFileName() string {
	return c.privateKeyFileName
}

func GetConfig() *Config {
	if config == nil {
		appIdAsInt, _ := strconv.Atoi(os.Getenv("GO_GITHUB_WIP_APP_ID"))
		config = &Config{
			appId:              int64(appIdAsInt),
			privateKeyFileName: os.Getenv("GO_GITHUB_WIP_APP_PRIVATE_KEY"),
		}
		// DEBUG START
		println(config.privateKeyFileName)
		data, err := ioutil.ReadFile(config.privateKeyFileName)
		if err != nil {
			log.Printf("[GetConfig] Error: %v\n", err)
		}
		println(string(data))
		// DEBUG END
	}
	return config
}
