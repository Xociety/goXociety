package main

import (
	"log"
	"math/rand"
	"sort"
	"strconv"

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
	for i := 0; i < len(countries); i++ {
		postsCountry, err := getPostsByRecentWithCountryNum(countries[i].CountryCode, categorySup, numPopularPostPerRefresh)
		if err != nil {
			continue
		}
		sort.Sort(postByPopular(postsCountry))
		// init post_user_read_index
		cr := connectRedis()
		hk := redisHashCountryPopularPostUserReadIndex + ":" + countries[i].CountryCode + ":category_id_" + categorySupStr
		err = cr.client.HMSet(hk, hf).Err()
		if err != nil {
			log.Println(err)
		}
		// country.sup_popular_posts
		if err := upsertSupPopularPostOnCountry(countries[i].CountryCode, postsCountry); err != nil {
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
				// init post_user_read_index
				hk := redisHashCityPopularPostUserReadIndex + ":" + countries[i].CountryCode + ":category_id_" + categorySupStr
				err = cr.client.HMSet(hk, hf).Err()
				if err != nil {
					log.Println(err)
				}
				// city.sup_popular_posts
				if err := upsertSupPopularPostOnCity(level, citiesLevel[k].CityID, posts); err != nil {
					log.Println(err)
				}
			}
		}
		cr.client.Close()
		log.Println("finish " + countries[i].CountryName + " posts")
	}
}
