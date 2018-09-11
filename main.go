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

	"github.com/chienfuchen32/goXociety/x/common"
	"github.com/chienfuchen32/goXociety/x/config"
	"github.com/chienfuchen32/goXociety/x/io"
	"github.com/chienfuchen32/handler"
	"google.golang.org/api/option"
)

func initSecret() {
	// secret
	// postgres
	postgresSecret := make(map[string]config.SecretPostgres)
	if file, err := ioutil.ReadFile(config.GlobalConfig[config.Env].PostgresSecretFolderPath + config.GlobalConfig[config.Env].PostgresSecretAuthFilename); err != nil {
		log.Panicln("postgres secret file miss", err)
	} else {
		if err := json.Unmarshal(file, &postgresSecret); err != nil {
			log.Panicln("postgres secret file parse fail", err)
		}
	}
	config.GlobalSecret.Postgres = postgresSecret[config.Env]
	config.PostgresConStr = config.GlobalConfig[config.Env].PostgresConStr + " " + config.GlobalSecret.Postgres.PostgresAuthStr
	// mongo
	mongoSecret := make(map[string]config.SecretMongo)
	if file, err := ioutil.ReadFile(config.GlobalConfig[config.Env].MongoSecretFolderPath + config.GlobalConfig[config.Env].MongoSecretAuthFilename); err != nil {
		log.Panicln("mongo secret file miss", err)
	} else {
		if err := json.Unmarshal(file, &mongoSecret); err != nil {
			log.Panicln("mongo secret file parse fail", err)
		}
	}
	config.GlobalSecret.Mongo = mongoSecret[config.Env]
}
func initGCP() {
	// GCP
	// cloud storage
	config.ClientOptionGoogleAPI = option.WithServiceAccountFile(config.GlobalConfig[config.Env].GCPSecretFolderPath + config.GlobalConfig[config.Env].GCPSecretFilename)
	// google map
	googleMap := config.SecretMap{}
	if file, err := ioutil.ReadFile(config.GlobalConfig[config.Env].GCPSecretFolderPath + config.GlobalConfig[config.Env].GoogleMapKeyFilename); err != nil {
		log.Panicln("google map secret file miss", err)
	} else {
		if err := json.Unmarshal(file, &googleMap); err != nil {
			log.Panicln("google map file parse fail", err)
		}
	}
	config.GoogleMapKey = googleMap.Key
}
func initData() {
	var err error
	c, err := io.ConnectPostgres()
	if err != nil {
		log.Fatalln(err)
	}
	defer c.DB.Close()
	cm, err := io.ConnectMongoDB()
	if err != nil {
		log.Fatalln(err)
	}
	defer cm.Session.Close()
	if config.CountryConfigAPI, err = io.GetCountries(&cm); err != nil {
		log.Println("country config")
	}
	if config.LanguageConfigAPI, err = io.GetLanguages(&c); err != nil {
		log.Fatalln("language config")
	}
	if config.PostTypeConfigAPI, err = io.GetPostType(&c); err == nil {
		IDs := []int{}
		values := []string{}
		for i := 0; i < len(config.PostTypeConfigAPI); i++ {
			IDs = append(IDs, config.PostTypeConfigAPI[i].PostTypeID)
			values = append(values, config.PostTypeConfigAPI[i].Value)
		}
		common.ConvertIDAndValue(IDs, values, true, config.PostTypeMapID2Type, config.PostTypeMapType2ID)
	} else {
		log.Fatalln("post type config")
	}
	if config.ReactionConfigAPI, err = io.GetReaction(&c); err == nil {
		IDs := []int{}
		values := []string{}
		for i := 0; i < len(config.ReactionConfigAPI); i++ {
			IDs = append(IDs, config.ReactionConfigAPI[i].ReactionID)
			values = append(values, config.ReactionConfigAPI[i].Value)
		}
		common.ConvertIDAndValue(IDs, values, false, config.ReactionsMapID2Description, nil)
	} else {
		log.Println("reaction")
	}
	config.GenderConfigAPI, err = io.GetGender(&c)
	if err != nil {
		log.Println("gender")
	}
	if config.CategoryConfigAPI, err = io.GetCategories(&c); err == nil {
		IDs := []int{}
		values := []string{}
		for i := 0; i < len(config.CategoryConfigAPI); i++ {
			IDs = append(IDs, config.CategoryConfigAPI[i].CategoryID)
			values = append(values, config.CategoryConfigAPI[i].CategoryName)
		}
		common.ConvertIDAndValue(IDs, values, false, config.CategoryMapID2Name, nil)
	} else {
		log.Println("category config")
	}
}
func startServer() {
	// graphql
	http.Handle(config.GraphqlRoute, func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// header user token
			reqCtx := context.Background()
			if userToken := r.Header.Get(config.UserTokenHeaderKey); userToken != "" {
				reqCtx = context.WithValue(reqCtx, config.ContextUserToken, userToken)
			}
			// file upload
			if file, _, err := r.FormFile(config.FileFormDataBodyKey); err == nil {
				reqCtx = context.WithValue(reqCtx, config.ContextKeyFile, file)
			}
			inner.ServeHTTP(w, r.WithContext(reqCtx))
		})
	}(handler.New(&handler.Config{
		Schema:     &io.GraphqlSchema,
		Pretty:     true,
		GraphiQL:   config.GraphqlGraphiql,
		Playground: config.GraphqlHandlerPlayground,
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
		Addr:           config.GlobalConfig[config.Env].ServerAddrBind + ":" + strconv.Itoa(config.GlobalConfig[config.Env].ServerPort),
		ReadTimeout:    5 * time.Minute,
		WriteTimeout:   5 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("xcociety graphql api server " + config.GlobalConfig[config.Env].ServerAddrBind + ":" + strconv.Itoa(config.GlobalConfig[config.Env].ServerPort))
	log.Fatal(server.ListenAndServeTLS(config.GlobalConfig[config.Env].ServerSecretFolderPath+config.GlobalConfig[config.Env].ServerSecretCertFilename, config.GlobalConfig[config.Env].ServerSecretFolderPath+config.GlobalConfig[config.Env].ServerSecretKeyFilename))
}

func init() {
	/*
		init steps:
		* config.Env
		* config
		* secret
		* others
	*/
	// config
	configFolerPath := "./config"
	switch os.Getenv("config.Env") {
	case "development":
		// default config.Env, configFolerPath
	case "staging":
		config.Env = os.Getenv("config.Env")
		// undefined
	case "production":
		config.Env = os.Getenv("config.Env")
	}
	if file, err := ioutil.ReadFile(configFolerPath + "/xocietyConfig.json"); err != nil {
		log.Panicln("xociety config file miss", err)
	} else {
		if err := json.Unmarshal(file, &config.GlobalConfig); err != nil {
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
