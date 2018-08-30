package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/go-redis/redis"
	geo "github.com/kellydunn/golang-geo"
	"github.com/manveru/faker"

	_ "github.com/lib/pq"
)

func genXuserFaker(fake *faker.Faker, r *rand.Rand) userDB {
	timeNow := getNowUnixTimestamp()
	return userDB{
		Username:    fake.UserName(),
		Email:       fake.Email(),
		Password:    "salted",
		Name:        fake.Name(),
		Phone:       fake.PhoneNumber(),
		Gender:      r.Intn(3),
		Bio:         fake.Sentence(r.Intn(3), true),
		Credit:      0,
		PhotoURL:    "",
		LanguageID:  13,
		CountryCode: "TWN",
		Timezone:    28800,
		LastIP:      fake.IPv4Address().String(),
		Updatetime:  timeNow,
		Createtime:  timeNow,
	}
}

func genPostFaker(userID, placeID int64, faker *faker.Faker, r *rand.Rand) postAPI {
	timeNow := getNowUnixTimestamp()
	postType := 0             // r.Intn(2)
	originWidth := 0          // plz check current sample size
	originHeight := 0         // plz check current sample size
	categoryID := r.Intn(16)  // plz check current category
	sampleID := r.Intn(7) + 1 // plz check current sample file
	if categoryID == 0 {
		sampleID = 1
	}
	switch postTypeMapID2Type[postType] {
	case mediaFormatJPG:
		originWidth = 1920
		originHeight = 1280
		if categoryID != 0 {
			originWidth = 1242
			originHeight = 2004
		}
	case mediaFormatHLS:
		originWidth = 1920
		originHeight = 1080
	}
	post := postAPI{
		User: userBasicAPI{
			UserID: userID,
		},
		Content: "test post",
		Blob: blobAPI{
			BlobID:       "sample/" + strconv.Itoa(categoryID) + "/" + strconv.Itoa(sampleID),
			OriginWidth:  originWidth,
			OriginHeight: originHeight,
		},
		Type:         postType,
		LikeCount:    0,
		DislikeCount: 0,
		CommentCount: 0,
		Place:        placeAPI{PlaceID: placeID},
		CategoryID:   categoryID,
		Public:       true,
		Createtime:   timeNow,
		Updatetime:   timeNow,
	}
	return post
}

func genPostReactionFaker(postID, userID int64, r *rand.Rand) reactionOnPostAPI {
	return reactionOnPostAPI{
		PostID:     postID,
		User:       userBasicAPI{UserID: userID},
		ReactionID: r.Intn(2),
		Createtime: getNowUnixTimestamp(),
	}
}

func genPostCommentFaker(postID, userID int64, faker *faker.Faker, r *rand.Rand) commentAPI {
	timestamp := getNowUnixTimestamp()
	return commentAPI{
		CommentID:    0,
		PostID:       postID,
		User:         userBasicAPI{UserID: userID},
		Comment:      faker.Sentence(10, true),
		LikeCount:    0,
		DislikeCount: 0,
		ReplyCount:   0,
		Createtime:   timestamp,
		Updatetime:   timestamp,
	}
}

func genPlaceFaker(geoPolygon *geo.Polygon, lat0, latR, lon0, lonR float64, r *rand.Rand) placeAPI {
	place := placeAPI{
		Name: "test",
		Lat:  float64(lat0 + latR*r.Float64()),
		Lon:  float64(lon0 + lonR*r.Float64()),
	}
	for {
		if geoPolygon.Contains(geo.NewPoint(place.Lat, place.Lon)) {
			break
		}
		place.Lat = float64(lat0 + latR*r.Float64())
		place.Lon = float64(lon0 + lonR*r.Float64())
	}
	cities, _ := getCityByLocation(place.Lat, place.Lon)
	if len(cities) > 0 {
		place.CityID1 = cities[0].Properties.CityID1
		place.CityID2 = cities[0].Properties.CityID2
		place.CityID3 = cities[0].Properties.CityID3
		place.CityID4 = cities[0].Properties.CityID4
		place.CityID5 = cities[0].Properties.CityID5
		place.CountryCode = cities[0].Properties.CountryCode
	}
	return place
}

