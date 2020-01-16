package config

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	appId              int64
	githubHost         string
	privateKeyFileName string
	prefixes           []string
	debug              bool
}

var (
	config *Config
)

func (c *Config) GetAppId() int64 {
	return c.appId
}

func (c *Config) GetGithubHost() string {
	return c.githubHost
}

func (c *Config) GetPrivateKeyFileName() string {
	return c.privateKeyFileName
}

func (c *Config) HasWipPrefix(title string) bool {
	for _, prefix := range c.prefixes {
		if strings.HasPrefix(title, prefix) {
			return true
		}
	}
	return false
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
	if len(config.prefixes) == 0 || (len(config.prefixes) == 1 && config.prefixes[0] == "") {
		log.Println("[Config.validate] No environment variable 'GO_GITHUB_WIP_PREFIXES' found, defaulting to '[WIP]' and 'WIP'")
		config.prefixes = []string{"[WIP]", "WIP"}
	} else {
		// Silently filter out invalid prefixes (e.g. empty string)
		var filteredPrefixes []string
		for _, prefix := range config.prefixes {
			if len(prefix) != 0 {
				// XXX: Should I trim the prefix and add a whitespace after it? to prevent titles like "WIPE" from
				// XXX: matching "WIP"
				filteredPrefixes = append(filteredPrefixes, prefix)
			}
		}
		config.prefixes = filteredPrefixes
	}
}

func Validate() {
	Get().validate()
}

func Get() *Config {
	if config == nil {
		appIdAsInt, _ := strconv.Atoi(os.Getenv("GO_GITHUB_WIP_APP_ID"))
		githubHost := os.Getenv("GITHUB_HOST")
		if len(githubHost) == 0 {
			githubHost = "github.com"
		}
		config = &Config{
			appId:              int64(appIdAsInt),
			githubHost:         githubHost,
			privateKeyFileName: os.Getenv("GO_GITHUB_WIP_APP_PRIVATE_KEY"),
			prefixes:           strings.Split(os.Getenv("GO_GITHUB_WIP_PREFIXES"), ","),
			debug:              os.Getenv("GO_GITHUB_WIP_DEBUG") == "true",
		}
	}
	return config
}

// Manually instantiate Config. This bypasses configuration validation. Used for testing
func Set(appId int64, privateKeyFileName string, prefixes []string, debug bool) {
	config = &Config{
		appId:              appId,
		prefixes:           prefixes,
		privateKeyFileName: privateKeyFileName,
		debug:              debug,
	}
}
