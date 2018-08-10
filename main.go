package main

import (
	"encoding/json"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"

	"google.golang.org/api/option"
)

var globalConfig map[string]config
var globalSecret secret
var env = "development"

// Graphql Server
const (
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
// all table name
)

var (
	postgresConStr = ""
)

// MongoDB
const (
	mongoDBXociety             = "xociety"
	mongoCollectionPostPopular = "post_popular"
	mongoCollectionPostRead    = "post_read"
)

const numPopularPostPerRefresh = 10000

// GCP cloud storage
var clientOptionGoogleAPI option.ClientOption

const (
	bucketPostCloudStorage   = "posts"
	bucketImagesCloudStorage = bucketPostCloudStorage + "/images"
	bucketVideosCloudStorage = bucketPostCloudStorage + "/videos"
	bucketUserCloudStorage   = "users"
)

const (
	mediaFormatJPG  = "jpg"
	mediaFormatHLS  = "hls"
	mediaFormatM3U8 = "m3u8"
	mediaFormatTS   = "ts"
)

func initSecret() {
	// secret
	postgresSecret := make(map[string]secretPostgres)
	if file, err := ioutil.ReadFile(globalConfig[env].PostgresSecretFolderPath + globalConfig[env].PostgresSecretAuthFilename); err != nil {
		log.Panicln("postgres secret file miss", err)
	} else {
		if err := json.Unmarshal(file, &postgresSecret); err != nil {
			log.Panicln("postgres secret file parse fail", err)
		}
	}
	globalSecret.Postgres = postgresSecret[env]
	postgresConStr = globalConfig[env].PostgresConStr + " " + globalSecret.Postgres.PostgresAuthStr
}
func initGCP() {
	// GCP
	// cloud storage
	clientOptionGoogleAPI = option.WithServiceAccountFile(globalConfig[env].GCPSecretFolderPath + globalConfig[env].GCPSecretFilename)
}
func initData() {

	loadConfigFromFile(globalConfig[env].PostgresConfigFolderPath+"country.csv", 0, 1, 3, false, countryMapID2Country, nil)
	loadConfigFromFile(globalConfig[env].PostgresConfigFolderPath+"country.csv", 0, 2, 3, false, countryMapID2CountryCode, nil)
	loadConfigFromFile(globalConfig[env].PostgresConfigFolderPath+"language.csv", 0, 1, 3, false, languageMapID2DisplayLanguage, nil)
	loadConfigFromFile(globalConfig[env].PostgresConfigFolderPath+"language.csv", 0, 2, 3, false, languageMapID2HlParameterValue, nil)
	loadConfigFromFile(globalConfig[env].PostgresConfigFolderPath+"gender.csv", 0, 1, 2, false, genderMapID2Description, nil)
	loadConfigFromFile(globalConfig[env].PostgresConfigFolderPath+"reaction.csv", 0, 1, 2, false, reactionsMapID2Description, nil)
	loadConfigFromFile(globalConfig[env].PostgresConfigFolderPath+"post_type.csv", 0, 1, 2, true, postTypeMapID2Type, postTypeMapType2ID)
	loadConfigFromFile(globalConfig[env].PostgresConfigFolderPath+"category.csv", 0, 1, 2, false, categoryMapID2Name, nil)

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

func init() {
	/*
		init steps:
		* env
		* config
		* secret
		* others
	*/
	// config
	configFolerPath := "./config"
	switch os.Getenv("env") {
	case "development":
		// default env, configFolerPath
	case "staging":
		env = os.Getenv("env")
		// undefined
	case "production":
		env = os.Getenv("env")
	}
	if file, err := ioutil.ReadFile(configFolerPath + "/xocietyConfig.json"); err != nil {
		log.Panicln("xociety config file miss", err)
	} else {
		if err := json.Unmarshal(file, &globalConfig); err != nil {
			log.Panicln("xociety config file parse fail", err)
		}
	}
	initSecret()
	initGCP()
	initData()
}

func main() {
	startServer()
}