func genPopularPostUpsert() {
	users, err := getAllUserID()
	if err != nil {
		log.Println("user", err)
	}
	usersID := []int64{}
	for i := 0; i < len(users); i++ {
		usersID = append(usersID, users[i].UserID)
	}
	for categoryID := range categoryMapID2Name {
		posts, err := getPostsByRecentNum(categoryID, numPopularPostPerRefresh)
		if err != nil {
			log.Println("post", err)
		}
		sort.Sort(postByPopular(posts))
		// init post_user_read_index
		if err := upsertInitPostUserReadIndex(categoryID, usersID); err != nil {
			log.Println(err)
		}
		// post_common
		if err := upsertPopularPostOnPostCommon(categoryID, posts); err != nil {
			log.Println(err)
		}
	}
}
func TestUploadImagesSample(t *testing.T) { // delete this in the future, this is for demo data structure
	type sampleImageCategory struct {
		category map[string]map[string]int
	}
	if os.Getenv("test_fakedata") != "true" || os.Getenv("test_fakedata") == "" {
		t.Skip("skip TestUploadSampleMedia")
	}
	sampleCategoryAppendIndex := func(c sampleImageCategory, cn string, scn []string, isc []int) {
		if len(scn) != len(isc) {
			return
		}
		if c.category[cn] == nil {
			c.category[cn] = make(map[string]int)
		}
		for i := range scn {
			c.category[cn][scn[i]] = isc[i]
		}
	}
	categorySample := sampleImageCategory{category: make(map[string]map[string]int)}

	list := []string{"Fragrance", "Makeup"}
	index := []int{15, 8}
	sampleCategoryAppendIndex(categorySample, "Beauty", list, index)
	list = []string{"Art", "Car", "Design", "Fitness", "Food", "Party", "Pets", "Shopping", "Sports", "Tech", "Travel"}
	index = []int{13, 9, 12, 10, 6, 5, 4, 11, 3, 2, 1}
	sampleCategoryAppendIndex(categorySample, "Life", list, index)
	list = []string{"Formal", "Street"}
	index = []int{14, 7}
	sampleCategoryAppendIndex(categorySample, "Style", list, index)
	list = []string{"Sup"}
	index = []int{0}
	sampleCategoryAppendIndex(categorySample, "Sup", list, index)
	// log.Println(categorySample)

	client, err := storage.NewClient(context.Background(), clientOptionGoogleAPI)
	defer client.Close()
	if err != nil {
		log.Println("gcp sdk client", err)
	}
	for k0, v0 := range categorySample.category {
		for k1, v1 := range v0 {
			for s := 1; s <= 7; s++ {
				// local
				ir, err := ioReaderFromFile("./development/upload/Pics/" + k0 + "/" + k1 + "/" + k1 + "/" + strconv.Itoa(s) + ".jpg")
				if err != nil {
					log.Println("err", err)
					continue
				}
				// cloud storage
				foldername := bucketImagesCloudStorage + "/sample/" + strconv.Itoa(v1) + "/" + strconv.Itoa(s) + "/"
				filename := "0.jpg"
				if err := writeAndMakePublicCloudStorageGCP(client, globalConfig[env].GCPBucketRootCloudStorage, foldername+filename, ir); err != nil {
					log.Println("upload failed: ", err)
					continue
				}
			}
		}
	}
}
func TestFakeDataSample(t *testing.T) {
	// createtime > 1533631289
	if os.Getenv("test_fakedata") != "true" || os.Getenv("test_fakedata") == "" {
		t.Skip("skip TestFakeData")
	}
	log.Println("start")
	scale := 100
	totalNumUser := 1 * scale
	totalNumPost := 10 * scale
	totalNumPlace := 1 * scale
	usersID := []int64{}
	fake, err := faker.New("en")
	if err != nil {
		log.Fatalln("faker", err)
	}
	r := rand.New(rand.NewSource(99))
	// random people
	log.Printf("start random %d people", totalNumUser)
	for i := 0; i < totalNumUser; i++ {
		user := genXuserFaker(fake, r)
		userID, err := userInsert(user)
		if err == nil {
			usersID = append(usersID, userID)
		}
		// log.Println("user", err)
	}
	// random follow
	log.Println("start random follow")
	for i := 0; i < len(usersID); i++ {
		totalNumFollow := r.Int63n(int64(len(usersID)))
		followingUserIDCheck := make(map[int64]bool)
		count := int64(0)
		follwerUserID := usersID[i]
		followingUserID := usersID[i]
		for count < totalNumFollow {
			followingUserID = usersID[r.Intn(len(usersID))]
			if followingUserID == follwerUserID || followingUserIDCheck[followingUserID] {
				continue
			}
			if _, err := follow(followingUserID, follwerUserID); err != nil {
				log.Println("follow", err)
				continue
			}
			followingUserIDCheck[followingUserID] = true
			count++
		}
	}
	// random place
	type countryLatLng struct {
		Lng float64 `json:"lng"`
		Lat float64 `json:"lat"`
	}
	type countryPolygon struct {
		Polygon []countryLatLng `json:"polygon"`
	}

	var countryPolygonTest countryPolygon
	lat0 := 21.5
	lat1 := 25.5
	lon0 := 119.8
	lon1 := 122.1
	latR := lat1 - lat0
	lonR := lon1 - lon0
	if file, err := ioutil.ReadFile("./development/geo/taiwan.json"); err == nil {
		if err1 := json.Unmarshal(file, &countryPolygonTest); err1 != nil {
			log.Fatalln("taiwan.json parse failed")
		}
	} else {
		log.Fatalln("taiwan.json file err")
	}
	var geoPoints []*geo.Point
	if len(countryPolygonTest.Polygon) > 0 {
		for i := 0; i < len(countryPolygonTest.Polygon); i++ {
			geoPoints = append(geoPoints, geo.NewPoint(countryPolygonTest.Polygon[i].Lat, countryPolygonTest.Polygon[i].Lng))
		}
		geoPoints = append(geoPoints, geo.NewPoint(countryPolygonTest.Polygon[0].Lat, countryPolygonTest.Polygon[0].Lng))
	}
	geoPolygon := geo.NewPolygon(geoPoints)
	placesID := []int64{}
	for i := 0; i < totalNumPlace; i++ {
		place := genPlaceFaker(geoPolygon, lat0, latR, lon0, lonR, r)
		placeID, err := placeInsert(place)
		if err == nil {
			placesID = append(placesID, placeID)
		}
	}
	// random post
	log.Printf("start random total %d post", totalNumPost)
	for i := 0; i < totalNumPost; i++ {
		if i%100 == 0 && i != 0 {
			log.Println("finished", i, "posts")
		}
		userID := usersID[r.Intn(len(usersID))]
		placeID := int64(0)
		if r.Intn(2) == 0 {
			placeID = placesID[r.Intn(len(placesID))]
		}
		post := genPostFaker(userID, placeID, fake, r)
		postID, err := postInsert(post)
		if err != nil {
			log.Println("post err", err)
		}
		userIDsNewReaction := append([]int64{}, usersID...) // prevent same person
		totalNumPostReaction := r.Intn(len(userIDsNewReaction))
		for j := 0; j < totalNumPostReaction; j++ {
			index := r.Intn(len(userIDsNewReaction))
			reactionOnPost := genPostReactionFaker(postID, userIDsNewReaction[index], r)
			if _, err = reactionOnPostSet(reactionOnPost); err != nil {
				log.Println("reaction err", err)
			}
			userIDsNewReaction = append(userIDsNewReaction[:index], userIDsNewReaction[index+1:]...)
		}
		totalNumPostComment := r.Intn(len(usersID))
		for j := 0; j < totalNumPostComment; j++ {
			index := r.Intn(len(usersID))
			commentOnPost := genPostCommentFaker(postID, usersID[index], fake, r)
			if _, err = commentOnPostInsert(commentOnPost); err != nil {
				log.Println("comment err", err)
			}
		}
	}
	// gen popular post
	log.Println("start popular post")
	genPopularPostUpsert()
}
func TestPopularPostUpsertSample(t *testing.T) {
	if os.Getenv("test_fakedata") != "true" || os.Getenv("test_fakedata") == "" {
		t.Skip("skip TestSortPostPopular")
	}
	genPopularPostUpsert()
}
func TestCategorySupCityPost(t *testing.T) {
	countries, err := getCountries()
	if err != nil {
		log.Println(err)
	}
	for i := 0; i < len(countries); i++ {
		postsCountry, err := getPostsByRecentWithCountryNum(countries[i].CountryCode, categorySup, numPopularPostPerRefresh)
		if err != nil {
			continue
		}
		sort.Sort(postByPopular(postsCountry))
		if err := upsertPopularPostOnCountry(countries[i].CountryCode, postsCountry); err != nil {
			log.Println(err)
		}
		for j := cityLevelRangeFirst; j <= cityLevelRangeLast; j++ {
			level := strconv.Itoa(j)
			citiesLevel, err := getCitiesLevelByCityIDLike(level, countries[i].CountryCode)
			if err != nil {
				log.Println(err)
			}
			for k := 0; k < len(citiesLevel); k++ {
				posts, err := getPostsByRecentWithCityNum(level, citiesLevel[k].CityID, categorySup, numPopularPostPerRefresh)
				if err != nil {
					log.Println(err)
				}
				sort.Sort(postByPopular(posts))
				if err := upsertPopularPostOnCity(level, citiesLevel[k].CityID, posts); err != nil {
					log.Println(err)
				}
				log.Println(len(posts))
			}
		}
	}
}
func TestRedisPopularPostUserReadIndex(t *testing.T) {
	startTimeTotal := time.Now()
	startTime := time.Now()
	cr := connectRedis()
	defer cr.client.Close()
	log.Printf("connect total took %fs\n", time.Since(startTime).Seconds())
	startTime = time.Now()
	val, err := cr.client.HGet(parseHashPopularPostUserReadIndex(1), "0").Result()
	if err == redis.Nil {
		log.Println("key does not exist")
	} else if err != nil {
		panic(err)
	} else {
		log.Println(val)
	}
	log.Printf("get total took %fs\n", time.Since(startTime).Seconds())
	startTime = time.Now()
	err = cr.client.HSet(parseHashPopularPostUserReadIndex(1), "0", 1).Err()
	if err != nil {
		log.Println(err)
	}
	log.Printf("set total took %fs\n", time.Since(startTime).Seconds())
	log.Printf("total took %fs\n", time.Since(startTimeTotal).Seconds())
}
func TestRedisBenchmark(t *testing.T) {
	startTime := time.Now()
	cr := connectRedis()
	defer cr.client.Close()
	var wg sync.WaitGroup
	numReq := 1
	wg.Add(numReq)
	log.Printf("connect total took %fs\n", time.Since(startTime).Seconds())
	for i := 0; i < numReq; i++ {
		go func() {
			defer wg.Done()
			val, err := cr.client.HGet("user1", "3").Result()
			if err == redis.Nil {
				log.Println("key does not exist")
			} else if err != nil {
				panic(err)
			} else {
				log.Println(val)
			}
		}()
	}
	wg.Wait()
	log.Printf("get total took %fs\n", time.Since(startTime).Seconds())
}
