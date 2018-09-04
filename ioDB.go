package main

import (
	"crypto/tls"
	"database/sql"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"strconv"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
)

// net
type connPostgres struct {
	db *sql.DB
}

type connMongo struct {
	session *mgo.Session
}
type connRedis struct {
	client *redis.Client
}

func connectPostgres() (connPostgres, error) {
	var err error
	c := connPostgres{}
	c.db, err = sql.Open("postgres", postgresConStr)
	if err != nil {
		log.Println("postgres connect", err)
		return c, err
	}
	return c, nil
}
func connectMongoDB() (connMongo, error) {
	var err error
	c := connMongo{}
	switch env {
	case "development":
		c.session, err = mgo.Dial(globalConfig[env].MongoConStr)
	case "staging":
		// https://godoc.org/github.com/globalsign/mgo#example-Dial--TlsConfig
		clientCertPEM, err := ioutil.ReadFile(globalConfig[env].MongoSecretFolderPath + globalConfig[env].MongoSecretCertFilename)
		if err != nil {
			return c, errors.New("mongo pem not found")
		}
		clientKeyPEM, err := ioutil.ReadFile(globalConfig[env].MongoSecretFolderPath + globalConfig[env].MongoSecretKeyFilename)
		if err != nil {
			return c, errors.New("mongo pem not found")
		}
		clientCert, err := tls.X509KeyPair(clientCertPEM, clientKeyPEM)
		if err != nil {
			return c, errors.New("mongo pem not found")
		}
		// clientCert.Leaf, err = x509.ParseCertificate(clientCert.Certificate[0])
		tlsConfig := &tls.Config{
			Certificates:       []tls.Certificate{clientCert},
			InsecureSkipVerify: true,
		}
		c.session, err = mgo.DialWithInfo(&mgo.DialInfo{
			Addrs: []string{globalConfig[env].MongoConStr},
			DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
				return tls.Dial("tcp", globalConfig[env].MongoConStr, tlsConfig)
			},
			Database: globalSecret.Mongo.MongoDatabase,
			Username: globalSecret.Mongo.MongoUsername,
			Password: globalSecret.Mongo.MongoPassword,
		})
	case "production":
		// https://godoc.org/github.com/globalsign/mgo#example-Dial--TlsConfig
		clientCertPEM, err := ioutil.ReadFile(globalConfig[env].MongoSecretFolderPath + globalConfig[env].MongoSecretCertFilename)
		if err != nil {
			return c, errors.New("mongo pem not found")
		}
		clientKeyPEM, err := ioutil.ReadFile(globalConfig[env].MongoSecretFolderPath + globalConfig[env].MongoSecretKeyFilename)
		if err != nil {
			return c, errors.New("mongo pem not found")
		}
		clientCert, err := tls.X509KeyPair(clientCertPEM, clientKeyPEM)
		if err != nil {
			return c, errors.New("mongo pem not found")
		}
		// clientCert.Leaf, err = x509.ParseCertificate(clientCert.Certificate[0])
		tlsConfig := &tls.Config{
			Certificates:       []tls.Certificate{clientCert},
			InsecureSkipVerify: true,
		}
		c.session, err = mgo.DialWithInfo(&mgo.DialInfo{
			Addrs: []string{globalConfig[env].MongoConStr},
			DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
				return tls.Dial("tcp", globalConfig[env].MongoConStr, tlsConfig)
			},
			Database: globalSecret.Mongo.MongoDatabase,
			Username: globalSecret.Mongo.MongoUsername,
			Password: globalSecret.Mongo.MongoPassword,
		})
	}
	if err != nil {
		log.Println("mongo connect", err)
		return c, err
	}
	return c, nil
}
func connectRedis() connRedis {
	c := connRedis{}
	c.client = redis.NewClient(&redis.Options{
		Addr:     globalConfig[env].RedisConStr,
		Password: "",
		DB:       redisDBPopularPostUserReadIndex,
	})
	return c
}

// auth
func checkSession(userToken string) (user xuserAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return user, errors.New("db connection")
	}
	defer c.db.Close()
	row := c.db.QueryRow(`
		SELECT 
		user_id, username, email, name, phone, 
		gender, bio, credit, photo_url, 
		language_id, country_code, 
		timezone, last_ip, createtime, updatetime 
		FROM xuser WHERE user_id=$1;`,
		userToken)
	if err := row.Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.Name,
		&user.Phone,
		&user.Gender,
		&user.Bio,
		&user.Credit,
		&user.PhotoURL,
		&user.LanguageID,
		&user.CountryCode,
		&user.Timezone,
		&user.LastIP,
		&user.Createtime,
		&user.Updatetime,
	); err != nil {
		return user, err
	}

	return user, nil
}

// [query]
func login(email, password string) (lc loginAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return lc, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `SELECT user_id FROM xuser WHERE email=$1 AND password=$2;`
	row := c.db.QueryRow(sqlStr, email, password)
	if err := row.Scan(
		&lc.Token,
	); err != nil {
		log.Println("login", err)
		return lc, errors.New("not valid")
	}
	return lc, nil
}

