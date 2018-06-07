package main

import (
	"encoding/json"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"

	"google.golang.org/api/option"
)

// Graphql Server
const (
	serverPort   = 4000
	graphqlRoute = "/graphql"
	graphiql     = true
)

// Postgres
const postgresConStr = "host=localhost port=31160 user=postgres password=mysecretpassword sslmode=disable"

// GCP
var clientAuthGCP clientAuthFromServiceAccountFileGCP

const clientAuthGCPFilePath = "./keyfileGCP.json"

// GCP cloud storage
var clientOptionGoogleAPI option.ClientOption

const (
	bucketRootCloudStorage   = "storage.1mthechildbride.com"
	bucketImagesCloudStorage = "images"
	bucketVideosCloudStorage = "videos"
)

const (
	postVideoDefaultName = "video.m3u8"
)

func init() {
	// GCP
	if file, err := ioutil.ReadFile(clientAuthGCPFilePath); err != nil {
		log.Panicln("gcp auth file miss", err)
	} else {
		if err := json.Unmarshal(file, &clientAuthGCP); err != nil {
			log.Panicln("gcp auth file parse fail", err)
		}
	}
	// cloud storage
	clientOptionGoogleAPI = option.WithServiceAccountFile(clientAuthGCPFilePath)
}

func main() {
	fmt.Println("hello xcociety")
	// err := startInsertXuserfaker()
	// fmt.Println("finished", err)

	startServer()
}
