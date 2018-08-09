package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/chienfuchen32/goXHandler"
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

func startServer() {
	h := func(inner http.Handler) http.Handler {
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
		Schema:   &graphqlSchema,
		Pretty:   true,
		GraphiQL: graphiql,
	}))
	http.Handle(graphqlRoute, h)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./view/index.html")
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
