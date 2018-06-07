package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/chienfuchen32/goXHandler"
)

var (
	contextKeyFile = contextKey("file")
)

type contextKey string

func (c contextKey) String() string {
	return "xociety context key " + string(c)
}

func startServer() {
	h := func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// header user token
			if file, _, err := r.FormFile("file"); err == nil {
				fileCtx := context.WithValue(context.Background(), contextKeyFile, file)
				inner.ServeHTTP(w, r.WithContext(fileCtx))
			} else {
				inner.ServeHTTP(w, r)
			}
		})
	}(handler.New(&handler.Config{
		Schema:   &graphqlSchema,
		Pretty:   true,
		GraphiQL: graphiql,
	}))
	http.Handle(graphqlRoute, h)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	server := &http.Server{
		Addr:           "127.0.0.1:" + strconv.Itoa(serverPort),
		ReadTimeout:    5 * time.Minute,
		WriteTimeout:   5 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(server.ListenAndServe())
}
