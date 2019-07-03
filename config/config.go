package config

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
)

type Config struct {
	appId              int64
	privateKeyFileName string
	debug              bool
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

func (c *Config) IsDebugging() bool {
	return c.debug
}

func (c *Config) validate() {
	if c.appId == 0 {
		panic(errors.New("environment variable 'GO_GITHUB_WIP_APP_ID' cannot be empty or 0"))
	}
	if len(c.privateKeyFileName) == 0 {
		panic(errors.New("environment variable 'GO_GITHUB_WIP_APP_PRIVATE_KEY' cannot be empty"))
	}
	data, err := ioutil.ReadFile(c.privateKeyFileName)
	if err != nil {
		panic(errors.New("unable to read private key: " + err.Error()))
	}
	if len(data) == 0 {
		panic("private key cannot be empty")
	}
}

func Get() *Config {
	if config == nil {
		appIdAsInt, _ := strconv.Atoi(os.Getenv("GO_GITHUB_WIP_APP_ID"))
		config = &Config{
			appId:              int64(appIdAsInt),
			privateKeyFileName: os.Getenv("GO_GITHUB_WIP_APP_PRIVATE_KEY"),
			debug:              os.Getenv("GO_GITHUB_WIP_DEBUG") == "true",
		}
		config.validate()
	}
	return config
}

// Manually instantiate Config. This bypasses configuration validation. Used for testing
func Set(appId int64, privateKeyFileName string, debug bool) {
	config = &Config{
		appId:              appId,
		privateKeyFileName: privateKeyFileName,
		debug:              debug,
	}
}
