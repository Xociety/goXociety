package main

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

type conn struct {
	db   *sql.DB
	name string
}

func connectDB(dbinfo string, name string) conn {
	var err error
	c := conn{name: name}
	c.db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		log.Println("connect", err)
	}
	return c
}

// auth
func checkSession(userToken string) (user xuserAPI, err error) {
	c := connectDB(postgresConStr, "PgSQL")
	defer c.db.Close()
	row := c.db.QueryRow(`
		SELECT 
		user_id, username, email, name, phone, 
		gender, bio, credit, photo_url, 
		language_id, country_id, 
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
		&user.CountryID,
		&user.Timezone,
		&user.LastIP,
		&user.Createtime,
		&user.Updatetime,
	); err != nil {
		// log.Println("checkSession", err)
		return user, err
	}
	return user, nil
}

// query
func login(email, password string) (lc loginAPI, err error) {
	c := connectDB(postgresConStr, "PgSQL")
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

func getXuserByID(userID string) (user xuserAPI, err error) {
	c := connectDB(postgresConStr, "PgSQL")
	defer c.db.Close()
	sqlStr := `
		SELECT 
		user_id, username, email, name, phone, 
		gender, bio, credit, photo_url, 
		language_id, country_id, 
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
		&user.CountryID,
		&user.Timezone,
		&user.LastIP,
		&user.Createtime,
		&user.Updatetime,
	); err != nil {
		// log.Println("getXuserByID", err)
		return user, errors.New("user not found")
	}
	return user, nil
}

