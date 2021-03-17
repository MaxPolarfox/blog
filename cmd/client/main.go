package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	uuid "github.com/satori/go.uuid"

	blogClient "github.com/MaxPolarfox/blog/pkg/client"
	"github.com/MaxPolarfox/blog/pkg/types"
	goToolsClient "github.com/MaxPolarfox/goTools/client"
)

const ServiceName = "blog"
const EnvironmentVariable = "APP_ENV"

func main() {
	// Load current environment
	env := os.Getenv(EnvironmentVariable)

	// load config options
	options := loadEnvironmentConfig(env)

	client := blogClient.NewBlogClient(options)
	defer client.Conn.Close()

	data := types.Blog{
		AuthorId: uuid.NewV4().String(),
		Title:    "My first Blog",
		Content:  "test content",
	}

	res, err := client.CreateBlog(context.Background(), data)
	if err != nil {
		log.Printf("error happened while calling CrateBlog: %v", err)
	} else {
		log.Printf("Response from server: %#v", *res)
	}
}

// loadEnvironmentConfig will use the environment string and concatenate to a proper config file to use
func loadEnvironmentConfig(env string) goToolsClient.Options {
	configFile := "config/" + ServiceName + "/client_" + env + ".json"
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		panic(err)
	}
	return parseConfigFile(configFile)
}

func parseConfigFile(configFile string) goToolsClient.Options {
	var opts goToolsClient.Options
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
