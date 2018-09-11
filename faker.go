package main

import (
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"

	"github.com/chienfuchen32/goXociety/x/common"
	"github.com/chienfuchen32/goXociety/x/config"
	"github.com/chienfuchen32/goXociety/x/io"
	"github.com/globalsign/mgo/bson"
	geo "github.com/kellydunn/golang-geo"
	"github.com/manveru/faker"
)

func genXuserFaker(fake *faker.Faker, r *rand.Rand) config.UserDB {
	timeNow := common.GetNowUnixTimestamp()
	return config.UserDB{
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

func genPostFaker(userID, placeID int64, faker *faker.Faker, r *rand.Rand) config.PostAPI {
	timeNow := common.GetNowUnixTimestamp()
	postType := 0             // r.Intn(2)
	originWidth := 0          // plz check current sample size
	originHeight := 0         // plz check current sample size
	categoryID := r.Intn(16)  // plz check current category
	sampleID := r.Intn(7) + 1 // plz check current sample file
	if categoryID == 0 {
		sampleID = 1
	}
	switch config.PostTypeMapID2Type[postType] {
	case config.MediaFormatJPG:
		originWidth = 1920
		originHeight = 1280
		if categoryID != 0 {
			originWidth = 1242
			originHeight = 2004
		}
	case config.MediaFormatHLS:
		originWidth = 1920
		originHeight = 1080
	}
	post := config.PostAPI{
		User: config.UserBasicAPI{
			UserID: userID,
		},
		Content: "test post",
		Blob: config.BlobAPI{
			BlobID:       "sample/" + strconv.Itoa(categoryID) + "/" + strconv.Itoa(sampleID),
			OriginWidth:  originWidth,
			OriginHeight: originHeight,
		},
		Type:         postType,
		LikeCount:    0,
		DislikeCount: 0,
		CommentCount: 0,
		Place:        config.PlaceAPI{PlaceID: placeID},
		CategoryID:   categoryID,
		Public:       true,
		Createtime:   timeNow,
		Updatetime:   timeNow,
	}
	return post
}

func genPostReactionFaker(postID, userID int64, r *rand.Rand) config.ReactionOnPostAPI {
	return config.ReactionOnPostAPI{
		PostID:     postID,
		User:       config.UserBasicAPI{UserID: userID},
		ReactionID: r.Intn(2),
		Createtime: common.GetNowUnixTimestamp(),
	}
}

func genPostCommentFaker(postID, userID int64, faker *faker.Faker, r *rand.Rand) config.CommentAPI {
	timestamp := common.GetNowUnixTimestamp()
	return config.CommentAPI{
		CommentID:    0,
		PostID:       postID,
		User:         config.UserBasicAPI{UserID: userID},
		Comment:      faker.Sentence(10, true),
		LikeCount:    0,
		DislikeCount: 0,
		ReplyCount:   0,
		Createtime:   timestamp,
		Updatetime:   timestamp,
	}
}

func genPlaceFaker(cm *config.ConnMongo, geoPolygon *geo.Polygon, lat0, latR, lon0, lonR float64, r *rand.Rand) config.PlaceAPI {
	place := config.PlaceAPI{
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
	cities, _ := io.GetCityByLocation(cm, place.Lat, place.Lon)
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

func genPopularPostUpsert(c *config.ConnPostgres, cm *config.ConnMongo) {
	users, err := io.GetAllUserID(c)
	if err != nil {
		log.Println("user", err)
	}
	usersID := []int64{}
	for i := 0; i < len(users); i++ {
		usersID = append(usersID, users[i].UserID)
	}
	for categoryID := range config.CategoryMapID2Name {
		posts, err := io.GetPostsByRecentNum(c, categoryID, config.NumPopularPostPerRefresh)
		if err != nil {
			log.Println("post", err)
		}
		sort.Sort(common.PostByPopular(posts))
		// post_common
		if err := io.UpsertPopularPostOnPostCommon(cm, categoryID, posts); err != nil {
			log.Println(err)
		}
	}
	log.Println("finish common posts")
}
func genSupPopularPostUpsert(c *config.ConnPostgres) {
	users, err := io.GetAllUserID(c)
	if err != nil {
		log.Println("user", err)
	}
	usersID := []int64{}
	usersIDStr := []string{}
	hf := make(map[string]interface{})
	for i := 0; i < len(users); i++ {
		usersID = append(usersID, users[i].UserID)
		userIDStr := strconv.FormatInt(users[i].UserID, 10)
		usersIDStr = append(usersIDStr, userIDStr)
		hf[userIDStr] = "0"
	}
	// categorySupStr := strconv.Itoa(config.CategorySup)
	cm, err := io.ConnectMongoDB()
	if err != nil {
		log.Println("mongo session", err)
	}
	countries, err := io.GetCountries(&cm)
	if err != nil {
		log.Println(err)
	}
	citiesAll := [][]config.CityAPI{}
	for j := config.CityLevelRangeFirst; j <= config.CityLevelRangeLast; j++ {
		level := strconv.Itoa(j)
		cities, err := io.GetCities(&cm, level)
		if err != nil {
			log.Panicln(err)
		}
		citiesAll = append(citiesAll, cities)
	}
	cm.Session.Close()
	log.Println("city popular post")
	for i := 0; i < len(countries); i++ {
		posts, err := io.GetPostsByRecentWithPlaceLikeNum(c, countries[i].CountryCode, config.CategorySup, config.NumPopularPostPerRefresh)
		sort.Sort(common.PostByPopular(posts))
		// init post_user_read_index
		// cr := connectRedis()
		// hk := config.RedisHashCountryPopularPostUserReadIndex + ":" + countries[i].CountryCode + ":" + categorySupStr
		// err = cr.client.HMSet(hk, hf).Err()
		// if err != nil {
		// 	log.Println(err)
		// }
		// country.sup_popular_posts
		c, err := io.ConnectMongoDB()
		if err != nil {
			log.Println("mongo session", err)
		}
		collection := c.Session.DB(config.MongoDBXociety).C(config.MongoCollectionCity)
		selector := bson.M{"level": "0", "country_code": countries[i].CountryCode}
		if _, err := collection.Upsert(selector, bson.M{"$set": bson.M{"sup_popular_posts": posts, "post_count": len(posts)}}); err != nil {
			log.Println("upsertPopularPostOnCountry", err)
		}
		if len(posts) > 0 {

			for j := config.CityLevelRangeFirst; j <= config.CityLevelRangeLast; j++ {
				level := strconv.Itoa(j)
				filteredCity := []config.CityAPI{}
				for k := 0; k < len(citiesAll[j-1]); k++ {
					cityID := ""
					switch j {
					case 1:
						cityID = citiesAll[j-1][k].CityID1
					case 2:
						cityID = citiesAll[j-1][k].CityID2
					case 3:
						cityID = citiesAll[j-1][k].CityID3
					case 4:
						cityID = citiesAll[j-1][k].CityID4
					case 5:
						cityID = citiesAll[j-1][k].CityID5
					}
					if strings.Index(cityID, countries[i].CountryCode) != -1 {
						filteredCity = append(filteredCity, citiesAll[j-1][k])
					}
				}
				for k := 0; k < len(filteredCity); k++ {
					fCityID := ""
					switch j {
					case 1:
						fCityID = filteredCity[k].CityID1
					case 2:
						fCityID = filteredCity[k].CityID2
					case 3:
						fCityID = filteredCity[k].CityID3
					case 4:
						fCityID = filteredCity[k].CityID4
					case 5:
						fCityID = filteredCity[k].CityID5
					}
					filteredPosts := []config.PostAPI{}
					for l := 0; l < len(posts); l++ {
						cityID := ""
						switch j {
						case 1:
							cityID = posts[l].Place.CityID1
						case 2:
							cityID = posts[l].Place.CityID2
						case 3:
							cityID = posts[l].Place.CityID3
						case 4:
							cityID = posts[l].Place.CityID4
						case 5:
							cityID = posts[l].Place.CityID5
						}
						if strings.Index(cityID, fCityID) != -1 {
							filteredPosts = append(filteredPosts, posts[l])
						}
					}
					// init post_user_read_index
					// hk := config.RedisHashCityPopularPostUserReadIndex + ":" + fCityID + ":" + categorySupStr
					// err = cr.client.HMSet(hk, hf).Err()
					// if err != nil {
					// 	log.Println(err)
					// }
					// city.sup_popular_posts
					collection := c.Session.DB(config.MongoDBXociety).C(config.MongoCollectionCity)
					selector := bson.M{"level": level, "city_id_" + level: fCityID}
					if _, err := collection.Upsert(selector, bson.M{"$set": bson.M{"sup_popular_posts": filteredPosts, "post_count": len(filteredPosts)}}); err != nil {
						log.Println("upsertPopularPostOnCity", err)
					}
				}
			}
		} else {
			collection := c.Session.DB(config.MongoDBXociety).C(config.MongoCollectionCity)
			selector := bson.M{"country_code": countries[i].CountryCode}
			if _, err := collection.Upsert(selector, bson.M{"$set": bson.M{"sup_popular_posts": 0, "post_count": 0}}); err != nil {
				log.Println("upsertPopularPostOnCity", err)
			}
		}
		c.Session.Close()
		// cr.client.Close()
		log.Println("finish " + countries[i].CountryName + "," + strconv.Itoa(len(posts)) + " posts")
	}
	log.Println("init post user read index, users", len(users))
	// mongo version of post user read index
	for u := 0; u < len(usersID); u++ {
		rID := config.PopularPostUserReadIndexAPI{
			UserID:                     usersID[u],
			CommonPopularPostIndex:     make(map[int]int),
			CitySupPopularPostIndex:    make(map[string]int),
			CountrySupPopularPostIndex: make(map[string]int),
		}
		c, err := io.ConnectMongoDB()
		if err != nil {
			log.Println("mongo session", err)
		}
		collection := c.Session.DB(config.MongoDBXociety).C(config.MongoCollectionPopularPostUserReadIndex)
		selector := bson.M{"user_id": rID.UserID}
		if _, err := collection.Upsert(selector, bson.M{"user_id": rID.UserID}); err != nil {
			log.Println("upser popular_post_user_read_index", err)
		}
		c.Session.Close()
	}

	rID := config.PopularPostUserReadIndexAPI{
		UserID:                     0,
		CommonPopularPostIndex:     make(map[int]int),
		CitySupPopularPostIndex:    make(map[string]int),
		CountrySupPopularPostIndex: make(map[string]int),
	}
	cm, err = io.ConnectMongoDB()
	if err != nil {
		log.Println("mongo session", err)
	}
	for categoryID := range config.CategoryMapID2Name {
		rID.CommonPopularPostIndex[categoryID] = 0
	}
	for i := 0; i < len(countries); i++ {
		rID.CountrySupPopularPostIndex[countries[i].CountryCode] = 0
	}
	for j := config.CityLevelRangeFirst; j <= config.CityLevelRangeLast; j++ {
		for k := 0; k < len(citiesAll[j-1]); k++ {
			CityID := ""
			switch j {
			case 1:
				CityID = citiesAll[j-1][k].CityID1
			case 2:
				CityID = citiesAll[j-1][k].CityID2
			case 3:
				CityID = citiesAll[j-1][k].CityID3
			case 4:
				CityID = citiesAll[j-1][k].CityID4
			case 5:
				CityID = citiesAll[j-1][k].CityID5
			}
			rID.CitySupPopularPostIndex[CityID] = 0
		}
	}
	collection := cm.Session.DB(config.MongoDBXociety).C(config.MongoCollectionPopularPostUserReadIndex)
	if _, err := collection.UpdateAll(bson.M{}, bson.M{"$set": bson.M{
		"common_popular_post_index":      rID.CommonPopularPostIndex,
		"city_sup_popular_post_index":    rID.CitySupPopularPostIndex,
		"country_sup_popular_post_index": rID.CountrySupPopularPostIndex,
	}}); err != nil {
		log.Println("upser popular_post_user_read_index", err)
	}
	cm.Session.Close()
	log.Println("finish sup city posts")
}
