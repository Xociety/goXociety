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
	hostname = "127.0.0.1"
	// hostname     = "0.0.0.0"
	serverPort   = 4000
	graphqlRoute = "/graphql"
	graphiql     = true
	// numPerRequest := 10
)

// api
var (
	countryConfigAPI  []countryAPI
	languageConfigAPI []languageAPI
	genderConfigAPI   []genderAPI
	reactionConfigAPI []reactionAPI
	postTypeConfigAPI []postTypeAPI
	categoryConfigAPI []categoryAPI
)

// application
const (
	sixHoursInSecond        = 6 * 60 * 60
	twelveHoursInSecond     = 12 * 60 * 60
	twentyFourHoursInSecond = 24 * 60 * 60
	sevenDaysInSecond       = 7 * 24 * 60 * 60
	twoMonthsInSecond       = 2 * 30 * 24 * 60 * 60
)

type forwardLookupMap map[int]string
type reverseLookupMap map[string]int

var (
	postTypeMapID2Type             = make(map[int]string) // example: [0: "jpg", 1: "hls" ...]
	postTypeMapType2ID             = make(map[string]int) // example: [jpg: "0", "hls": 1 ...]
	reactionsMapID2Description     = make(map[int]string) // example: [0: "like", 1: "dislike"]
	countryMapID2Country           = make(map[int]string) // example: [0: "Afghanistan"]
	countryMapID2CountryCode       = make(map[int]string) // example: [0: "af"]
	languageMapID2DisplayLanguage  = make(map[int]string) // example: [0: "Afrikaans"]
	languageMapID2HlParameterValue = make(map[int]string) // example: [0: "af"]
	genderMapID2Description        = make(map[int]string) // example: [0: "not known"]
	categoryMapID2Name             = make(map[int]string) // example: [0: "travel"]
)

// PostgreSQL
const (
	postgresConStr     = "host=localhost port=30749 user=postgres password=mysecretpassword sslmode=disable"
	dbFilePathPostgres = "./database/postgres/"

// const postgresConStr = "host=my-release-postgresql port=5432 user=postgres password=mysecretpassword sslmode=disable"
)

// MongoDB
const (
	mongoConStr                = "localhost:30668"
	mongoDBXociety             = "xociety"
	mongoCollectionPostPopular = "post_popular"
	mongoCollectionPostRead    = "post_read"
)

const numPopularPostPerRefresh = 10000

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
	loadConfigFromFile(dbFilePathPostgres+"country.csv", 0, 1, 3, false, countryMapID2Country, nil)
	loadConfigFromFile(dbFilePathPostgres+"country.csv", 0, 2, 3, false, countryMapID2CountryCode, nil)
	loadConfigFromFile(dbFilePathPostgres+"language.csv", 0, 1, 3, false, languageMapID2DisplayLanguage, nil)
	loadConfigFromFile(dbFilePathPostgres+"language.csv", 0, 2, 3, false, languageMapID2HlParameterValue, nil)
	loadConfigFromFile(dbFilePathPostgres+"gender.csv", 0, 1, 2, false, genderMapID2Description, nil)
	loadConfigFromFile(dbFilePathPostgres+"reaction.csv", 0, 1, 2, false, reactionsMapID2Description, nil)
	loadConfigFromFile(dbFilePathPostgres+"post_type.csv", 0, 1, 2, true, postTypeMapID2Type, postTypeMapType2ID)
	loadConfigFromFile(dbFilePathPostgres+"category.csv", 0, 1, 2, false, categoryMapID2Name, nil)

	// api
	for k, v := range countryMapID2Country {
		countryConfigAPI = append(countryConfigAPI, countryAPI{
			CountryID:   k,
			Country:     v,
			CountryCode: countryMapID2CountryCode[k],
		})
	}
	for k, v := range languageMapID2DisplayLanguage {
		languageConfigAPI = append(languageConfigAPI, languageAPI{
			LanguageID:      k,
			DisplayLanguage: v,
			Value:           languageMapID2HlParameterValue[k],
		})
	}
	for k, v := range genderMapID2Description {
		genderConfigAPI = append(genderConfigAPI, genderAPI{
			GenderID: k,
			Value:    v,
		})
	}
	for k, v := range categoryMapID2Name {
		categoryConfigAPI = append(categoryConfigAPI, categoryAPI{
			CategoryID:   k,
			CategoryName: v,
		})
	}
	for k, v := range reactionsMapID2Description {
		reactionConfigAPI = append(reactionConfigAPI, reactionAPI{
			ReactionID: k,
			Value:      v,
		})
	}
	for k, v := range postTypeMapID2Type {
		postType := postTypeAPI{
			PostTypeID: k,
			Value:      v,
		}
		switch v {
		case mediaFormatJPG:
			postType.FileFormat = []string{mediaFormatJPG}
		case mediaFormatHLS:
			postType.FileFormat = []string{mediaFormatM3U8, mediaFormatTS}
		}
		postTypeConfigAPI = append(postTypeConfigAPI, postType)
	}
}

func main() {
	startServer()
}
