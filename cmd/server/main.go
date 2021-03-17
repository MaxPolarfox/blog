package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/MaxPolarfox/blog/pkg/blog"
	"github.com/MaxPolarfox/blog/pkg/types"
	"github.com/MaxPolarfox/goTools/mongoDB"
)

const ServiceName = "blog"
const EnvironmentVariable = "APP_ENV"

func main() {
	// Load current environment
	env := os.Getenv(EnvironmentVariable)

	// load config options
	options := loadEnvironmentConfig(env)

	// DB
	blogCollection := mongoDB.NewMongo(options.DB.Blog)

	s := blog.NewService(options, blogCollection)

	s.Start()
}

// loadEnvironmentConfig will use the environment string and concatenate to a proper config file to use
func loadEnvironmentConfig(env string) types.Options {
	configFile := "config/" + ServiceName + "/server_" + env + ".json"
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		panic(err)
	}
	return parseConfigFile(configFile)
}

func parseConfigFile(configFile string) types.Options {
	var opts types.Options
	byts, err := ioutil.ReadFile(configFile)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(byts, &opts)
	if err != nil {
		panic(err)
	}

	return opts
}