func getXuserByUsername(username string) (user xuserAPI, err error) {
	c := connectDB(postgresConStr, "PgSQL")
	defer c.db.Close()
	sqlStr := `
		SELECT 
		user_id, username, email, name, phone, 
		gender, bio, credit, photo_url, 
		language_id, country_id, 
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
		&user.CountryID,
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

func getFollowingList(followerUserID string, page int) (users []xuserFollowingAPI, err error) {
	numPerRequest := 10
	c := connectDB(postgresConStr, "PgSQL")
	defer c.db.Close()
	sqlStr := `
		SELECT xuser.user_id, xuser.username, xuser.name, xuser.photo_url, follow.createtime
		FROM follow 
		FULL OUTER JOIN xuser ON follow.following_user_id = xuser.user_id
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
		user := xuserFollowingAPI{}
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

func getFollwerList(followerUserID string, page int) (users []xuserFollowerAPI, err error) {
	numPerRequest := 10
	c := connectDB(postgresConStr, "PgSQL")
	defer c.db.Close()
	sqlStr := `
		SELECT xuser.user_id, xuser.username, xuser.name, xuser.photo_url, follow.createtime
		FROM follow 
		FULL OUTER JOIN xuser ON follow.follower_user_id = xuser.user_id
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
		user := xuserFollowerAPI{}
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

func checkIfFollowing(followingUserID, followerUserID string) (isFollowing bool, err error) {
	c := connectDB(postgresConStr, "PgSQL")
	defer c.db.Close()
	count := 0
	sqlStr := `SELECT COUNT(*) FROM follow WHERE following_user_id=$1 AND follower_user_id=$2;`
	err = c.db.QueryRow(sqlStr, followingUserID, followerUserID).Scan(&count)
	if err != nil {
		return count == 1, err
	}
	return count == 1, err
}

func getPostsRecent(categoryID, page int) (posts []postAPI) {
	numPerRequest := 10
	timestamp := int(time.Now().Unix()) - twoMonthsInSecond
	c := connectDB(postgresConStr, "PgSQL")
	defer c.db.Close()
	sqlStr := `
		SELECT 
		post_id, user_id, content, blob_id, type, 
		like_count, dislike_count, country_id, 
		category_id, createtime, updatetime 
		FROM post 
		WHERE category_id=$1 AND createtime>=$2 
		ORDER BY createtime DESC OFFSET $3 LIMIT $4;
	`
	rows, err := c.db.Query(sqlStr, categoryID, timestamp, page*numPerRequest, numPerRequest)
	if err != nil {
		log.Println("getPostsRecent", err)
	}
	defer rows.Close()
	for rows.Next() {
		post := postAPI{}
		if err := rows.Scan(
			&post.PostID,
			&post.UserID,
			&post.Content,
			&post.BlobID,
			&post.Type,
			&post.LikeCount,
			&post.DislikeCount,
			&post.CountryID,
			&post.CategoryID,
			&post.Createtime,
			&post.Updatetime,
		); err != nil {
			log.Println("errrr", err)
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
	}
	return posts
}

func getFollowingUsersPosts(userID string, page int) (posts []postAPI) { // not done yet
	numPerRequest := 10
	c := connectDB(postgresConStr, "PgSQL")
	defer c.db.Close()
	sqlStr := `
		SELECT 
		post.post_id, post.user_id, post.content, post.blob_id, post.type, 
		post.like_count, post.dislike_count, post.country_id, 
		post.category_id, post.createtime, post.updatetime 
		FROM post 
		FULL OUTER JOIN follow ON follow.following_user_id = post.user_id
		WHERE follow.follower_user_id= $1
		ORDER BY post.createtime
		DESC OFFSET $2 LIMIT $3;
	`
	rows, err := c.db.Query(sqlStr, userID, page*numPerRequest, numPerRequest)
	if err != nil {
		log.Println("getFollowingUsersPosts", err)
	}
	defer rows.Close()
	for rows.Next() {
		post := postAPI{}
		if err := rows.Scan(
			&post.PostID,
			&post.UserID,
			&post.Content,
			&post.BlobID,
			&post.Type,
			&post.LikeCount,
			&post.DislikeCount,
			&post.CountryID,
			&post.CategoryID,
			&post.Createtime,
			&post.Updatetime,
		); err != nil {
			log.Println(err)
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
	}
	// log.Println("posts", posts)
	return posts
}

// mutation
func follow(followingUserID, followerUserID string) (user xuserFollowingAPI, err error) {
	c := connectDB(postgresConStr, "PgSQL")
	defer c.db.Close()
	timestamp := getNowUnixTimestamp()
	sqlStr := `
		INSERT INTO follow 
		(following_user_id, follower_user_id, valid, createtime, updatetime) 
		values($1,$2,$3,$4,$5) returning createtime;
	`
	err = c.db.QueryRow(sqlStr, followingUserID, followerUserID, true, timestamp, timestamp).Scan(&user.FollowingTime)
	if err != nil {
		return user, errors.New("you've followed this user")
	}
	return user, nil
}

func unfollow(followingUserID, followerUserID string) (ds deleteStatusAPI, err error) {
	c := connectDB(postgresConStr, "PgSQL")
	defer c.db.Close()
	sqlStr := `
		DELETE FROM follow 
		WHERE following_user_id=$1 AND follower_user_id=$2;
	`
	res, err := c.db.Exec(sqlStr,
		followingUserID, followerUserID)
	if err != nil {
		return ds, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return ds, err
	}
	ds.RowsAffected = int(count)
	return ds, nil
}

func postNow(post postDB) (postID string) {
	c := connectDB(postgresConStr, "PgSQL")
	defer c.db.Close()
	sqlStr := `
		INSERT INTO post 
		(
			user_id, content, blob_id, type, 
			like_count, dislike_count,
			country_id, category_id, public, createtime, updatetime
		) VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING post_id;
	`
	err := c.db.QueryRow(sqlStr,
		post.UserID,
		post.Content,
		post.BlobID,
		post.Type,
		0,
		0,
		post.CountryID,
		post.CategoryID,
		post.Public,
		post.Createtime,
		post.Updatetime,
	).Scan(&postID)
	// tag
	// hashtag
	// post_hashtag
	if err != nil {
		log.Println(err)
	}
	return postID
}
