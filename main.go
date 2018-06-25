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
)

// application
const (
	sixHoursInSecond        = 6 * 3600
	twelveHoursInSecond     = 12 * 3600
	twentyFourHoursInSecond = 24 * 3600
	sevenDaysInSecond       = 7 * 24 * 3600
	twoMonthsInSecond       = 2 * 30 * 24 * 3600
)

type forwardLookupMap map[int]string
type reverseLookupMap map[string]int

var (
	postTypeMapID2Type                 = make(map[int]string) // example: [0: "jpg", 1: "hls" ...]
	postTypeMapType2ID                 = make(map[string]int) // example: [jpg: "0", "hls": 1 ...]
	reactionsTypeMapID2Description     = make(map[int]string) // example: [0: "like", 1: "dislike"]
	countryTypeMapID2Country           = make(map[int]string) // example: [0: "Afghanistan"]
	countryTypeMapID2CountryCode       = make(map[int]string) // example: [0: "af"]
	languageTypeMapID2DisplayLanguage  = make(map[int]string) // example: [0: "Afrikaans"]
	languageTypeMapID2HlParameterValue = make(map[int]string) // example: [0: "af"]
	genderTypeMapID2Description        = make(map[int]string) // example: [0: "not known"]
)

// Postgres
const postgresConStr = "host=localhost port=30749 user=postgres password=mysecretpassword sslmode=disable"

// const postgresConStr = "host=my-release-postgresql port=5432 user=postgres password=MGmQClLFup sslmode=disable"

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
	loadConfigFromFile("./database/country.csv", 0, 1, 3, false, countryTypeMapID2Country, nil)
	loadConfigFromFile("./database/country.csv", 0, 2, 3, false, countryTypeMapID2CountryCode, nil)
	loadConfigFromFile("./database/language.csv", 0, 1, 3, false, languageTypeMapID2DisplayLanguage, nil)
	loadConfigFromFile("./database/language.csv", 0, 2, 3, false, languageTypeMapID2HlParameterValue, nil)
	loadConfigFromFile("./database/gender.csv", 0, 1, 2, false, genderTypeMapID2Description, nil)
	loadConfigFromFile("./database/reaction.csv", 0, 1, 2, false, reactionsTypeMapID2Description, nil)
	loadConfigFromFile("./database/post_type.csv", 0, 1, 2, true, postTypeMapID2Type, postTypeMapType2ID)

	// api
	for k, v := range countryTypeMapID2Country {
		countryConfigAPI = append(countryConfigAPI, countryAPI{
			CountryID:   k,
			Country:     v,
			CountryCode: countryTypeMapID2CountryCode[k],
		})
	}
	for k, v := range languageTypeMapID2DisplayLanguage {
		languageConfigAPI = append(languageConfigAPI, languageAPI{
			LanguageID:      k,
			DisplayLanguage: v,
			Value:           languageTypeMapID2HlParameterValue[k],
		})
	}
	for k, v := range genderTypeMapID2Description {
		genderConfigAPI = append(genderConfigAPI, genderAPI{
			GenderID: k,
			Value:    v,
		})
	}
	for k, v := range reactionsTypeMapID2Description {
		reactionConfigAPI = append(reactionConfigAPI, reactionAPI{
			ReactionID: k,
			Value:      v,
		})
	}
	for k, v := range postTypeMapID2Type {
		postTypeConfigAPI = append(postTypeConfigAPI, postTypeAPI{
			PostTypeID: k,
			Value:      v,
		})
	}
}

func main() {
	startServer()
}
