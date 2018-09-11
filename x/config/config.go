package config

import (
	"math/rand"

	"google.golang.org/api/option"
)

var GlobalConfig map[string]Config
var GlobalSecret Secret
var Env = "development"

// Graphql Server
const (
	GraphqlRoute             = "/graphql"
	GraphqlGraphiql          = false
	GraphqlHandlerPlayground = true
	// config.NumPerRequest := 10
)

// api
var (
	CountryConfigAPI  []CityAPI
	LanguageConfigAPI []LanguageAPI
	GenderConfigAPI   []GenderAPI
	ReactionConfigAPI []ReactionAPI
	PostTypeConfigAPI []PostTypeAPI
	CategoryConfigAPI []CategoryAPI
)

// application
const (
	SixHoursInSecond        = 6 * 60 * 60
	TwelveHoursInSecond     = 12 * 60 * 60
	TwentyFourHoursInSecond = 24 * 60 * 60
	SevenDaysInSecond       = 7 * 24 * 60 * 60
	TwoMonthsInSecond       = 2 * 30 * 24 * 60 * 60
	CategorySup             = 0
)

var (
	PostTypeMapID2Type         = make(map[int]string) // example: [0: "jpg", 1: "hls" ...]
	PostTypeMapType2ID         = make(map[string]int) // example: [jpg: "0", "hls": 1 ...]
	ReactionsMapID2Description = make(map[int]string) // example: [0: "like", 1: "dislike"]
	CategoryMapID2Name         = make(map[int]string) // example: [0: "travel"]
)

// PostgreSQL
const (
// all table name
)

var (
	PostgresConStr = ""
)

// MongoDB
const (
	MongoDBXociety                          = "xociety"
	MongoCollectionCityGeometry             = "city_geometry"
	MongoCollectionCity                     = "city"
	MongoCollectionPostCommon               = "post_common"
	MongoCollectionPostUserRead             = "post_user_read"
	MongoCollectionPopularPostUserReadIndex = "popular_post_user_read_index"
	MongoGeoNearSearchInKM                  = 100 * 1000
	MongoTimeout                            = 0
	CityLevelRangeFirst                     = 1
	CityLevelRangeLast                      = 5
)

// Redis

const (
	RedisDBPopularPostUserReadIndex          = 10
	RedisHashCommonPopularPostUserReadIndex  = "common"
	RedisHashCountryPopularPostUserReadIndex = "country"
	RedisHashCityPopularPostUserReadIndex    = "city"
)

const NumPopularPostPerRefresh = 10000
const NumPerRequest = 10

// GCP cloud storage
var ClientOptionGoogleAPI option.ClientOption

// Google map
var GoogleMapKey string

// random seed
var RandSeed = rand.New(rand.NewSource(99))

const RadiusGoogleMap = 200

const (
	BucketPostCloudStorage   = "posts"
	BucketImagesCloudStorage = BucketPostCloudStorage + "/images"
	BucketVideosCloudStorage = BucketPostCloudStorage + "/videos"
	BucketUserCloudStorage   = "users"
)

const (
	MediaFormatJPG  = "jpg"
	MediaFormatHLS  = "hls"
	MediaFormatM3U8 = "m3u8"
	MediaFormatTS   = "ts"
)

const (
	UserTokenHeaderKey  = "User-Token"
	FileFormDataBodyKey = "file"
)

var (
	ContextUserToken = contextKey(UserTokenHeaderKey)
	ContextKeyFile   = contextKey(FileFormDataBodyKey)
)

type contextKey string

func (c contextKey) String() string {
	return "xociety context key " + string(c)
}
