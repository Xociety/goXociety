package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	// "github.com/graphql-go/handler"
	"github.com/chienfuchen32/handler"
)

func customHandler(h *handler.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ctx := context.Context.NewContext(r)
		// log.Debugf(ctx, "Request came: %v", r.RequestURI)
		// log.Println("w", w, "r", r)
		// Call context handler to serve HTTP
		// h.ContextHandler(ctx, w, r)
	})
}
func startServer() {
	_ = importJSONDataFromGraphqlFile("./data/data.json", &graphqlSampleData)
	// http.HandleFunc(graphqlRoute, func(w http.ResponseWriter, r *http.Request) {
	// 	user := struct {
	// 		ID   int    `json:"id"`
	// 		Name string `json:"name"`
	// 	}{1, "cool user"}
	// 	// reqStr := r.URL.Query().Get("query")
	// 	// if r.Method == "POST" {
	// 	// 	if body, err := ioutil.ReadAll(r.Body); err == nil {
	// 	// 		reqStr = string(body)
	// 	// 	}
	// 	// }
	// 	result := graphql.Do(graphql.Params{
	// 		Schema: graphqlSchema,
	// 		// RequestString: reqStr,
	// 		Context: context.WithValue(context.Background(), "currentUser", user),
	// 	})
	// 	// graphql.
	// 	if len(result.Errors) > 0 {
	// 		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	// 	}
	// 	json.NewEncoder(w).Encode(result)
	// })
	// handleGraphiql := handler.New(&handler.Config{
	// 	Schema:   &graphqlSchema,
	// 	Pretty:   true,
	// 	GraphiQL: true,
	// })
	// http.Handle(graphiqlRoute, handleGraphiql)

	h := func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}{1, "cool user"}
			// some kind of authentication here
			// user := parseJwtUser(r)
			// if !user.canDoGraphQL() {
			// 	// do something to tell the user
			// 	// ...
			// 	// ...
			// }
			// do normal graphql
			innerCtx := context.WithValue(r.Context(), "currentUser", user)
			// w.Header().Set("yo", "yo")
			outterCtx := context.WithValue(context.Background(), "yo", "yo")
			// log.Println(r.WithContext(innerCtx).WithContext(outterCtx))
			// var fileCtx context.Context
			if file, fileHeader, err := r.FormFile("file"); err == nil {
				if blob, err := ioutil.ReadAll(file); err == nil {
					// log.Println("query", r.FormValue("query"))
					// fileCtx = context.WithValue(context.Background(), "file", blob)
					log.Println("file", fileHeader.Filename, fileHeader.Header, fileHeader.Size, len(blob))
					// ioutil.WriteFile("./development/Weekly Milestone.pages", blob, 0644)
				}
			}
			inner.ServeHTTP(w, r.WithContext(innerCtx).WithContext(outterCtx)) //.WithContext(fileCtx))
		})
	}(handler.New(&handler.Config{
		Schema:   &graphqlSchema,
		Pretty:   true,
		GraphiQL: true,
	}))
	http.Handle(graphqlRoute, h)
	server := &http.Server{
		Addr:           ":" + strconv.Itoa(serverPort),
		ReadTimeout:    5 * time.Minute,
		WriteTimeout:   5 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(server.ListenAndServe())
}
