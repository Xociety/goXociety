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
	"sync"
	"testing"
)

var uriGraphqlTest string

func init() {
	// please make sure goXociety server is up
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	uriGraphqlTest = "https://localhost:" + strconv.Itoa(globalConfig[env].ServerPort) + graphqlRoute
}

// query

// reaction
func TestGraphqlQueryReaction(t *testing.T) {
	q := []byte(`{ "query": "{ reaction { value reaction_id } }" }`)
	resp, err := http.Post(uriGraphqlTest, "Content-Type: application/json", bytes.NewBuffer(q))
	if err != nil {
		log.Panicln(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicln(err)
	}
	log.Println(string(b))
}

func TestGraphqlQueryCityPostCount(t *testing.T) {
	level := 1
	cityID := "AFG_1_1"
	str := `
	{
		"operationName":"",
		"variables":{"level":` + strconv.Itoa(level) + ` ,"city_id":"` + cityID + `"},
		"query": "query ($level: Int, $city_id: String) {\n city_post_count(level: $level, city_id: $city_id) {\n post_count\n city_id\n name\n }\n}\n"
	}`
	q := []byte(str)
	resp, err := http.Post(uriGraphqlTest, "Content-Type: application/json", bytes.NewBuffer(q))
	if err != nil {
		log.Panicln(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicln(err)
	}
	log.Println(string(b))
}

// place
func TestGraphqlQueryPlaceByLocation(t *testing.T) {
	lat := float64(24.9995389)
	lon := float64(121.498405)
	name := "天籟之音 音樂藝術中心"
	str := `
	{
		"operationName":"",
		"variables":{"lat":` + strconv.FormatFloat(lat, 'f', -1, 64) + ` ,"lon":` + strconv.FormatFloat(lon, 'f', -1, 64) + `, "name": "` + name + `", "page_token": ""},
		"query": "query ($lat: Float, $lon: Float, $name: String, $page_token: String) {\n place_by_location(lat: $lat, lon: $lon, name: $name, page_token: $page_token) {\n place{\n place_id\n name\n lat\n lon \n} next_page_token }\n}\n"
	}`
	q := []byte(str)
	resp, err := http.Post(uriGraphqlTest, "Content-Type: application/json", bytes.NewBuffer(q))
	if err != nil {
		log.Panicln(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicln(err)
	}
	log.Println(string(b))
}
func TestGraphqlQueryPlaceByName(t *testing.T) {
	q := []byte(`{ "query": "{ place_by_name(name:\"\", page_token:\"\") { place{place_id name lat lon} next_page_token } }" }`)
	resp, err := http.Post(uriGraphqlTest, "Content-Type: application/json", bytes.NewBuffer(q))
	if err != nil {
		log.Panicln(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicln(err)
	}
	log.Println(string(b))
}

// post
func TestGraphqlQueryPostsByPopular(t *testing.T) {
	q := []byte(`{ "query": "{ posts_by_popular (category_id:0, page:0) { post_id like_count dislike_count comment_count } }" }`)
	body := bytes.NewBuffer(q)
	client := &http.Client{}
	req, err := http.NewRequest("POST", uriGraphqlTest, body)
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

// mutation
// place
func TestGraphqlMutaionPlaceInsert(t *testing.T) {
	lat := float64(24.9995389)
	lon := float64(121.498405)
	name := "天籟之音 音樂藝術中心"
	str := `
	{
		"operationName":"",
		"variables":{"lat":` + strconv.FormatFloat(lat, 'f', -1, 64) + ` ,"lon": ` + strconv.FormatFloat(lon, 'f', -1, 64) + `, "name":"` + name + `"},
		"query": "mutation ($lat: Float, $lon: Float, $name: String) {\n place_insert(lat: $lat, lon: $lon, name: $name) {\n place_id\n country_code\n city_id_1\n city_id_2\n city_id_3\n city_id_4\n city_id_5\n lat\n lon\n name }\n}\n"
	}`
	q := []byte(str)
	body := bytes.NewBuffer(q)
	client := &http.Client{}
	req, err := http.NewRequest("POST", uriGraphqlTest, body)
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

// post
func TestGraphqlMutationPostInsert(t *testing.T) {
	// form-data
	q := `
		mutation {
			post_insert(
				type:1,
				origin_width:1920,
				origin_height:1280,
				content: "ya",
				category_id:1,
				place_id: 0,
			){
				post_id
			}
		}
	`
	params := make(map[string]string)
	params["query"] = q
	paramName := "file"
	// path := "./development/upload/sample/image.tar.gz"
	path := "./development/upload/sample/playlist.tar.gz"
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

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	file.Close()
	err = writer.Close()
	if err != nil {
		log.Panicln(err)
	}
	type reqLock struct {
		req   *http.Request
		mutex sync.Mutex
	}
	req, err := http.NewRequest("POST", uriGraphqlTest, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Token", "1")
	r := reqLock{
		req:   req,
		mutex: sync.Mutex{},
	}
	var wg sync.WaitGroup
	numReq := 100000
	for i := 0; i < numReq; i++ {
		wg.Add(1)
		go func(j int, rj *reqLock) {
			defer wg.Done()
			client := &http.Client{}
			rj.mutex.Lock()
			resp, err := client.Do(rj.req)
			rj.mutex.Unlock()
			_, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Panicln(err)
			}
			if j%100 == 0 {
				log.Println(j)
			}
			// log.Println(string(b))
		}(i, &r)
	}
	wg.Wait()
	log.Println("finished")
}
func TestGraphqlMutationPostPopularRead(t *testing.T) {
	for i := 0; i < 20; i++ {
		q := []byte(`{ "query": "mutation { post_popular_read (category_id:0, index_read: 0) { post_id like_count dislike_count comment_count } }" }`)
		body := bytes.NewBuffer(q)
		client := &http.Client{}
		req, err := http.NewRequest("POST", uriGraphqlTest, body)
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
