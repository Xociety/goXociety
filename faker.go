package main

import (
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"

	"github.com/globalsign/mgo/bson"
	geo "github.com/kellydunn/golang-geo"
	"github.com/manveru/faker"
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
		if err := upsertInitCommonPostUserReadIndex(categoryID, usersID); err != nil {
			log.Println(err)
		}
		// post_common
		if err := upsertPopularPostOnPostCommon(categoryID, posts); err != nil {
			log.Println(err)
		}
	}
	log.Println("finish common posts")
}
func genSupPopularPostUpsert() {
	users, err := getAllUserID()
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
	categorySupStr := strconv.Itoa(categorySup)
	countries, err := getCountries()
	if err != nil {
		log.Println(err)
	}
	citiesAll := [][]city2API{}
	for j := cityLevelRangeFirst; j <= cityLevelRangeLast; j++ {
		level := strconv.Itoa(j)
		cities, err := getCities(level)
		if err != nil {
			log.Panicln(err)
		}
		citiesAll = append(citiesAll, cities)
	}
	for i := 0; i < len(countries); i++ {
		posts, err := getPostsByRecentWithPlaceLikeNum(countries[i].CountryCode, categorySup, numPopularPostPerRefresh)
		sort.Sort(postByPopular(posts))
		// init post_user_read_index
		cr := connectRedis()
		hk := redisHashCountryPopularPostUserReadIndex + ":" + countries[i].CountryCode + ":" + categorySupStr
		err = cr.client.HMSet(hk, hf).Err()
		if err != nil {
			log.Println(err)
		}
		// country.sup_popular_posts
		c, err := connectMongoDB()
		if err != nil {
			log.Println("mongo session", err)
		}
		collection := c.session.DB(mongoDBXociety).C(mongoCollectionCity2)
		selector := bson.M{"level": "0", "country_code": countries[i].CountryCode}
		if _, err := collection.Upsert(selector, bson.M{"$set": bson.M{"sup_popular_posts": posts, "post_count": len(posts)}}); err != nil {
			log.Println("upsertPopularPostOnCountry", err)
		}
		for j := cityLevelRangeFirst; j <= cityLevelRangeLast; j++ {
			level := strconv.Itoa(j)
			filteredCity := []city2API{}
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
				filteredPosts := []postAPI{}
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
				hk := redisHashCityPopularPostUserReadIndex + ":" + fCityID + ":" + categorySupStr
				err = cr.client.HMSet(hk, hf).Err()
				if err != nil {
					log.Println(err)
				}
				// city.sup_popular_posts
				collection := c.session.DB(mongoDBXociety).C(mongoCollectionCity2)
				selector := bson.M{"level": level, "city_id": fCityID}
				if _, err := collection.Upsert(selector, bson.M{"$set": bson.M{"sup_popular_posts": filteredPosts, "post_count": len(filteredPosts)}}); err != nil {
					log.Println("upsertPopularPostOnCity", err)
				}
			}
		}
		c.session.Close()
		cr.client.Close()
		log.Println("finish " + countries[i].CountryName + " posts")
	}
	log.Println("finish sup city posts")
}
