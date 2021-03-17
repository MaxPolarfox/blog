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

	ctx := context.Background()

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

	createdRes, err := client.CreateBlog(ctx, data)
	if err != nil {
		log.Printf("Create blog error: %v", err)
	} else {
		log.Printf("Blog was creates: %#v", *createdRes)
	}

	// Should return blog
	_, err = client.ReadBlog(ctx, "123456789")
	if err != nil {
		log.Printf("Read blog error: %v", err)
	}

	// Should return not found
	_, err = client.ReadBlog(ctx, *createdRes)
	if err != nil {
		log.Printf("Read blog error: %v", err)
	}

	// Should return blog
	readRes, err := client.ReadBlog(ctx, *createdRes)
	if err != nil {
		log.Printf("Read blog error: %v", err)
	} else {
		log.Printf("Blog was read: %#v", *readRes)
	}

	// Should return slice of blogs
	listRes, err := client.ListBlog(ctx)
	if err != nil {
		log.Printf("List blog error: %v", err)
	} else {
		log.Printf("Listed blogs: %#v", listRes)
	}

	// Should return not found
	readRes.Content = "Updated content"
	err = client.UpdateBlog(ctx, types.Blog{ID: "12345"})
	if err != nil {
		log.Printf("Update blog error: %v", err)
	}

	// Should successfully update blog
	readRes.Content = "Updated content"
	err = client.UpdateBlog(ctx, *readRes)
	if err != nil {
		log.Printf("Update blog error: %v", err)
	} else {
		log.Printf("Successfully updated blog: %v", readRes.ID)
	}

	// Should successfully delete blog
	deleteRes, err := client.DeleteBlog(ctx, readRes.ID)
	if err != nil {
		log.Printf("Delete blog error: %v", err)
	} else {
		log.Printf("Successfully deleted blog: %v", *deleteRes)
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
