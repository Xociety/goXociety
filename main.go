package main

import (
	"encoding/json"
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
	// numPerRequest := 10
)

// application
const (
	sixHoursInSecond        = 6 * 3600
	twelveHoursInSecond     = 12 * 3600
	twentyFourHoursInSecond = 24 * 3600
	sevenDaysInSecond       = 7 * 24 * 3600
	twoMonthsInSecond       = 2 * 30 * 24 * 3600
)

var postTypeMapID2Type = make(map[int]string) // example: [0: "jpg", 1: "hls" ...]
var postTypeMapType2ID = make(map[string]int) // example: [jpg: "0", "hls": 1 ...]

// Postgres
const postgresConStr = "host=localhost port=31160 user=postgres password=mysecretpassword sslmode=disable"

// GCP
var clientAuthGCP clientAuthFromServiceAccountFileGCP

const clientAuthGCPFilePath = "./keyfileGCP.json"

// GCP cloud storage
var clientOptionGoogleAPI option.ClientOption

const (
	bucketRootCloudStorage   = "storage.1mthechildbride.com" // "storagejp.1mthechildbride.com", "storagetw.1mthechildbride.com", "storageasia.1mthechildbride.com"
	bucketImagesCloudStorage = "images"
	bucketVideosCloudStorage = "videos"
)

const (
	mediaFormatJPG  = "jpg"
	mediaFormatHLS  = "hls"
	mediaFormatM3U8 = "m3u8"
	mediaFormatTS   = "ts"
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
	// config
	loadPostTypeFromFConfig("./database/post_type.csv")
}

func main() {
	startServer()
}
