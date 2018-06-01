package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/chienfuchen32/handler"
)

var (
	contextKeyFileOrigin       = contextKey("file origin")
	contextKeyFileOriginHeader = contextKey("file origin header")
	contextKeyFileSmall        = contextKey("file small")
)

type contextKey string

func (c contextKey) String() string {
	return "xociety context key " + string(c)
}

func startServer() {
	_ = importJSONDataFromGraphqlFile("./data/data.json", &graphqlSampleData)
	h := func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if readerOrigin, fileHeader, err := r.FormFile("file_origin"); err == nil {
				fileOrigin, err := ioutil.ReadAll(readerOrigin)
				if err != nil {
					log.Println("fileOrigin err", err)
				}
				log.Println("file", fileHeader.Filename, fileHeader.Header, fileHeader.Size, len(fileOrigin))
				fileOriginCtx := context.WithValue(context.Background(), contextKeyFileOrigin, fileOrigin)
				inner.ServeHTTP(w, r.WithContext(fileOriginCtx))
			} else {
				inner.ServeHTTP(w, r)
			}
		})
	}(handler.New(&handler.Config{
		Schema:   &graphqlSchema,
		Pretty:   true,
		GraphiQL: true,
	}))
	http.Handle(graphqlRoute, h)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	server := &http.Server{
		Addr:           ":" + strconv.Itoa(serverPort),
		ReadTimeout:    5 * time.Minute,
		WriteTimeout:   5 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(server.ListenAndServe())
}
