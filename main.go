package main

import (
	"context"
	"encoding/json"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/chienfuchen32/handler"
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
	categorySup             = 0
)

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
	mongoDBXociety              = "xociety"
	mongoCollectionCity         = "city"
	mongoCollectionCountry      = "country"
	mongoCollectionCityLevel    = "city_level_"
	mongoCollectionPostCommon   = "post_common"
	mongoCollectionPostUserRead = "post_user_read"
	mongoGeoNearSearchInKM      = 100 * 1000
	mongoTimeout                = 0
	cityLevelRangeFirst         = 1
	cityLevelRangeLast          = 5
)

// Redis

const (
	redisDBPopularPostUserReadIndex          = 10
	redisHashCommonPopularPostUserReadIndex  = "common"
	redisHashCountryPopularPostUserReadIndex = "country"
	redisHashCityPopularPostUserReadIndex    = "city"
)

const numPopularPostPerRefresh = 10000
const numPerRequest = 10

// GCP cloud storage
var clientOptionGoogleAPI option.ClientOption

// Google map
var googleMapKey string

const radiusGoogleMap = 200

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

const (
	userTokenHeaderKey  = "User-Token"
	fileFormDataBodyKey = "file"
)

var (
	contextUserToken = contextKey(userTokenHeaderKey)
	contextKeyFile   = contextKey(fileFormDataBodyKey)
)

type contextKey string

func (c contextKey) String() string {
	return "xociety context key " + string(c)
}

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
	// google map
	googleMap := secretMap{}
	if file, err := ioutil.ReadFile(globalConfig[env].GCPSecretFolderPath + globalConfig[env].GoogleMapKeyFilename); err != nil {
		log.Panicln("google map secret file miss", err)
	} else {
		if err := json.Unmarshal(file, &googleMap); err != nil {
			log.Panicln("google map file parse fail", err)
		}
	}
	googleMapKey = googleMap.Key
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
}
func startServer() {
	// graphql
	http.Handle(graphqlRoute, func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// header user token
			reqCtx := context.Background()
			if userToken := r.Header.Get(userTokenHeaderKey); userToken != "" {
				reqCtx = context.WithValue(reqCtx, contextUserToken, userToken)
			}
			// file upload
			if file, _, err := r.FormFile(fileFormDataBodyKey); err == nil {
				reqCtx = context.WithValue(reqCtx, contextKeyFile, file)
			}
			inner.ServeHTTP(w, r.WithContext(reqCtx))
		})
	}(handler.New(&handler.Config{
		Schema:     &graphqlSchema,
		Pretty:     true,
		GraphiQL:   graphqlGraphiql,
		Playground: graphqlHandlerPlayground,
	})))
	// index
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./view/index.html")
	})
	// logo
	http.HandleFunc("/logo-background.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./view/logo-background.png")
	})
	// upload
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./view/upload.html")
	})
	http.HandleFunc("/upload/sample/image.tar.gz", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./development/upload/sample/image.tar.gz")
	})
	http.HandleFunc("/upload/sample/playlist.tar.gz", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./development/upload/sample/playlist.tar.gz")
	})
	server := &http.Server{
		Addr:           globalConfig[env].ServerAddrBind + ":" + strconv.Itoa(globalConfig[env].ServerPort),
		ReadTimeout:    5 * time.Minute,
		WriteTimeout:   5 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("xcociety graphql api server " + globalConfig[env].ServerAddrBind + ":" + strconv.Itoa(globalConfig[env].ServerPort))
	log.Fatal(server.ListenAndServeTLS(globalConfig[env].ServerSecretFolderPath+globalConfig[env].ServerSecretCertFilename, globalConfig[env].ServerSecretFolderPath+globalConfig[env].ServerSecretKeyFilename))
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