// common
func getCategories() (categories []categoryAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return categories, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		SELECT 
		category_id, category_name
		FROM category;
	`
	rows, err := c.db.Query(sqlStr)
	if err != nil {
		log.Println("getCategories", err)
		return categories, err
	}
	defer rows.Close()
	for rows.Next() {
		category := categoryAPI{}
		if err := rows.Scan(
			&category.CategoryID,
			&category.CategoryName,
		); err != nil {
			log.Println(err)
			return categories, err
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return categories, err
	}
	return categories, nil
}
func getGender() (genderList []genderAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return genderList, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		SELECT 
		gender_id, value
		FROM gender;
	`
	rows, err := c.db.Query(sqlStr)
	if err != nil {
		log.Println("getGender", err)
		return genderList, err
	}
	defer rows.Close()
	for rows.Next() {
		gender := genderAPI{}
		if err := rows.Scan(
			&gender.GenderID,
			&gender.Value,
		); err != nil {
			log.Println(err)
			return genderList, err
		}
		genderList = append(genderList, gender)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return genderList, err
	}
	return genderList, nil
}
func getLanguages() (languages []languageAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return languages, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		SELECT 
		language_id, display_language, value
		FROM language;
	`
	rows, err := c.db.Query(sqlStr)
	if err != nil {
		log.Println("getLanguages", err)
		return languages, err
	}
	defer rows.Close()
	for rows.Next() {
		language := languageAPI{}
		if err := rows.Scan(
			&language.LanguageID,
			&language.DisplayLanguage,
			&language.Value,
		); err != nil {
			log.Println(err)
			return languages, err
		}
		languages = append(languages, language)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return languages, err
	}
	return languages, nil
}
func getPostType() (postTypeList []postTypeAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return postTypeList, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		SELECT 
		post_type_id, value
		FROM post_type;
	`
	rows, err := c.db.Query(sqlStr)
	if err != nil {
		log.Println("getPostType", err)
		return postTypeList, err
	}
	defer rows.Close()
	for rows.Next() {
		postType := postTypeAPI{}
		if err := rows.Scan(
			&postType.PostTypeID,
			&postType.Value,
		); err != nil {
			log.Println(err)
			return postTypeList, err
		}
		switch postType.Value {
		case mediaFormatJPG:
			postType.FileFormat = []string{mediaFormatJPG}
		case mediaFormatHLS:
			postType.FileFormat = []string{mediaFormatM3U8, mediaFormatTS}
		}
		postTypeList = append(postTypeList, postType)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return postTypeList, err
	}
	return postTypeList, nil
}
func getReaction() (reactionList []reactionAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return reactionList, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		SELECT 
		reaction_id, value
		FROM reaction;
	`
	rows, err := c.db.Query(sqlStr)
	if err != nil {
		log.Println("getReaction", err)
		return reactionList, err
	}
	defer rows.Close()
	for rows.Next() {
		reaction := reactionAPI{}
		if err := rows.Scan(
			&reaction.ReactionID,
			&reaction.Value,
		); err != nil {
			log.Println(err)
			return reactionList, err
		}
		reactionList = append(reactionList, reaction)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return reactionList, err
	}
	return reactionList, nil
}

// user
func getUserByUserID(userID int64) (user xuserAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return user, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		SELECT 
		user_id, username, email, name, phone, 
		gender, bio, credit, photo_url, 
		language_id, country_code, 
		timezone, last_ip, createtime, updatetime 
		FROM xuser WHERE user_id=$1;
	`
	row := c.db.QueryRow(sqlStr, userID)
	if err := row.Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.Name,
		&user.Phone,
		&user.Gender,
		&user.Bio,
		&user.Credit,
		&user.PhotoURL,
		&user.LanguageID,
		&user.CountryCode,
		&user.Timezone,
		&user.LastIP,
		&user.Createtime,
		&user.Updatetime,
	); err != nil {
		log.Println("getXuserByID", err)
		return user, errors.New("user not found")
	}
	return user, nil
}
func getUserByUsername(username string) (user xuserAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return user, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		SELECT 
		user_id, username, email, name, phone, 
		gender, bio, credit, photo_url, 
		language_id, country_code, 
		timezone, last_ip, createtime, updatetime 
		FROM xuser WHERE username=$1;
	`
	row := c.db.QueryRow(sqlStr, username)
	if err := row.Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.Name,
		&user.Phone,
		&user.Gender,
		&user.Bio,
		&user.Credit,
		&user.PhotoURL,
		&user.LanguageID,
		&user.CountryCode,
		&user.Timezone,
		&user.LastIP,
		&user.Createtime,
		&user.Updatetime,
	); err != nil {
		// log.Println("getXuserByUsername", err)
		return user, errors.New("user not found")
	}
	return user, nil
}

// follow
func getUsersByFollowing(followerUserID int64, page int) (users []userFollowingAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return users, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		SELECT xuser.user_id, xuser.username, xuser.name, xuser.photo_url, follow.createtime
		FROM follow 
		JOIN xuser ON follow.following_user_id = xuser.user_id
		WHERE follow.follower_user_id=$1 AND follow.valid=true 
		ORDER BY follow.createtime DESC OFFSET $2 LIMIT $3;
	`
	rows, err := c.db.Query(sqlStr, followerUserID, page*numPerRequest, numPerRequest)
	if err != nil {
		log.Println("getFollowingList1", err)
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		user := userFollowingAPI{}
		if err := rows.Scan(
			&user.UserID,
			&user.UserName,
			&user.Name,
			&user.PhotoURL,
			&user.FollowingTime,
		); err != nil {
			log.Println("err", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		log.Println("getFollowingList2", err)
		return users, err
	}
	return users, nil
}
func getUsersByFollower(followerUserID int64, page int) (users []userFollowerAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return users, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		SELECT xuser.user_id, xuser.username, xuser.name, xuser.photo_url, follow.createtime
		FROM follow 
		JOIN xuser ON follow.follower_user_id = xuser.user_id
		WHERE follow.following_user_id=$1 AND follow.valid=true 
		ORDER BY follow.createtime DESC OFFSET $2 LIMIT $3;
	`
	rows, err := c.db.Query(sqlStr, followerUserID, page*numPerRequest, numPerRequest)
	if err != nil {
		log.Println("getFollwerList1", err)
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		user := userFollowerAPI{}
		if err := rows.Scan(
			&user.UserID,
			&user.UserName,
			&user.Name,
			&user.PhotoURL,
			&user.FollowingTime,
		); err != nil {
			log.Println("err", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		log.Println("getFollwerList2", err)
		return users, err
	}
	return users, nil
}
func checkUserIfFollowing(followingUserID, followerUserID int64) (isFollowing bool, err error) {
	c, err := connectPostgres()
	if err != nil {
		return isFollowing, errors.New("db connection")
	}
	defer c.db.Close()
	count := 0
	sqlStr := `SELECT COUNT(*) FROM follow WHERE following_user_id=$1 AND follower_user_id=$2;`
	err = c.db.QueryRow(sqlStr, followingUserID, followerUserID).Scan(&count)
	if err != nil {
		return count == 1, err
	}
	return count == 1, err
}

// country, city
func getCountries() (countries []cityAPI, err error) {
	c, err := connectMongoDB()
	if err != nil {
		log.Println("mongo session", err)
		return countries, err
	}
	defer c.session.Close()
	collection := c.session.DB(mongoDBXociety).C(mongoCollectionCity)
	q := bson.M{"level": "0"}
	selector := bson.M{"country_code": 1, "country_name": 1}
	if err := collection.Find(q).Select(selector).All(&countries); err != nil {
		log.Println("getCountries", err)
		return countries, err
	}
	return countries, nil
}
func getCountry(countryCode string) (country cityAPI, err error) {
	c, err := connectMongoDB()
	if err != nil {
		log.Println("mongo session", err)
		return country, err
	}
	defer c.session.Close()
	collection := c.session.DB(mongoDBXociety).C(mongoCollectionCity)
	q := bson.M{"level": "0", "country_code": countryCode}
	selector := bson.M{"country_code": 1, "country_name": 1}
	if err := collection.Find(q).Select(selector).One(&country); err != nil {
		log.Println("getCountry", err)
		return country, err
	}
	return country, nil
}
func getCities(level string) (cities []cityAPI, err error) {
	c, err := connectMongoDB()
	if err != nil {
		log.Println("mongo session", err)
		return cities, err
	}
	defer c.session.Close()
	collection := c.session.DB(mongoDBXociety).C(mongoCollectionCity)
	q := bson.M{"level": level}
	if err := collection.Find(q).All(&cities); err != nil {
		log.Println("getCities", err)
		return cities, err
	}
	return cities, nil
}
func getCity(level string, cityID string) (city cityAPI, err error) {
	c, err := connectMongoDB()
	if err != nil {
		log.Println("mongo session", err)
		return city, err
	}
	defer c.session.Close()
	collection := c.session.DB(mongoDBXociety).C(mongoCollectionCity)
	q := bson.M{"level": level, "city_id_" + level: cityID}
	if err := collection.Find(q).One(&city); err != nil {
		log.Println("getCity", err)
		return city, err
	}
	return city, nil
}
func getCitiesLevelByCityIDLike(level string, cityID string) (cities []cityAPI, err error) {
	c, err := connectMongoDB()
	if err != nil {
		log.Println("mongo session", err)
		return cities, err
	}
	defer c.session.Close()
	collection := c.session.DB(mongoDBXociety).C(mongoCollectionCity)
	q := bson.M{
		"level":            level,
		"city_id_" + level: bson.M{"$regex": bson.RegEx{Pattern: cityID, Options: ""}}}
	if err := collection.Find(q).All(&cities); err != nil {
		log.Println("getCitiesLevelByCityIDLike", err)
		return cities, err
	}
	return cities, nil
}
func getCityByLocation(lat, lon float64) (cities []cityGeometryAPI, err error) {
	c, err := connectMongoDB()
	if err != nil {
		log.Println("mongo session", err)
		return cities, err
	}
	defer c.session.Close()
	numPerRequest := 5
	collection := c.session.DB(mongoDBXociety).C(mongoCollectionCityGeometry)
	s := bson.M{"properties": 1}
	// geoIntersects
	q := bson.M{
		"geometry": bson.M{
			"$geoIntersects": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{lon, lat},
				},
			},
		},
	}
	if err := collection.Find(q).Select(s).Limit(numPerRequest).All(&cities); err != nil {
		log.Println("getCityByLocation", err)
		return cities, err
	}
	if len(cities) == 0 {
		// near
		q = bson.M{
			"geometry": bson.M{
				"$near": bson.M{
					"$geometry": bson.M{
						"type":        "Point",
						"coordinates": []float64{lon, lat},
					},
					"$maxDistance": mongoGeoNearSearchInKM,
				},
			},
		}
		if err := collection.Find(q).Select(s).Limit(numPerRequest).All(&cities); err != nil {
			log.Println("getCityByLocation", err)
			return cities, err
		}
	}
	return cities, nil
}

// place
func getPlacesByPlacesGCP(placesGCP []placeAPI) (places []placeAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return places, errors.New("db connection")
	}
	defer c.db.Close()
	name := []string{}
	lat := []float64{}
	lon := []float64{}
	for i := 0; i < len(placesGCP); i++ {
		name = append(name, placesGCP[i].Name)
		lat = append(lat, placesGCP[i].Lat)
		lon = append(lon, placesGCP[i].Lon)
	}
	sqlStr, args := parsePlaceSelectAllSQL(name, lat, lon)
	rows, err := c.db.Query(sqlStr, args...)
	if err != nil {
		log.Println("getPlacesByPlacesGCP", err)
		return places, err
	}
	defer rows.Close()
	for rows.Next() {
		place := placeAPI{}
		if err := rows.Scan(
			&place.PlaceID,
			&place.CountryCode,
			&place.CityID1,
			&place.CityID2,
			&place.CityID3,
			&place.CityID4,
			&place.CityID5,
			&place.Lat,
			&place.Lon,
			&place.Name,
		); err != nil {
			log.Println(err)
			return places, err
		}
		places = append(places, place)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return places, err
	}
	return places, nil
}

// post
func getPostsByRecent(categoryID, page int) (posts []postAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return posts, errors.New("db connection")
	}
	defer c.db.Close()
	timestamp := getNowUnixTimestamp() - twoMonthsInSecond // sixHoursInSecond
	sqlStr := `
		SELECT 
		post.post_id,
		post.user_id, xuser.username, xuser.name, xuser.photo_url,
		post.content, post.blob_id, post.origin_width, post.origin_height, post.type, 
		post.like_count, post.dislike_count, post.comment_count,
		place.place_id, place.country_code,
		place.city_id_1, place.city_id_2, place.city_id_3,
		place.city_id_4, place.city_id_5,
		place.lat, place.lon, place.name,
		post.category_id, post.createtime, post.updatetime 
		FROM post 
		JOIN xuser ON xuser.user_id = post.user_id
		JOIN place ON place.place_id = post.place_id
		WHERE post.category_id=$1 AND post.createtime>=$2 
		ORDER BY post.createtime DESC OFFSET $3 LIMIT $4;
	`
	rows, err := c.db.Query(sqlStr, categoryID, timestamp, page*numPerRequest, numPerRequest)
	if err != nil {
		log.Println("getPostsRecent", err)
		return posts, err
	}
	defer rows.Close()
	for rows.Next() {
		post := postAPI{}
		place := placeAPI{}
		var placeID interface{}
		if err := rows.Scan(
			&post.PostID,
			&post.User.UserID,
			&post.User.Username,
			&post.User.Name,
			&post.User.PhotoURL,
			&post.Content,
			&post.Blob.BlobID,
			&post.Blob.OriginWidth,
			&post.Blob.OriginHeight,
			&post.Type,
			&post.LikeCount,
			&post.DislikeCount,
			&post.CommentCount,
			&placeID,
			&place.CountryCode,
			&place.CityID1,
			&place.CityID2,
			&place.CityID3,
			&place.CityID4,
			&place.CityID5,
			&place.Lat,
			&place.Lon,
			&place.Name,
			&post.CategoryID,
			&post.Createtime,
			&post.Updatetime,
		); err != nil {
			log.Println("errrr", err)
			return posts, err
		}
		placeIDCheck, isOK := placeID.(int64)
		if isOK {
			place.PlaceID = placeIDCheck
		}
		post.Place = place
		post.Blob.BlobID = makeBlobURL(post)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return posts, err
	}
	return posts, err
}
func getPostsByRecentNum(categoryID, numPost int) (posts []postAPI, err error) {
	// combine this with func getPostsByRecent
	c, err := connectPostgres()
	if err != nil {
		return posts, errors.New("db connection")
	}
	defer c.db.Close()
	timestamp := getNowUnixTimestamp() - twoMonthsInSecond
	sqlStr := `
		SELECT 
		post.post_id,
		post.user_id, xuser.username, xuser.name, xuser.photo_url,
		post.content, post.blob_id, post.origin_width, post.origin_height, post.type, 
		post.like_count, post.dislike_count, post.comment_count,
		place.place_id, place.country_code,
		place.city_id_1, place.city_id_2, place.city_id_3,
		place.city_id_4, place.city_id_5,
		place.lat, place.lon, place.name,
		post.category_id, post.createtime, post.updatetime 
		FROM post 
		JOIN xuser ON xuser.user_id = post.user_id
		JOIN place ON place.place_id = post.place_id
		WHERE post.category_id=$1 AND post.createtime>=$2 
		ORDER BY post.createtime DESC OFFSET $3 LIMIT $4;
	`
	rows, err := c.db.Query(sqlStr, categoryID, timestamp, 0, numPost)
	if err != nil {
		log.Println("getPostsByRecentNum", err)
		return posts, err
	}
	defer rows.Close()
	for rows.Next() {
		post := postAPI{}
		place := placeAPI{}
		var placeID interface{}
		if err := rows.Scan(
			&post.PostID,
			&post.User.UserID,
			&post.User.Username,
			&post.User.Name,
			&post.User.PhotoURL,
			&post.Content,
			&post.Blob.BlobID,
			&post.Blob.OriginWidth,
			&post.Blob.OriginHeight,
			&post.Type,
			&post.LikeCount,
			&post.DislikeCount,
			&post.CommentCount,
			&placeID,
			&place.CountryCode,
			&place.CityID1,
			&place.CityID2,
			&place.CityID3,
			&place.CityID4,
			&place.CityID5,
			&place.Lat,
			&place.Lon,
			&place.Name,
			&post.CategoryID,
			&post.Createtime,
			&post.Updatetime,
		); err != nil {
			log.Println("errrr", err)
			return posts, err
		}
		placeIDCheck, isOK := placeID.(int64)
		if isOK {
			place.PlaceID = placeIDCheck
		}
		post.Place = place
		post.Blob.BlobID = makeBlobURL(post)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return posts, err
	}
	return posts, nil
}
func getPostsByRecentWithCountry(countryCode string, categoryID, page int) (posts []postAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return posts, errors.New("db connection")
	}
	defer c.db.Close()
	timestamp := getNowUnixTimestamp() - twoMonthsInSecond
	if categoryID == categorySup {
		timestamp = getNowUnixTimestamp() - twentyFourHoursInSecond
	}
	sqlStr := `
		SELECT 
		post.post_id, 
		post.user_id, xuser.username, xuser.name, xuser.photo_url,
		post.content, post.blob_id, post.origin_width, post.origin_height, post.type,
		post.like_count, post.dislike_count, post.comment_count,
		place.place_id, place.country_code,
		place.city_id_1, place.city_id_2, place.city_id_3,
		place.city_id_4, place.city_id_5,
		place.lat, place.lon, place.name,
		post.category_id, post.createtime, post.updatetime 
		FROM post
		JOIN xuser ON xuser.user_id = post.user_id
		JOIN place ON place.place_id = post.place_id
		WHERE post.category_id=$1 AND place.country_code = $2 AND post.createtime >= $3
		ORDER BY post.createtime
		DESC OFFSET $4 LIMIT $5;
	`
	rows, err := c.db.Query(sqlStr, categoryID, countryCode, timestamp, page*numPerRequest, numPerRequest)
	if err != nil {
		log.Println("getPostsByRecentWithCountryCode", err)
		return posts, err
	}
	defer rows.Close()
	for rows.Next() {
		post := postAPI{}
		place := placeAPI{}
		var placeID interface{}
		if err := rows.Scan(
			&post.PostID,
			&post.User.UserID,
			&post.User.Username,
			&post.User.Name,
			&post.User.PhotoURL,
			&post.Content,
			&post.Blob.BlobID,
			&post.Blob.OriginWidth,
			&post.Blob.OriginHeight,
			&post.Type,
			&post.LikeCount,
			&post.DislikeCount,
			&post.CommentCount,
			&placeID,
			&place.CountryCode,
			&place.CityID1,
			&place.CityID2,
			&place.CityID3,
			&place.CityID4,
			&place.CityID5,
			&place.Lat,
			&place.Lon,
			&place.Name,
			&post.CategoryID,
			&post.Createtime,
			&post.Updatetime,
		); err != nil {
			log.Println(err)
			return posts, err
		}
		placeIDCheck, isOK := placeID.(int64)
		if isOK {
			place.PlaceID = placeIDCheck
		}
		post.Place = place
		post.Blob.BlobID = makeBlobURL(post)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return posts, err
	}
	// log.Println("posts", posts)
	return posts, nil
}
func getPostsByRecentWithCity(level, cityID string, categoryID, page int) (posts []postAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return posts, errors.New("db connection")
	}
	defer c.db.Close()
	timestamp := getNowUnixTimestamp() - twoMonthsInSecond
	if categoryID == categorySup {
		timestamp = getNowUnixTimestamp() - twentyFourHoursInSecond
	}
	sqlStr := `
		SELECT 
		post.post_id, 
		post.user_id, xuser.username, xuser.name, xuser.photo_url,
		post.content, post.blob_id, post.origin_width, post.origin_height, post.type,
		post.like_count, post.dislike_count, post.comment_count,
		place.place_id, place.country_code,
		place.city_id_1, place.city_id_2, place.city_id_3,
		place.city_id_4, place.city_id_5,
		place.lat, place.lon, place.name,
		post.category_id, post.createtime, post.updatetime 
		FROM post
		JOIN xuser ON xuser.user_id = post.user_id
		JOIN place ON place.place_id = post.place_id
		WHERE post.category_id=$1 AND place.city_id_` + level + ` = $2 AND post.createtime >= $3
		ORDER BY post.createtime
		DESC OFFSET $4 LIMIT $5;
	`
	// log.Println(sqlStr)
	rows, err := c.db.Query(sqlStr, categoryID, cityID, timestamp, page*numPerRequest, numPerRequest)
	if err != nil {
		log.Println("getPostsByRecentWithCity", err)
		return posts, err
	}
	defer rows.Close()
	for rows.Next() {
		post := postAPI{}
		place := placeAPI{}
		var placeID interface{}
		if err := rows.Scan(
			&post.PostID,
			&post.User.UserID,
			&post.User.Username,
			&post.User.Name,
			&post.User.PhotoURL,
			&post.Content,
			&post.Blob.BlobID,
			&post.Blob.OriginWidth,
			&post.Blob.OriginHeight,
			&post.Type,
			&post.LikeCount,
			&post.DislikeCount,
			&post.CommentCount,
			&placeID,
			&place.CountryCode,
			&place.CityID1,
			&place.CityID2,
			&place.CityID3,
			&place.CityID4,
			&place.CityID5,
			&place.Lat,
			&place.Lon,
			&place.Name,
			&post.CategoryID,
			&post.Createtime,
			&post.Updatetime,
		); err != nil {
			log.Println(err)
			return posts, err
		}
		placeIDCheck, isOK := placeID.(int64)
		if isOK {
			place.PlaceID = placeIDCheck
		}
		post.Place = place
		post.Blob.BlobID = makeBlobURL(post)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return posts, err
	}
	// log.Println("posts", posts)
	return posts, nil
}
func getPostsByRecentWithPlaceLikeNum(countryCode string, categoryID, numPost int) (posts []postAPI, err error) {
	// combine this with func getPostsByRecentWithPlaceLikeNum
	c, err := connectPostgres()
	if err != nil {
		return posts, errors.New("db connection")
	}
	defer c.db.Close()
	timestamp := getNowUnixTimestamp() - twoMonthsInSecond
	if categoryID == categorySup && env != "development" {
		timestamp = getNowUnixTimestamp() - twentyFourHoursInSecond
	}
	sqlStr := `
		SELECT 
		post.post_id,
		post.user_id, xuser.username, xuser.name, xuser.photo_url,
		post.content, post.blob_id, post.origin_width, post.origin_height, post.type, 
		post.like_count, post.dislike_count, post.comment_count,
		place.place_id, place.country_code,
		place.city_id_1, place.city_id_2, place.city_id_3,
		place.city_id_4, place.city_id_5,
		place.lat, place.lon, place.name,
		post.category_id, post.createtime, post.updatetime 
		FROM post 
		JOIN xuser ON xuser.user_id = post.user_id
		JOIN place ON place.place_id = post.place_id
		WHERE post.category_id=$1 AND 
		(place.country_code = $2 OR place.city_id_1 LIKE $2 || '%' OR place.city_id_2 LIKE $2 || '%' OR
		place.city_id_3 LIKE $2 || '%' OR place.city_id_4 LIKE $2 || '%' OR place.city_id_5 LIKE $2 || '%') AND 
		post.createtime >= $3
		ORDER BY post.createtime DESC OFFSET $4 LIMIT $5;
	`
	rows, err := c.db.Query(sqlStr, categoryID, countryCode, timestamp, 0, numPost)
	if err != nil {
		log.Println("getPostsByRecentWithPlaceLikeNum", err)
		return posts, err
	}
	defer rows.Close()
	for rows.Next() {
		post := postAPI{}
		place := placeAPI{}
		var placeID interface{}
		if err := rows.Scan(
			&post.PostID,
			&post.User.UserID,
			&post.User.Username,
			&post.User.Name,
			&post.User.PhotoURL,
			&post.Content,
			&post.Blob.BlobID,
			&post.Blob.OriginWidth,
			&post.Blob.OriginHeight,
			&post.Type,
			&post.LikeCount,
			&post.DislikeCount,
			&post.CommentCount,
			&placeID,
			&place.CountryCode,
			&place.CityID1,
			&place.CityID2,
			&place.CityID3,
			&place.CityID4,
			&place.CityID5,
			&place.Lat,
			&place.Lon,
			&place.Name,
			&post.CategoryID,
			&post.Createtime,
			&post.Updatetime,
		); err != nil {
			log.Println("errrr", err)
			return posts, err
		}
		placeIDCheck, isOK := placeID.(int64)
		if isOK {
			place.PlaceID = placeIDCheck
		}
		post.Place = place
		post.Blob.BlobID = makeBlobURL(post)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return posts, err
	}
	return posts, nil
}
func getPostsByFollowingUsers(userID int64, page int) (posts []postAPI, err error) { // not done yet
	timestamp := getNowUnixTimestamp() - twoMonthsInSecond
	c, err := connectPostgres()
	if err != nil {
		return posts, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		SELECT 
		post.post_id, 
		post.user_id, xuser.username, xuser.name, xuser.photo_url,
		post.content, post.blob_id, post.origin_width, post.origin_height, post.type,
		post.like_count, post.dislike_count, post.comment_count,
		place.place_id, place.country_code,
		place.city_id_1, place.city_id_2, place.city_id_3,
		place.city_id_4, place.city_id_5,
		place.lat, place.lon, place.name,
		post.category_id, post.createtime, post.updatetime 
		FROM post
		JOIN xuser ON xuser.user_id = post.user_id
		JOIN follow ON follow.following_user_id = post.user_id
		JOIN place ON place.place_id = post.place_id
		WHERE follow.follower_user_id= $1 AND post.createtime>=$2
		ORDER BY post.createtime
		DESC OFFSET $3 LIMIT $4;
	`
	rows, err := c.db.Query(sqlStr, userID, timestamp, page*numPerRequest, numPerRequest)
	if err != nil {
		log.Println("getPostsFollowing", err)
		return posts, err
	}
	defer rows.Close()
	for rows.Next() {
		post := postAPI{}
		place := placeAPI{}
		var placeID interface{}
		if err := rows.Scan(
			&post.PostID,
			&post.User.UserID,
			&post.User.Username,
			&post.User.Name,
			&post.User.PhotoURL,
			&post.Content,
			&post.Blob.BlobID,
			&post.Blob.OriginWidth,
			&post.Blob.OriginHeight,
			&post.Type,
			&post.LikeCount,
			&post.DislikeCount,
			&post.CommentCount,
			&placeID,
			&place.CountryCode,
			&place.CityID1,
			&place.CityID2,
			&place.CityID3,
			&place.CityID4,
			&place.CityID5,
			&place.Lat,
			&place.Lon,
			&place.Name,
			&post.CategoryID,
			&post.Createtime,
			&post.Updatetime,
		); err != nil {
			log.Println(err)
			return posts, err
		}
		placeIDCheck, isOK := placeID.(int64)
		if isOK {
			place.PlaceID = placeIDCheck
		}
		post.Place = place
		post.Blob.BlobID = makeBlobURL(post)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return posts, err
	}
	// log.Println("posts", posts)
	return posts, nil
}
func getPostsByFollowingUsersWithCountry(userID int64, countryCode string, categoryID, page int) (posts []postAPI, err error) { // not done yet
	timestamp := getNowUnixTimestamp() - twoMonthsInSecond
	c, err := connectPostgres()
	if err != nil {
		return posts, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		SELECT 
		post.post_id, 
		post.user_id, xuser.username, xuser.name, xuser.photo_url,
		post.content, post.blob_id, post.origin_width, post.origin_height, post.type,
		post.like_count, post.dislike_count, post.comment_count,
		place.place_id, place.country_code,
		place.city_id_1, place.city_id_2, place.city_id_3,
		place.city_id_4, place.city_id_5,
		place.lat, place.lon, place.name,
		post.category_id, post.createtime, post.updatetime 
		FROM post
		JOIN xuser ON xuser.user_id = post.user_id
		JOIN follow ON follow.following_user_id = post.user_id
		JOIN place ON place.place_id = post.place_id
		WHERE follow.follower_user_id= $1 AND post.createtime>=$2 AND 
		post.category_id=$3 AND place.country_code=$4
		ORDER BY post.createtime
		DESC OFFSET $5 LIMIT $6;
	`
	rows, err := c.db.Query(sqlStr, userID, timestamp, categoryID, countryCode, page*numPerRequest, numPerRequest)
	if err != nil {
		log.Println("getPostsByFollowingUsersWithCountry", err)
		return posts, err
	}
	defer rows.Close()
	for rows.Next() {
		post := postAPI{}
		place := placeAPI{}
		var placeID interface{}
		if err := rows.Scan(
			&post.PostID,
			&post.User.UserID,
			&post.User.Username,
			&post.User.Name,
			&post.User.PhotoURL,
			&post.Content,
			&post.Blob.BlobID,
			&post.Blob.OriginWidth,
			&post.Blob.OriginHeight,
			&post.Type,
			&post.LikeCount,
			&post.DislikeCount,
			&post.CommentCount,
			&placeID,
			&place.CountryCode,
			&place.CityID1,
			&place.CityID2,
			&place.CityID3,
			&place.CityID4,
			&place.CityID5,
			&place.Lat,
			&place.Lon,
			&place.Name,
			&post.CategoryID,
			&post.Createtime,
			&post.Updatetime,
		); err != nil {
			log.Println(err)
			return posts, err
		}
		placeIDCheck, isOK := placeID.(int64)
		if isOK {
			place.PlaceID = placeIDCheck
		}
		post.Place = place
		post.Blob.BlobID = makeBlobURL(post)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return posts, err
	}
	// log.Println("posts", posts)
	return posts, nil
}
func getPostsByFollowingUsersWithCity(userID int64, level, cityID string, categoryID, page int) (posts []postAPI, err error) { // not done yet
	timestamp := getNowUnixTimestamp() - twoMonthsInSecond
	c, err := connectPostgres()
	if err != nil {
		return posts, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		SELECT 
		post.post_id, 
		post.user_id, xuser.username, xuser.name, xuser.photo_url,
		post.content, post.blob_id, post.origin_width, post.origin_height, post.type,
		post.like_count, post.dislike_count, post.comment_count,
		place.place_id, place.country_code,
		place.city_id_1, place.city_id_2, place.city_id_3,
		place.city_id_4, place.city_id_5,
		place.lat, place.lon, place.name,
		post.category_id, post.createtime, post.updatetime 
		FROM post
		JOIN xuser ON xuser.user_id = post.user_id
		JOIN follow ON follow.following_user_id = post.user_id
		JOIN place ON place.place_id = post.place_id
		WHERE follow.follower_user_id= $1 AND post.createtime>=$2 AND
		post.category_id=$3 AND place.city_id_` + level + `=$3
		ORDER BY post.createtime
		DESC OFFSET $4 LIMIT $5;
	`
	rows, err := c.db.Query(sqlStr, userID, timestamp, categoryID, cityID, page*numPerRequest, numPerRequest)
	if err != nil {
		log.Println("getPostsByFollowingUsersWithCity", err)
		return posts, err
	}
	defer rows.Close()
	for rows.Next() {
		post := postAPI{}
		place := placeAPI{}
		var placeID interface{}
		if err := rows.Scan(
			&post.PostID,
			&post.User.UserID,
			&post.User.Username,
			&post.User.Name,
			&post.User.PhotoURL,
			&post.Content,
			&post.Blob.BlobID,
			&post.Blob.OriginWidth,
			&post.Blob.OriginHeight,
			&post.Type,
			&post.LikeCount,
			&post.DislikeCount,
			&post.CommentCount,
			&placeID,
			&place.CountryCode,
			&place.CityID1,
			&place.CityID2,
			&place.CityID3,
			&place.CityID4,
			&place.CityID5,
			&place.Lat,
			&place.Lon,
			&place.Name,
			&post.CategoryID,
			&post.Createtime,
			&post.Updatetime,
		); err != nil {
			log.Println(err)
			return posts, err
		}
		placeIDCheck, isOK := placeID.(int64)
		if isOK {
			place.PlaceID = placeIDCheck
		}
		post.Place = place
		post.Blob.BlobID = makeBlobURL(post)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return posts, err
	}
	// log.Println("posts", posts)
	return posts, nil
}
func getPostsByUser(userID int64, page int) (posts []postAPI, err error) { // not done yet
	c, err := connectPostgres()
	if err != nil {
		return posts, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		SELECT 
		post.post_id, 
		post.user_id, xuser.username, xuser.name, xuser.photo_url,
		post.content, post.blob_id, post.origin_width, post.origin_height, post.type,
		post.like_count, post.dislike_count, post.comment_count,
		place.place_id, place.country_code,
		place.city_id_1, place.city_id_2, place.city_id_3,
		place.city_id_4, place.city_id_5,
		place.lat, place.lon, place.name,
		post.category_id, post.createtime, post.updatetime 
		FROM post
		JOIN xuser ON xuser.user_id = post.user_id
		JOIN place ON place.place_id = post.place_id
		WHERE post.user_id= $1
		ORDER BY post.createtime
		DESC OFFSET $2 LIMIT $3;
	`
	rows, err := c.db.Query(sqlStr, userID, page*numPerRequest, numPerRequest)
	if err != nil {
		log.Println("getPostsUser", err)
		return posts, err
	}
	defer rows.Close()
	for rows.Next() {
		post := postAPI{}
		place := placeAPI{}
		var placeID interface{}
		if err := rows.Scan(
			&post.PostID,
			&post.User.UserID,
			&post.User.Username,
			&post.User.Name,
			&post.User.PhotoURL,
			&post.Content,
			&post.Blob.BlobID,
			&post.Blob.OriginWidth,
			&post.Blob.OriginHeight,
			&post.Type,
			&post.LikeCount,
			&post.DislikeCount,
			&post.CommentCount,
			&placeID,
			&place.CountryCode,
			&place.CityID1,
			&place.CityID2,
			&place.CityID3,
			&place.CityID4,
			&place.CityID5,
			&place.Lat,
			&place.Lon,
			&place.Name,
			&post.CategoryID,
			&post.Createtime,
			&post.Updatetime,
		); err != nil {
			log.Println(err)
			return posts, err
		}
		placeIDCheck, isOK := placeID.(int64)
		if isOK {
			place.PlaceID = placeIDCheck
		}
		post.Place = place
		post.Blob.BlobID = makeBlobURL(post)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return posts, err
	}
	// log.Println("posts", posts)
	return posts, nil
}
func getPostByPostIDUserID(postID, userID int64) (count int, err error) {
	c, err := connectPostgres()
	if err != nil {
		return count, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		SELECT 
		COUNT(*)
		FROM post
		WHERE post_id = $1 AND user_id= $2;
	`
	row := c.db.QueryRow(sqlStr, postID, userID)
	err = row.Scan(&count)
	if err != nil {
		log.Println("getPostByPostIDUserID", err)
		return count, err
	}
	return count, err
}
func getPostsByPlaceID(placeID int64, page int) (posts []postAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return posts, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `SELECT lat, lon, name FROM place WHERE place_id=$1;`
	place := placeAPI{PlaceID: placeID}
	row := c.db.QueryRow(sqlStr, placeID)
	if err := row.Scan(
		&place.Lat,
		&place.Lon,
		&place.Name,
	); err != nil {
		log.Println("getPostsByPlaceID", err)
		return posts, errors.New("place not found")
	}
	sqlStr = `
		SELECT 
		post.post_id, 
		post.user_id, xuser.username, xuser.name, xuser.photo_url,
		post.content, post.blob_id, post.origin_width, post.origin_height, post.type,
		post.like_count, post.dislike_count, post.comment_count,
		post.category_id, post.createtime, post.updatetime 
		FROM post
		JOIN xuser ON xuser.user_id = post.user_id
		WHERE post.place_id= $1
		ORDER BY post.createtime
		DESC OFFSET $2 LIMIT $3;
	`
	rows, err := c.db.Query(sqlStr, placeID, page*numPerRequest, numPerRequest)
	if err != nil {
		log.Println("getPostsByPlaceID", err)
		return posts, err
	}
	defer rows.Close()
	for rows.Next() {
		post := postAPI{}
		if err := rows.Scan(
			&post.PostID,
			&post.User.UserID,
			&post.User.Username,
			&post.User.Name,
			&post.User.PhotoURL,
			&post.Content,
			&post.Blob.BlobID,
			&post.Blob.OriginWidth,
			&post.Blob.OriginHeight,
			&post.Type,
			&post.LikeCount,
			&post.DislikeCount,
			&post.CommentCount,
			&post.CategoryID,
			&post.Createtime,
			&post.Updatetime,
		); err != nil {
			log.Println(err)
			return posts, err
		}
		post.Blob.BlobID = makeBlobURL(post)
		post.Place = place
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return posts, err
	}
	// log.Println("posts", posts)
	return posts, nil
}
func getPostsByPopular(userID int64, categoryID, page int) (posts []postAPI, err error) { // not done yet
	c, err := connectMongoDB()
	if err != nil {
		log.Println("mongo session", err)
		return posts, err
	}
	defer c.session.Close()
	cr := connectRedis()
	defer cr.client.Close()
	// read_index
	collection := c.session.DB(mongoDBXociety).C(mongoCollectionPopularPostUserReadIndex)
	ppuri := popularPostUserReadIndexAPI{}
	q := bson.M{"user_id": userID}
	collection.Find(q).One(&ppuri)
	popularPostUserReadIndex := ppuri.CommonPopularPostIndex[categoryID]
	// read
	weekTimestamp := getNowUnixWeekTimestamp()
	// ** you can wrap func as a transaction
	collection = c.session.DB(mongoDBXociety).C(mongoCollectionPostUserRead)
	q = bson.M{"user_id": userID, "category_id": categoryID, "week_timestamp": weekTimestamp}
	r := postUserReadAPI{}
	collection.Find(q).One(&r)
	// popular_common
	collection = c.session.DB(mongoDBXociety).C(mongoCollectionPostCommon)
	q = bson.M{"category_id": categoryID}
	p := postCommonAPI{}
	if err := collection.Find(q).One(&p); err != nil {
		log.Println("getPostsByPopular2", err)
		return posts, err
	}
	count := 0
	for i := popularPostUserReadIndex; i < len(p.PopularPosts); i++ {
		if count >= numPerRequest {
			break
		}
		if r.PopularPosts[p.PopularPosts[i].PostID] == 0 {
			posts = append(posts, p.PopularPosts[i])
			count++
		}
	}
	return posts, nil
}
func getSupPostsByPopularWithCountry(countryCode string, userID int64, page int) (posts []postAPI, err error) { // not done yet
	c, err := connectMongoDB()
	if err != nil {
		log.Println("mongo session", err)
		return posts, err
	}
	defer c.session.Close()
	// read_index
	collection := c.session.DB(mongoDBXociety).C(mongoCollectionPopularPostUserReadIndex)
	ppuri := popularPostUserReadIndexAPI{}
	q := bson.M{"user_id": userID}
	collection.Find(q).One(&ppuri)
	popularPostUserReadIndex := ppuri.CountrySupPopularPostIndex[countryCode]
	// read
	weekTimestamp := getNowUnixWeekTimestamp()
	// ** you can wrap func as a transaction
	collection = c.session.DB(mongoDBXociety).C(mongoCollectionPostUserRead)
	q = bson.M{"user_id": userID, "category_id": categorySup, "week_timestamp": weekTimestamp}
	r := postUserReadAPI{}
	collection.Find(q).One(&r)
	// sup popular post on country
	collection = c.session.DB(mongoDBXociety).C(mongoCollectionCity)
	q = bson.M{"level": "0", "country_code": countryCode}
	p := cityAPI{}
	if err := collection.Find(q).One(&p); err != nil {
		log.Println("getPostsBySupPopularWithCountry2", err)
		return posts, err
	}
	count := 0
	for i := popularPostUserReadIndex; i < len(p.SupPopularPosts); i++ {
		if count >= numPerRequest {
			break
		}
		if r.PopularPosts[p.SupPopularPosts[i].PostID] == 0 {
			posts = append(posts, p.SupPopularPosts[i])
			count++
		}
	}
	return posts, nil
}
func getSupPostsByPopularWithCity(level, cityID string, userID int64, page int) (posts []postAPI, err error) { // not done yet
	c, err := connectMongoDB()
	if err != nil {
		log.Println("mongo session", err)
		return posts, err
	}
	defer c.session.Close()
	// read_index
	collection := c.session.DB(mongoDBXociety).C(mongoCollectionPopularPostUserReadIndex)
	ppuri := popularPostUserReadIndexAPI{}
	q := bson.M{"user_id": userID}
	collection.Find(q).One(&ppuri)
	popularPostUserReadIndex := ppuri.CitySupPopularPostIndex[cityID]
	// read
	weekTimestamp := getNowUnixWeekTimestamp()
	// ** you can wrap func as a transaction
	collection = c.session.DB(mongoDBXociety).C(mongoCollectionPostUserRead)
	q = bson.M{"user_id": userID, "category_id": categorySup, "week_timestamp": weekTimestamp}
	r := postUserReadAPI{}
	collection.Find(q).One(&r)
	// sup popular post on city
	collection = c.session.DB(mongoDBXociety).C(mongoCollectionCity)
	q = bson.M{"level": level, "city_id_" + level: cityID}
	p := cityAPI{}
	if err := collection.Find(q).One(&p); err != nil {
		log.Println("getPostsBySupPopularWithCity", err)
		return posts, err
	}
	count := 0
	for i := popularPostUserReadIndex; i < len(p.SupPopularPosts); i++ {
		if count >= numPerRequest {
			break
		}
		if r.PopularPosts[p.SupPopularPosts[i].PostID] == 0 {
			posts = append(posts, p.SupPopularPosts[i])
			count++
		}
	}
	return posts, nil
}
func getPostsByHashtag(hashtagID int64, page int) (posts []postAPI, err error) { // not done yet
	c, err := connectPostgres()
	if err != nil {
		return posts, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		SELECT 
		post.post_id, 
		post.user_id, xuser.username, xuser.name, xuser.photo_url,
		post.content, post.blob_id, post.origin_width, post.origin_height, post.type,
		post.like_count, post.dislike_count, post.comment_count,
		post.category_id, post.createtime, post.updatetime 
		FROM post
		JOIN xuser ON xuser.user_id = post.user_id
		JOIN post_hashtag ON post_hashtag.post_id = post.post_id
		WHERE post_hashtag.hashtag_id= $1
		ORDER BY post.createtime
		DESC OFFSET $2 LIMIT $3;
	`
	rows, err := c.db.Query(sqlStr, hashtagID, page*numPerRequest, numPerRequest)
	if err != nil {
		log.Println("getPostsByHashtag", err)
		return posts, err
	}
	defer rows.Close()
	for rows.Next() {
		post := postAPI{}
		if err := rows.Scan(
			&post.PostID,
			&post.User.UserID,
			&post.User.Username,
			&post.User.Name,
			&post.User.PhotoURL,
			&post.Content,
			&post.Blob.BlobID,
			&post.Blob.OriginWidth,
			&post.Blob.OriginHeight,
			&post.Type,
			&post.LikeCount,
			&post.DislikeCount,
			&post.CommentCount,
			&post.CategoryID,
			&post.Createtime,
			&post.Updatetime,
		); err != nil {
			log.Println(err)
			return posts, err
		}
		post.Blob.BlobID = makeBlobURL(post)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return posts, err
	}
	// log.Println("posts", posts)
	return posts, nil
}
func getPostsByTag(userID int64, page int) (posts []postAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return posts, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		SELECT 
		post.post_id, 
		post.user_id, xuser.username, xuser.name, xuser.photo_url,
		post.content, post.blob_id, post.origin_width, post.origin_height, post.type,
		post.like_count, post.dislike_count, post.comment_count,
		post.category_id, post.createtime, post.updatetime 
		FROM post_tag_xuser
		JOIN xuser ON xuser.user_id = post_tag_xuser.user_id
		JOIN post ON post_tag_xuser.post_id = post.post_id
		WHERE post_tag_xuser.user_id= $1 AND post_tag_xuser.valid = true
		ORDER BY post.createtime
		DESC OFFSET $2 LIMIT $3;
	`
	rows, err := c.db.Query(sqlStr, userID, page*numPerRequest, numPerRequest)
	if err != nil {
		log.Println("getPostsByTag", err)
		return posts, err
	}
	defer rows.Close()
	for rows.Next() {
		post := postAPI{}
		if err := rows.Scan(
			&post.PostID,
			&post.User.UserID,
			&post.User.Username,
			&post.User.Name,
			&post.User.PhotoURL,
			&post.Content,
			&post.Blob.BlobID,
			&post.Blob.OriginWidth,
			&post.Blob.OriginHeight,
			&post.Type,
			&post.LikeCount,
			&post.DislikeCount,
			&post.CommentCount,
			&post.CategoryID,
			&post.Createtime,
			&post.Updatetime,
		); err != nil {
			log.Println(err)
			return posts, err
		}
		post.Blob.BlobID = makeBlobURL(post)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return posts, err
	}
	// log.Println("posts", posts)
	return posts, nil
}

// hashtags
func getHashtags(value string, page int) (hashtags []hashtagAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return hashtags, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		SELECT 
		hashtag_id, value, count
		FROM hashtag
		WHERE value LIKE $1
		ORDER BY hashtag.count
		DESC OFFSET $2 LIMIT $3;
	`
	rows, err := c.db.Query(sqlStr, value+"%", page*numPerRequest, numPerRequest)
	if err != nil {
		log.Println("getHashtags", err)
		return hashtags, err
	}
	defer rows.Close()
	for rows.Next() {
		hashtag := hashtagAPI{}
		if err := rows.Scan(
			&hashtag.HashtagID,
			&hashtag.Value,
			&hashtag.Count,
		); err != nil {
			log.Println(err)
			return hashtags, err
		}
		hashtags = append(hashtags, hashtag)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return hashtags, err
	}
	// log.Println("hashtags", hashtags)
	return hashtags, nil
}

// tags
func getAllTagsByPost(postID int64) (tags []tagOnPostAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return tags, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		SELECT 
		post_tag_xuser.post_id,
		post_tag_xuser.user_id, xuser.username, xuser.name,
		post_tag_xuser.x, post_tag_xuser.y,
		post_tag_xuser.valid,
		post_tag_xuser.createtime, post_tag_xuser.updatetime
		FROM post_tag_xuser
		JOIN xuser ON xuser.user_id = post_tag_xuser.user_id
		WHERE post_tag_xuser.post_id = $1 AND post_tag_xuser.valid = true;
	`
	rows, err := c.db.Query(sqlStr, postID)
	if err != nil {
		log.Println("getAllTagsByPost", err)
		return tags, err
	}
	defer rows.Close()
	for rows.Next() {
		tag := tagOnPostAPI{}
		if err := rows.Scan(
			&tag.PostID,
			&tag.User.UserID,
			&tag.User.Username,
			&tag.User.Name,
			&tag.X,
			&tag.Y,
			&tag.Valid,
			&tag.Createtime,
			&tag.Updatetime,
		); err != nil {
			log.Println(err)
			return tags, err
		}
		tags = append(tags, tag)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return tags, err
	}
	// log.Println("tags", hashtags)
	return tags, nil
}

// comment
func getCommentsOnPost(postID int64, page int) (comments []commentAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return comments, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		SELECT 
		comment.comment_id, 
		comment.user_id, xuser.username, xuser.name, 
		comment.comment, 
		comment.like_count, comment.dislike_count, comment.reply_count,
		comment.createtime, comment.updatetime 
		FROM comment JOIN xuser ON comment.user_id = xuser.user_id 
		WHERE post_id=$1 
		ORDER BY createtime DESC OFFSET $2 LIMIT $3;
	`
	rows, err := c.db.Query(sqlStr, postID, page*numPerRequest, numPerRequest)
	if err != nil {
		log.Println("getCommentsPost", err)
	}
	defer rows.Close()
	for rows.Next() {
		comment := commentAPI{}
		comment.PostID = postID
		if err := rows.Scan(
			&comment.CommentID,
			&comment.User.UserID,
			&comment.User.Username,
			&comment.User.Name,
			&comment.Comment,
			&comment.LikeCount,
			&comment.DislikeCount,
			&comment.ReplyCount,
			&comment.Createtime,
			&comment.Updatetime,
		); err != nil {
			log.Println("errrr", err)
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
	}
	return comments, nil
}

// reply
func getRepliesOnComment(commentID int64, page int) (replies []replyAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return replies, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		SELECT 
		reply.reply_id, 
		reply.user_id, xuser.username, xuser.name, 
		reply.reply, 
		reply.like_count, reply.dislike_count,
		reply.createtime, reply.updatetime 
		FROM reply JOIN xuser ON reply.user_id = xuser.user_id 
		WHERE comment_id=$1 
		ORDER BY createtime DESC OFFSET $2 LIMIT $3;
	`
	rows, err := c.db.Query(sqlStr, commentID, page*numPerRequest, numPerRequest)
	if err != nil {
		log.Println("getRepliesOnComment", err)
	}
	defer rows.Close()
	for rows.Next() {
		reply := replyAPI{}
		reply.CommentID = commentID
		if err := rows.Scan(
			&reply.ReplyID,
			&reply.User.UserID,
			&reply.User.Username,
			&reply.User.Name,
			&reply.Reply,
			&reply.LikeCount,
			&reply.DislikeCount,
			&reply.Createtime,
			&reply.Updatetime,
		); err != nil {
			log.Println("errrr", err)
		}
		replies = append(replies, reply)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
	}
	return replies, nil
}

// reaction
func getReactionsOnPost(postID int64, page int) (reactionsOnPost []reactionOnPostAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return reactionsOnPost, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		SELECT 
		post_reaction.post_id, 
		post_reaction.user_id, xuser.username, xuser.name, 
		post_reaction.reaction_id, post_reaction.createtime 
		FROM public.post_reaction JOIN xuser on post_reaction.user_id = xuser.user_id
		WHERE post_id=$1
		ORDER BY createtime DESC OFFSET $2 LIMIT $3;
	`
	rows, err := c.db.Query(sqlStr, postID, page*numPerRequest, numPerRequest)
	if err != nil {
		log.Println("getReactionsOnPost", err)
	}
	defer rows.Close()
	for rows.Next() {
		reactionOnPost := reactionOnPostAPI{}
		if err := rows.Scan(
			&reactionOnPost.PostID,
			&reactionOnPost.User.UserID,
			&reactionOnPost.User.Username,
			&reactionOnPost.User.Name,
			&reactionOnPost.ReactionID,
			&reactionOnPost.Createtime,
		); err != nil {
			log.Println("errrr", err)
		}
		reactionsOnPost = append(reactionsOnPost, reactionOnPost)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
	}
	return reactionsOnPost, nil
}
func getReactionsOnComment(commentID int64, page int) (reactionsOnComment []reactionOnCommentAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return reactionsOnComment, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		SELECT 
		comment_reaction.comment_id, 
		comment_reaction.user_id, xuser.username, xuser.name, 
		comment_reaction.reaction_id, comment_reaction.createtime 
		FROM public.comment_reaction JOIN xuser on comment_reaction.user_id = xuser.user_id
		WHERE comment_id=$1
		ORDER BY createtime DESC OFFSET $2 LIMIT $3;
	`
	rows, err := c.db.Query(sqlStr, commentID, page*numPerRequest, numPerRequest)
	if err != nil {
		log.Println("getReactionsOnComment", err)
	}
	defer rows.Close()
	for rows.Next() {
		reactionOnComment := reactionOnCommentAPI{}
		if err := rows.Scan(
			&reactionOnComment.CommentID,
			&reactionOnComment.User.UserID,
			&reactionOnComment.User.Username,
			&reactionOnComment.User.Name,
			&reactionOnComment.ReactionID,
			&reactionOnComment.Createtime,
		); err != nil {
			log.Println("errrr", err)
		}
		reactionsOnComment = append(reactionsOnComment, reactionOnComment)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
	}
	return reactionsOnComment, nil
}
func getReactionsOnReply(replyID int64, page int) (reactionsOnReply []reactionOnReplyAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return reactionsOnReply, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		SELECT 
		reply_reaction.reply_id, 
		reply_reaction.user_id, xuser.username, xuser.name, 
		reply_reaction.reaction_id, reply_reaction.createtime 
		FROM public.reply_reaction JOIN xuser on reply_reaction.user_id = xuser.user_id
		WHERE reply_id=$1
		ORDER BY createtime DESC OFFSET $2 LIMIT $3;
	`
	rows, err := c.db.Query(sqlStr, replyID, page*numPerRequest, numPerRequest)
	if err != nil {
		log.Println("getReactionsOnReply", err)
	}
	defer rows.Close()
	for rows.Next() {
		reactionOnReply := reactionOnReplyAPI{}
		if err := rows.Scan(
			&reactionOnReply.ReplyID,
			&reactionOnReply.User.UserID,
			&reactionOnReply.User.Username,
			&reactionOnReply.User.Name,
			&reactionOnReply.ReactionID,
			&reactionOnReply.Createtime,
		); err != nil {
			log.Println("errrr", err)
		}
		reactionsOnReply = append(reactionsOnReply, reactionOnReply)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
	}
	return reactionsOnReply, nil
}

// [mutation]
func userInsert(user userDB) (userID int64, err error) {
	c, err := connectPostgres()
	if err != nil {
		return userID, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		INSERT INTO xuser 
		(
			username, email, password, name, phone, 
			gender, bio, credit, photo_url, 
			language_id, country_code, 
			timezone, last_ip, createtime, updatetime 
		) VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) RETURNING user_id;
	`
	err = c.db.QueryRow(sqlStr,
		user.Username,
		user.Email,
		user.Password,
		user.Name,
		user.Phone,
		user.Gender,
		user.Bio,
		user.Credit,
		user.PhotoURL,
		user.LanguageID,
		user.CountryCode,
		user.Timezone,
		user.LastIP,
		user.Createtime,
		user.Updatetime,
	).Scan(&userID)
	if err != nil {
		// log.Println(err)
		return userID, err
	}
	return userID, err
}

// follow
func follow(followingUserID, followerUserID int64) (us updateStatusAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return us, errors.New("db connection")
	}
	defer c.db.Close()
	timestamp := getNowUnixTimestamp()
	sqlStr := `
		INSERT INTO follow 
		(following_user_id, follower_user_id, valid, createtime, updatetime) 
		values($1,$2,$3,$4,$5) returning createtime;
	`
	res, err := c.db.Exec(sqlStr, followingUserID, followerUserID, true, timestamp, timestamp)
	if err != nil {
		return us, errors.New("you've followed this user")
	}
	count, err := res.RowsAffected()
	if err != nil {
		return us, err
	}
	us.RowsAffected = int(count)
	return us, nil
}
func unfollow(followingUserID, followerUserID int64) (us updateStatusAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return us, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		DELETE FROM follow 
		WHERE following_user_id=$1 AND follower_user_id=$2;
	`
	res, err := c.db.Exec(sqlStr,
		followingUserID, followerUserID)
	if err != nil {
		return us, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return us, err
	}
	us.RowsAffected = int(count)
	return us, nil
}

// place
func placeInsert(place placeAPI) (placeID int64, err error) {
	c, err := connectPostgres()
	if err != nil {
		return placeID, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		INSERT INTO place 
		(country_code, city_id_1, city_id_2, city_id_3, city_id_4, city_id_5, lat, lon, name) 
		values($1,$2,$3,$4,$5,$6,$7,$8,$9) returning place_id;
	`
	err = c.db.QueryRow(sqlStr,
		place.CountryCode,
		place.CityID1,
		place.CityID2,
		place.CityID3,
		place.CityID4,
		place.CityID5,
		place.Lat,
		place.Lon,
		place.Name,
	).Scan(&placeID)
	if err != nil {
		// log.Println(err)
		return placeID, err
	}
	return placeID, nil
}

// post
func postInsert(post postAPI) (postID int64, err error) {
	c, err := connectPostgres()
	if err != nil {
		return postID, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		INSERT INTO post 
		(
			user_id, content, blob_id, origin_width, origin_height, type, 
			like_count, dislike_count, comment_count, place_id, 
			category_id, public, createtime, updatetime
		) VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING post_id;
	`
	var placeID interface{} // which is nil => null in postgres
	if post.Place.PlaceID > 0 {
		placeID = post.Place.PlaceID
	}
	err = c.db.QueryRow(sqlStr,
		post.User.UserID,
		post.Content,
		post.Blob.BlobID,
		post.Blob.OriginWidth,
		post.Blob.OriginHeight,
		post.Type,
		post.LikeCount,
		post.DislikeCount,
		post.CommentCount,
		placeID,
		post.CategoryID,
		post.Public,
		post.Createtime,
		post.Updatetime,
	).Scan(&postID)
	if err != nil {
		// log.Println(err)
		return postID, err
	}
	return postID, err
}
func postUpdate(post postAPI) (us updateStatusAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return us, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		UPDATE post 
		SET content=$1, place_id=$2,
		category_id=$3, updatetime=$4
		WHERE post_id=$5 AND user_id=$6;
	`
	var placeID interface{} // which is nil
	if post.Place.PlaceID > 0 {
		placeID = post.Place.PlaceID
	}
	res, err := c.db.Exec(sqlStr,
		post.Content, placeID,
		post.CategoryID, post.Updatetime,
		post.PostID, post.User.UserID)
	if err != nil {
		return us, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return us, err
	}
	us.RowsAffected = int(count)
	return us, nil
}
func postDelete(post postAPI) (us updateStatusAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return us, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		DELETE FROM post
		WHERE post_id=$1 AND user_id=$2;
	`
	res, err := c.db.Exec(sqlStr, post.PostID, post.User.UserID)
	if err != nil {
		return us, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return us, err
	}
	us.RowsAffected = int(count)
	return us, nil
}
func postPopularRead(categoryID int, indexRead int, userID int64) (posts []postAPI, err error) {
	c, err := connectMongoDB()
	if err != nil {
		log.Println("mongo session", err)
		return posts, err
	}
	defer c.session.Close()
	// read_index
	collection := c.session.DB(mongoDBXociety).C(mongoCollectionPopularPostUserReadIndex)
	ppuri := popularPostUserReadIndexAPI{}
	q := bson.M{"user_id": userID}
	collection.Find(q).One(&ppuri)
	popularPostUserReadIndex := ppuri.CommonPopularPostIndex[categoryID]
	// read
	weekTimestamp := getNowUnixWeekTimestamp()
	postsRead := make(map[int64]int)
	// ** you can wrap func as a transaction
	collection = c.session.DB(mongoDBXociety).C(mongoCollectionPostUserRead)
	for t := weekTimestamp; t > weekTimestamp-2*sevenDaysInSecond; t -= sevenDaysInSecond {
		u := postUserReadAPI{}
		q := bson.M{"user_id": userID, "category_id": categoryID, "week_timestamp": t}
		if err := collection.Find(q).One(&u); err == nil {
			for k, v := range u.PopularPosts {
				postsRead[k] = v
			}
		}
	}
	// post_common
	collection = c.session.DB(mongoDBXociety).C(mongoCollectionPostCommon)
	u := postCommonAPI{}
	q = bson.M{"category_id": categoryID}
	if err := collection.Find(q).One(&u); err != nil {
		log.Println("postPopularRead", err)
		return posts, err
	}
	timestamp := getNowUnixTimestamp()
	postsReadNew := make(map[int64]int)
	count := 0
	if indexRead >= 0 {
		last := popularPostUserReadIndex + indexRead + 1
		if last >= len(u.PopularPosts) {
			last = len(u.PopularPosts)
		}
		for i := popularPostUserReadIndex; i < last; i++ {
			postsReadNew[u.PopularPosts[i].PostID] = timestamp // for record read post
		}
	}
	if len(postsReadNew) > 0 {
		collection = c.session.DB(mongoDBXociety).C(mongoCollectionPostUserRead)
		if _, err := collection.Upsert(
			bson.M{"user_id": userID, "category_id": categoryID, "week_timestamp": weekTimestamp},
			bson.M{"$set": parsePopularPostReadObjectMongo(postsReadNew)}); err != nil {
			log.Println("upsert", err)
		}
	}
	for k, v := range postsReadNew {
		postsRead[k] = v
	}
	for i := popularPostUserReadIndex + indexRead + 1; i < len(u.PopularPosts); i++ {
		if count > numPerRequest {
			break
		}
		if postsRead[u.PopularPosts[i].PostID] == 0 {
			posts = append(posts, u.PopularPosts[i]) // for current query popular post list
			count++
		}
	}
	if count > 0 {
		collection = c.session.DB(mongoDBXociety).C(mongoCollectionPopularPostUserReadIndex)
		if _, err := collection.Upsert(
			bson.M{"user_id": userID},
			bson.M{"$set": bson.M{
				"common_popular_post_index." + strconv.Itoa(categoryID): popularPostUserReadIndex + indexRead + 1,
			}}); err != nil {
			log.Println("upsert", err)
		}
	}
	return posts, nil
}
func supPostPopularReadCountry(countryCode string, indexRead int, userID int64) (posts []postAPI, err error) {
	c, err := connectMongoDB()
	if err != nil {
		log.Println("mongo session", err)
		return posts, err
	}
	defer c.session.Close()
	// read_index
	collection := c.session.DB(mongoDBXociety).C(mongoCollectionPopularPostUserReadIndex)
	ppuri := popularPostUserReadIndexAPI{}
	q := bson.M{"user_id": userID}
	collection.Find(q).One(&ppuri)
	popularPostUserReadIndex := ppuri.CountrySupPopularPostIndex[countryCode]
	// read
	weekTimestamp := getNowUnixWeekTimestamp()
	postsRead := make(map[int64]int)
	// ** you can wrap func as a transaction
	collection = c.session.DB(mongoDBXociety).C(mongoCollectionPostUserRead)
	for t := weekTimestamp; t > weekTimestamp-2*sevenDaysInSecond; t -= sevenDaysInSecond {
		u := postUserReadAPI{}
		q := bson.M{"user_id": userID, "category_id": categorySup, "week_timestamp": t}
		if err := collection.Find(q).One(&u); err == nil {
			for k, v := range u.PopularPosts {
				postsRead[k] = v
			}
		}
	}
	// post_common
	collection = c.session.DB(mongoDBXociety).C(mongoCollectionCity)
	u := cityAPI{}
	q = bson.M{"level": "0", "country_code": countryCode}
	if err := collection.Find(q).One(&u); err != nil {
		log.Println("postSupPopularReadCountry", err)
		return posts, err
	}
	timestamp := getNowUnixTimestamp()
	postsReadNew := make(map[int64]int)
	count := 0
	if indexRead >= 0 {
		last := popularPostUserReadIndex + indexRead + 1
		if last >= len(u.SupPopularPosts) {
			last = len(u.SupPopularPosts)
		}
		for i := popularPostUserReadIndex; i < last; i++ {
			postsReadNew[u.SupPopularPosts[i].PostID] = timestamp // for record read post
		}
	}
	if len(postsReadNew) > 0 {
		collection = c.session.DB(mongoDBXociety).C(mongoCollectionPostUserRead)
		if _, err := collection.Upsert(
			bson.M{"user_id": userID, "category_id": categorySup, "week_timestamp": weekTimestamp},
			bson.M{"$set": parsePopularPostReadObjectMongo(postsReadNew)}); err != nil {
			log.Println("upsert", err)
		}
	}
	for k, v := range postsReadNew {
		postsRead[k] = v
	}
	for i := popularPostUserReadIndex + indexRead + 1; i < len(u.SupPopularPosts); i++ {
		if count > numPerRequest {
			break
		}
		if postsRead[u.SupPopularPosts[i].PostID] == 0 {
			posts = append(posts, u.SupPopularPosts[i]) // for current query popular post list
			count++
		}
	}
	if count > 0 {
		collection = c.session.DB(mongoDBXociety).C(mongoCollectionPopularPostUserReadIndex)
		if _, err := collection.Upsert(
			bson.M{"user_id": userID},
			bson.M{"$set": bson.M{
				"country_sup_popular_post_index." + countryCode: popularPostUserReadIndex + indexRead + 1,
			}}); err != nil {
			log.Println("upsert", err)
		}
	}
	return posts, nil
}
func supPostPopularReadCity(level, cityID string, indexRead int, userID int64) (posts []postAPI, err error) {
	c, err := connectMongoDB()
	if err != nil {
		log.Println("mongo session", err)
		return posts, err
	}
	defer c.session.Close()
	// read_index
	collection := c.session.DB(mongoDBXociety).C(mongoCollectionPopularPostUserReadIndex)
	ppuri := popularPostUserReadIndexAPI{}
	q := bson.M{"user_id": userID}
	collection.Find(q).One(&ppuri)
	popularPostUserReadIndex := ppuri.CitySupPopularPostIndex[cityID]
	// read
	weekTimestamp := getNowUnixWeekTimestamp()
	postsRead := make(map[int64]int)
	// ** you can wrap func as a transaction
	collection = c.session.DB(mongoDBXociety).C(mongoCollectionPostUserRead)
	for t := weekTimestamp; t > weekTimestamp-2*sevenDaysInSecond; t -= sevenDaysInSecond {
		u := postUserReadAPI{}
		q := bson.M{"user_id": userID, "category_id": categorySup, "week_timestamp": t}
		if err := collection.Find(q).One(&u); err == nil {
			for k, v := range u.PopularPosts {
				postsRead[k] = v
			}
		}
	}
	// post_common
	collection = c.session.DB(mongoDBXociety).C(mongoCollectionCity)
	u := cityAPI{}
	q = bson.M{"level": level, "city_id_" + level: cityID}
	if err := collection.Find(q).One(&u); err != nil {
		log.Println("postPopularRead", err)
		return posts, err
	}
	timestamp := getNowUnixTimestamp()
	postsReadNew := make(map[int64]int)
	count := 0
	if indexRead >= 0 {
		last := popularPostUserReadIndex + indexRead + 1
		if last >= len(u.SupPopularPosts) {
			last = len(u.SupPopularPosts)
		}
		for i := popularPostUserReadIndex; i < last; i++ {
			postsReadNew[u.SupPopularPosts[i].PostID] = timestamp // for record read post
		}
	}
	if len(postsReadNew) > 0 {
		collection = c.session.DB(mongoDBXociety).C(mongoCollectionPostUserRead)
		if _, err := collection.Upsert(
			bson.M{"user_id": userID, "category_id": categorySup, "week_timestamp": weekTimestamp},
			bson.M{"$set": parsePopularPostReadObjectMongo(postsReadNew)}); err != nil {
			log.Println("upsert", err)
		}
	}
	for k, v := range postsReadNew {
		postsRead[k] = v
	}
	for i := popularPostUserReadIndex + indexRead + 1; i < len(u.SupPopularPosts); i++ {
		if count > numPerRequest {
			break
		}
		if postsRead[u.SupPopularPosts[i].PostID] == 0 {
			posts = append(posts, u.SupPopularPosts[i]) // for current query popular post list
			count++
		}
	}
	if count > 0 {
		collection = c.session.DB(mongoDBXociety).C(mongoCollectionPopularPostUserReadIndex)
		if _, err := collection.Upsert(
			bson.M{"user_id": userID},
			bson.M{"$set": bson.M{
				"city_sup_popular_post_index." + cityID: popularPostUserReadIndex + indexRead + 1,
			}}); err != nil {
			log.Println("upsert", err)
		}
	}
	return posts, nil
}

// hashtags
func hashtagInsert(hashtags []string) (hashtagsID []int64, err error) {
	c, err := connectPostgres()
	if err != nil {
		return hashtagsID, errors.New("db connection")
	}
	defer c.db.Close()
	for i := 0; i < len(hashtags); i++ {
		hashtagID := int64(0)
		sqlStr := `
			UPDATE hashtag SET count = hashtag.count + 1 WHERE hashtag.value=$1 RETURNING hashtag_id
		`
		_ = c.db.QueryRow(sqlStr,
			hashtags[i],
		).Scan(&hashtagID)
		// if err != nil {
		// 	log.Println("this part means that hashtag value's not existed", err)
		// }
		if hashtagID == 0 {
			sqlStr = `
				INSERT INTO hashtag (value, count) VALUES($1, 1) RETURNING hashtag_id
			`
			err = c.db.QueryRow(sqlStr,
				hashtags[i],
			).Scan(&hashtagID)
			if err != nil {
				return hashtagsID, err
			}
		}
		hashtagsID = append(hashtagsID, hashtagID)
	}
	return hashtagsID, err
}
func hashtagOnPostSet(postID int64, hashtagsID []int64) (us updateStatusAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return us, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStrInsert, sqlStrDelete, args := parseHashtagOnPostSQL(postID, hashtagsID)
	res, err := c.db.Exec(sqlStrInsert, args...)
	if err != nil {
		return us, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return us, err
	}
	us.RowsAffected = int(count)
	// delete not current use
	_, err = c.db.Exec(sqlStrDelete, args...)
	if err != nil {
		return us, err
	}
	return us, nil
}

// ** you can wrap func as a transaction

// post_tags
func tagOnPostUpdate(postID int64, tag tagOnPostSetAPI) (us updateStatusAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return us, errors.New("db connection")
	}
	defer c.db.Close()
	timestamp := getNowUnixTimestamp()
	sqlStr := `
		UPDATE post_tag_xuser 
		SET
		x=$1, y=$2,
		updatetime=$3
		WHERE post_id=$4 AND user_id=$5;
	`
	res, err := c.db.Exec(sqlStr,
		tag.X, tag.Y,
		timestamp,
		postID, tag.UserID,
	)
	if err != nil {
		return us, err
	}
	count, err := res.RowsAffected()
	// if err != nil {
	// 	return us, err
	// }
	if count == 0 {
		sqlStr = `
			INSERT INTO post_tag_xuser (
				post_id, user_id,
				x, y,
				valid,
				createtime, updatetime
			)
			VALUES ($1,$2,$3,$4,$5,$6,$7);
		`
		res, err := c.db.Exec(sqlStr,
			postID, tag.UserID,
			tag.X, tag.Y,
			false,
			timestamp, timestamp,
		)
		if err != nil {
			return us, err
		}
		count, err = res.RowsAffected()
		if err != nil {
			return us, err
		}
	}
	us.RowsAffected = int(count)
	return us, nil
}
func tagsOnPostSet(postID int64, tags []tagOnPostSetAPI) (us updateStatusAPI, err error) {
	if len(tags) == 0 {
		return us, nil
	}
	c, err := connectPostgres()
	if err != nil {
		return us, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr, args := parseTagOnPostInsertSQL(postID, tags)
	res, err := c.db.Exec(sqlStr, args...)
	if err != nil {
		return us, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return us, err
	}
	us.RowsAffected = int(count)
	return us, nil
}
func tagOnPostConfirm(postID, userID int64) (us updateStatusAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return us, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		UPDATE post_tag_xuser 
		SET
		valid=$1 
		WHERE post_id=$2 AND user_id=$3;
	`
	res, err := c.db.Exec(sqlStr,
		true,
		postID, userID,
	)
	if err != nil {
		return us, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return us, err
	}
	us.RowsAffected = int(count)
	return us, nil
}
func tagOnPostDelete(postID, userID int64) (us updateStatusAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return us, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		DELETE FROM post_tag_xuser 
		WHERE post_id=$1 AND user_id=$2;
	`
	res, err := c.db.Exec(sqlStr,
		postID, userID)
	if err != nil {
		return us, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return us, err
	}
	us.RowsAffected = int(count)
	return us, nil
}

// comment
func commentOnPostInsert(comment commentAPI) (commentID int64, err error) {
	c, err := connectPostgres()
	if err != nil {
		return commentID, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		INSERT INTO comment 
		(
			post_id, user_id, comment, 
			like_count, dislike_count, reply_count,
			createtime, updatetime
		) VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8) RETURNING comment_id;
	`
	err = c.db.QueryRow(sqlStr,
		comment.PostID,
		comment.User.UserID,
		comment.Comment,
		comment.LikeCount,
		comment.DislikeCount,
		comment.ReplyCount,
		comment.Createtime,
		comment.Updatetime,
	).Scan(&commentID)
	if err != nil {
		log.Println(err)
	}
	// update post.comment_count
	sqlStr = `
		UPDATE post
		SET comment_count =
		(SELECT COUNT(*) FROM comment
		WHERE comment.post_id = post.post_id AND comment.post_id = $1) WHERE post_id = $1;
	`
	_, err = c.db.Exec(sqlStr, comment.PostID)
	if err != nil {
		return commentID, err
	}
	return commentID, nil
}
func commentOnPostUpdate(comment commentAPI) (us updateStatusAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return us, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		UPDATE comment 
		SET comment=$1, updatetime=$2
		WHERE comment_id=$3 AND user_id=$4;
	`
	res, err := c.db.Exec(sqlStr,
		comment.Comment, comment.Updatetime,
		comment.CommentID, comment.User.UserID,
	)
	if err != nil {
		return us, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return us, err
	}
	us.RowsAffected = int(count)
	return us, nil
}
func commentOnPostDelete(comment commentAPI) (us updateStatusAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return us, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		DELETE FROM comment 
		WHERE comment_id=$1 AND user_id=$2;
	`
	res, err := c.db.Exec(sqlStr,
		comment.CommentID, comment.User.UserID)
	if err != nil {
		return us, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return us, err
	}
	us.RowsAffected = int(count)
	// update post.comment_count
	sqlStr = `
		UPDATE post
		SET comment_count =
		(SELECT COUNT(*) FROM comment
		WHERE comment.post_id = post.post_id AND comment.post_id = $1) WHERE post_id = $1;
	`
	_, err = c.db.Exec(sqlStr, comment.PostID)
	if err != nil {
		return us, err
	}
	return us, nil
}

// reply
func replyOnCommentInsert(reply replyAPI) (replyID int64, err error) {
	c, err := connectPostgres()
	if err != nil {
		return replyID, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		INSERT INTO reply 
		(
			comment_id, user_id, reply, 
			like_count, dislike_count,
			createtime, updatetime
		) VALUES 
		($1, $2, $3, $4, $5, $6, $7) RETURNING reply_id;
	`
	err = c.db.QueryRow(sqlStr,
		reply.CommentID,
		reply.User.UserID,
		reply.Reply,
		reply.LikeCount,
		reply.DislikeCount,
		reply.Createtime,
		reply.Updatetime,
	).Scan(&replyID)
	if err != nil {
		log.Println(err)
	}
	// update comment.reply_count
	sqlStr = `
		UPDATE comment
		SET reply_count =
		(SELECT COUNT(*) FROM reply
		WHERE reply.comment_id = comment.comment_id AND comment.comment_id = $1) WHERE comment_id = $1;
	`
	_, err = c.db.Exec(sqlStr, reply.CommentID)
	if err != nil {
		return replyID, err
	}
	return replyID, nil
}
func replyOnCommentUpdate(reply replyAPI) (us updateStatusAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return us, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		UPDATE reply 
		SET reply=$1, updatetime=$2
		WHERE reply_id=$3 AND user_id=$4;
	`
	res, err := c.db.Exec(sqlStr,
		reply.Reply, reply.Updatetime,
		reply.ReplyID, reply.User.UserID,
	)
	if err != nil {
		return us, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return us, err
	}
	us.RowsAffected = int(count)
	return us, nil
}
func replyOnCommentDelete(reply replyAPI) (us updateStatusAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return us, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		DELETE FROM reply 
		WHERE reply_id=$1 AND user_id=$2;
	`
	res, err := c.db.Exec(sqlStr,
		reply.ReplyID, reply.User.UserID)
	if err != nil {
		return us, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return us, err
	}
	us.RowsAffected = int(count)
	// update comment.reply_count
	sqlStr = `
		UPDATE comment
		SET reply_count =
		(SELECT COUNT(*) FROM reply
		WHERE reply.comment_id = comment.comment_id AND comment.comment_id = $1) WHERE comment_id = $1;
	`
	_, err = c.db.Exec(sqlStr, reply.CommentID)
	if err != nil {
		return us, err
	}
	return us, nil
}

// reaction
func reactionOnPostSet(reactionOnPost reactionOnPostAPI) (us updateStatusAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return us, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		INSERT INTO post_reaction (
			post_id,user_id,reaction_id,createtime
		) 
		VALUES($1,$2,$3,$4) 
		ON CONFLICT ON CONSTRAINT post_reaction_post_user_unique DO 
		UPDATE SET post_id=$1, user_id=$2, reaction_id=$3, createtime=$4;
	`
	res, err := c.db.Exec(sqlStr,
		reactionOnPost.PostID, reactionOnPost.User.UserID, reactionOnPost.ReactionID, reactionOnPost.Createtime)
	if err != nil {
		return us, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return us, err
	}
	us.RowsAffected = int(count)
	// update post.like_count || post.dislike_count
	sqlStr = parseReactionCountSQL("post_reaction", "reaction_id", "post", "post_id")
	_, err = c.db.Exec(sqlStr, reactionOnPost.PostID)
	if err != nil {
		return us, err
	}
	return us, nil
}
func reactionOnPostDelete(reactionOnPost reactionOnPostAPI) (us updateStatusAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return us, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		DELETE FROM post_reaction 
		WHERE post_id=$1 AND user_id=$2;
	`
	res, err := c.db.Exec(sqlStr,
		reactionOnPost.PostID, reactionOnPost.User.UserID)
	if err != nil {
		return us, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return us, err
	}
	us.RowsAffected = int(count)
	// update post.like_count || post.dislike_count
	sqlStr = parseReactionCountSQL("post_reaction", "reaction_id", "post", "post_id")
	_, err = c.db.Exec(sqlStr, reactionOnPost.PostID)
	if err != nil {
		return us, err
	}
	return us, nil
}

func reactionOnCommentSet(reactionOnComment reactionOnCommentAPI) (us updateStatusAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return us, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		INSERT INTO comment_reaction (
			comment_id,user_id,reaction_id,createtime
		) 
		VALUES($1,$2,$3,$4) 
		ON CONFLICT ON CONSTRAINT comment_reaction_comment_user_unique DO 
		UPDATE SET comment_id=$1, user_id=$2, reaction_id=$3, createtime=$4;
	`
	res, err := c.db.Exec(sqlStr,
		reactionOnComment.CommentID, reactionOnComment.User.UserID, reactionOnComment.ReactionID, reactionOnComment.Createtime)
	if err != nil {
		return us, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return us, err
	}
	us.RowsAffected = int(count)
	// update comment.like_count || comment.dislike_count
	sqlStr = parseReactionCountSQL("comment_reaction", "reaction_id", "comment", "comment_id")
	_, err = c.db.Exec(sqlStr, reactionOnComment.CommentID)
	if err != nil {
		return us, err
	}
	return us, nil
}
func reactionOnCommentDelete(reactionOnComment reactionOnCommentAPI) (us updateStatusAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return us, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		DELETE FROM comment_reaction 
		WHERE comment_id=$1 AND user_id=$2;
	`
	res, err := c.db.Exec(sqlStr,
		reactionOnComment.CommentID, reactionOnComment.User.UserID)
	if err != nil {
		return us, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return us, err
	}
	us.RowsAffected = int(count)
	// update comment.like_count || comment.dislike_count
	sqlStr = parseReactionCountSQL("comment_reaction", "reaction_id", "comment", "comment_id")
	_, err = c.db.Exec(sqlStr, reactionOnComment.CommentID)
	if err != nil {
		return us, err
	}
	return us, nil
}

func reactionOnReplySet(reactionOnReply reactionOnReplyAPI) (us updateStatusAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return us, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		INSERT INTO reply_reaction (
			reply_id,user_id,reaction_id,createtime
		) 
		VALUES($1,$2,$3,$4) 
		ON CONFLICT ON CONSTRAINT reply_reaction_reply_user_unique DO 
		UPDATE SET reply_id=$1, user_id=$2, reaction_id=$3, createtime=$4;
	`
	res, err := c.db.Exec(sqlStr,
		reactionOnReply.ReplyID, reactionOnReply.User.UserID, reactionOnReply.ReactionID, reactionOnReply.Createtime)
	if err != nil {
		return us, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return us, err
	}
	us.RowsAffected = int(count)
	// update reply.like_count || reply.dislike_count
	sqlStr = parseReactionCountSQL("reply_reaction", "reaction_id", "reply", "reply_id")
	_, err = c.db.Exec(sqlStr, reactionOnReply.ReplyID)
	if err != nil {
		return us, err
	}
	return us, nil
}
func reactionOnReplyDelete(reactionOnReply reactionOnReplyAPI) (us updateStatusAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return us, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		DELETE FROM reply_reaction 
		WHERE reply_id=$1 AND user_id=$2;
	`
	res, err := c.db.Exec(sqlStr,
		reactionOnReply.ReplyID, reactionOnReply.User.UserID)
	if err != nil {
		return us, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return us, err
	}
	us.RowsAffected = int(count)
	// update reply.like_count || reply.dislike_count
	sqlStr = parseReactionCountSQL("reply_reaction", "reaction_id", "reply", "reply_id")
	_, err = c.db.Exec(sqlStr, reactionOnReply.ReplyID)
	if err != nil {
		return us, err
	}
	return us, nil
}

// cronjob
func getAllUserID() (users []xuserAPI, err error) {
	c, err := connectPostgres()
	if err != nil {
		return users, errors.New("db connection")
	}
	defer c.db.Close()
	sqlStr := `
		SELECT 
		user_id 
		FROM xuser
		ORDER BY createtime ASC;
	`
	rows, err := c.db.Query(sqlStr)
	if err != nil {
		log.Println("getAllUserID", err)
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		user := xuserAPI{}
		if err := rows.Scan(
			&user.UserID,
		); err != nil {
			log.Println("err", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		log.Println("getAllUserID2", err)
		return users, err
	}
	return users, nil
}
func getPostsReadByUser(categoryID, weekTimestamp int, userID int64) (posts map[int64]int, err error) {
	c, err := connectMongoDB()
	if err != nil {
		log.Println("mongo session", err)
		return posts, err
	}
	defer c.session.Close()
	posts = make(map[int64]int)
	collection := c.session.DB(mongoDBXociety).C(mongoCollectionPostUserRead)
	for t := weekTimestamp; t > weekTimestamp-2*sevenDaysInSecond; t -= sevenDaysInSecond {
		u := postUserReadAPI{}
		q := bson.M{"user_id": userID, "category_id": categoryID, "week_timestamp": t}
		if err := collection.Find(q).One(&u); err == nil {
			for k, v := range u.PopularPosts { // to merge different week
				posts[k] = v
			}
		}
	}
	return posts, nil
}

func upsertPopularPostOnPostCommon(categoryID int, posts []postAPI) (err error) {
	c, err := connectMongoDB()
	if err != nil {
		log.Println("mongo session", err)
		return err
	}
	defer c.session.Close()
	collection := c.session.DB(mongoDBXociety).C(mongoCollectionPostCommon)
	selector := bson.M{"category_id": categoryID}
	if _, err := collection.Upsert(selector, bson.M{"$set": bson.M{"popular_posts": posts}}); err != nil {
		log.Println("upsertPopularPostOnPostCommon", err)
		return err
	}
	return nil
}
