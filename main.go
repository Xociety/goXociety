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
	graphqlRoute             = "/graphql"
	graphqlGraphiql          = false
	graphqlHandlerPlayground = true
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
	postTypeMapID2Type         = make(map[int]string) // example: [0: "jpg", 1: "hls" ...]
	postTypeMapType2ID         = make(map[string]int) // example: [jpg: "0", "hls": 1 ...]
	reactionsMapID2Description = make(map[int]string) // example: [0: "like", 1: "dislike"]
	categoryMapID2Name         = make(map[int]string) // example: [0: "travel"]
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
	mongoDBXociety                      = "xociety"
	mongoCollectionCity                 = "city"
	mongoCollectionCountry              = "country"
	mongoCollectionCityLevel            = "city_level_"
	mongoCollectionPostPopular          = "post_popular"
	mongoCollectionPostPopularCommon    = "post_popular_common"
	mongoCollectionPostPopularReadIndex = "post_popular_read_index"
	mongoCollectionPostPopularRead      = "post_popular_read"
	mongoGeoNearSearchInKM              = 100 * 1000
	mongoTimeout                        = 0
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
	// postgres
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
	// mongo
	mongoSecret := make(map[string]secretMongo)
	if file, err := ioutil.ReadFile(globalConfig[env].MongoSecretFolderPath + globalConfig[env].MongoSecretAuthFilename); err != nil {
		log.Panicln("mongo secret file miss", err)
	} else {
		if err := json.Unmarshal(file, &mongoSecret); err != nil {
			log.Panicln("mongo secret file parse fail", err)
		}
	}
	globalSecret.Mongo = mongoSecret[env]
}
func initGCP() {
	// GCP
	// cloud storage
	clientOptionGoogleAPI = option.WithServiceAccountFile(globalConfig[env].GCPSecretFolderPath + globalConfig[env].GCPSecretFilename)
}
func initData() {
	var err error
	if countryConfigAPI, err = getCountries(); err != nil {
		log.Println("country config")
	}
	if languageConfigAPI, err = getLanguages(); err != nil {
		log.Fatalln("language config")
	}
	if postTypeConfigAPI, err = getPostType(); err == nil {
		IDs := []int{}
		values := []string{}
		for i := 0; i < len(postTypeConfigAPI); i++ {
			IDs = append(IDs, postTypeConfigAPI[i].PostTypeID)
			values = append(values, postTypeConfigAPI[i].Value)
		}
		convertIDAndValue(IDs, values, true, postTypeMapID2Type, postTypeMapType2ID)
	} else {
		log.Fatalln("post type config")
	}
	if reactionConfigAPI, err = getReaction(); err == nil {
		IDs := []int{}
		values := []string{}
		for i := 0; i < len(reactionConfigAPI); i++ {
			IDs = append(IDs, reactionConfigAPI[i].ReactionID)
			values = append(values, reactionConfigAPI[i].Value)
		}
		convertIDAndValue(IDs, values, false, reactionsMapID2Description, nil)
	} else {
		log.Println("reaction")
	}
	genderConfigAPI, err = getGender()
	if err != nil {
		log.Println("gender")
	}
	if categoryConfigAPI, err = getCategories(); err == nil {
		IDs := []int{}
		values := []string{}
		for i := 0; i < len(categoryConfigAPI); i++ {
			IDs = append(IDs, categoryConfigAPI[i].CategoryID)
			values = append(values, categoryConfigAPI[i].CategoryName)
		}
		convertIDAndValue(IDs, values, false, categoryMapID2Name, nil)
	} else {
		log.Println("category config")
	}
	log.Println(postTypeMapID2Type, postTypeMapType2ID)
	log.Println(reactionsMapID2Description)
	log.Println(categoryMapID2Name)
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
