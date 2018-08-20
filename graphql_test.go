package main

import (
	"bytes"
	"crypto/tls"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestGraphqlQueryReaction(t *testing.T) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	q := []byte(`{ "query": "{ reaction { value reaction_id } }" }`)
	resp, err := http.Post("https://localhost:"+strconv.Itoa(globalConfig[env].ServerPort)+graphqlRoute, "Content-Type: application/json", bytes.NewBuffer(q))
	if err != nil {
		log.Panicln(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicln(err)
	}
	log.Println(string(b))
}

func TestGraphqlMutationPostInsert(t *testing.T) {
	q := `
		mutation {
			post_insert(
				country_id: 206,
				type:0,
				origin_width:1920,
				origin_height:1280,
				content: "ya",
				category_id:0,
			){
				post_id
			}
		}
	`
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	params := make(map[string]string)
	params["query"] = q
	paramName := "file"
	path := "./development/upload/sample/image.tar.gz"
	uri := "https://localhost:" + strconv.Itoa(globalConfig[env].ServerPort) + graphqlRoute
	for i := 0; i < 10000; i++ {
		file, err := os.Open(path)
		if err != nil {
			log.Panicln(err)
		}
		// defer file.Close()
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile(paramName, filepath.Base(path))
		if err != nil {
			log.Panicln(err)
		}
		_, err = io.Copy(part, file)

		file.Close()
		for key, val := range params {
			_ = writer.WriteField(key, val)
		}
		err = writer.Close()
		if err != nil {
			log.Panicln(err)
		}
		if i%100 == 0 {
			log.Println(i)
		}
		client := &http.Client{}
		req, err := http.NewRequest("POST", uri, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("User-Token", "1")
		resp, err := client.Do(req)
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("res", err)
			log.Panicln(err)
		}
		// log.Println(string(b))
	}
	log.Println("finished")
}

func TestGraphqlQueryPostsByPopular0(t *testing.T) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	q := []byte(`{ "query": "{ posts_by_popular_old (category_id:0, page:0) { post_id like_count dislike_count comment_count } }" }`)
	body := bytes.NewBuffer(q)
	uri := "https://localhost:" + strconv.Itoa(globalConfig[env].ServerPort) + graphqlRoute
	client := &http.Client{}
	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Token", "1")
	resp, err := client.Do(req)
	if err != nil {
		log.Panicln(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicln(err)
	}
	log.Println(string(b))
}
func TestGraphqlMutationPostsByPopular0(t *testing.T) {
	for i := 0; i < 20; i++ {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		q := []byte(`{ "query": "mutation { post_popular_read_old (category_id:0, index_read: 0) { post_id like_count dislike_count comment_count } }" }`)
		body := bytes.NewBuffer(q)
		uri := "https://localhost:" + strconv.Itoa(globalConfig[env].ServerPort) + graphqlRoute
		client := &http.Client{}
		req, err := http.NewRequest("POST", uri, body)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Token", "2")
		resp, err := client.Do(req)
		if err != nil {
			log.Panicln(err)
		}
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Panicln(err)
		}
		// log.Println(string(b))
	}
}
func TestGraphqlQueryPostsByPopular1(t *testing.T) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	q := []byte(`{ "query": "{ posts_by_popular (category_id:0, page:0) { post_id like_count dislike_count comment_count } }" }`)
	body := bytes.NewBuffer(q)
	uri := "https://localhost:" + strconv.Itoa(globalConfig[env].ServerPort) + graphqlRoute
	client := &http.Client{}
	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Token", "1")
	resp, err := client.Do(req)
	if err != nil {
		log.Panicln(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicln(err)
	}
	log.Println(string(b))
}

func TestGraphqlMutationPostsByPopular1(t *testing.T) {
	for i := 0; i < 20; i++ {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		q := []byte(`{ "query": "mutation { post_popular_read (category_id:0, index_read: 0) { post_id like_count dislike_count comment_count } }" }`)
		body := bytes.NewBuffer(q)
		uri := "https://localhost:" + strconv.Itoa(globalConfig[env].ServerPort) + graphqlRoute
		client := &http.Client{}
		req, err := http.NewRequest("POST", uri, body)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Token", "2")
		resp, err := client.Do(req)
		if err != nil {
			log.Panicln(err)
		}
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Panicln(err)
		}
		// log.Println(string(b))
	}
}
