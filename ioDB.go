package main

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	_ "github.com/lib/pq"
)

// net
type connPostgres struct {
	db *sql.DB
}

type connMongo struct {
	session *mgo.Session
}

func connectPostgres(dbinfo string) (connPostgres, error) {
	var err error
	c := connPostgres{}
	c.db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		log.Println("postgres connect", err)
		return c, err
	}
	return c, nil
}
func connectMongoDB(dbinfo string) (connMongo, error) {
	var err error
	c := connMongo{}
	c.session, err = mgo.Dial(dbinfo)
	if err != nil {
		log.Println("mongo connect", err)
		return c, err
	}
	return c, nil
}

// auth
func checkSession(userToken string) (user xuserAPI, err error) {
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return user, errors.New("db connection")
	}
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

// [query]
func login(email, password string) (lc loginAPI, err error) {
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return lc, errors.New("db connection")
	}
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

// user
func getUserByUserID(userID int64) (user xuserAPI, err error) {
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return user, errors.New("db connection")
	}
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
func getUserByUsername(username string) (user xuserAPI, err error) {
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return user, errors.New("db connection")
	}
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

// follow
func getUsersByFollowing(followerUserID int64, page int) (users []userFollowingAPI, err error) {
	numPerRequest := 10
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return users, errors.New("db connection")
	}
	sqlStr := `
		SELECT xuser.user_id, xuser.username, xuser.name, xuser.photo_url, follow.createtime
		FROM follow 
		INNER JOIN xuser ON follow.following_user_id = xuser.user_id
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
	numPerRequest := 10
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return users, errors.New("db connection")
	}
	sqlStr := `
		SELECT xuser.user_id, xuser.username, xuser.name, xuser.photo_url, follow.createtime
		FROM follow 
		INNER JOIN xuser ON follow.follower_user_id = xuser.user_id
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
	count := 0
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return isFollowing, errors.New("db connection")
	}
	sqlStr := `SELECT COUNT(*) FROM follow WHERE following_user_id=$1 AND follower_user_id=$2;`
	err = c.db.QueryRow(sqlStr, followingUserID, followerUserID).Scan(&count)
	if err != nil {
		return count == 1, err
	}
	return count == 1, err
}

// post
func getPostsByRecentPage(categoryID, page int) (posts []postAPI, err error) {
	numPerRequest := 10
	timestamp := int(time.Now().Unix()) - twoMonthsInSecond // sixHoursInSecond
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return posts, errors.New("db connection")
	}
	sqlStr := `
		SELECT 
		post.post_id,
		post.user_id, xuser.username, xuser.name,
		post.content, post.blob_id, post.origin_width, post.origin_height, post.type, 
		post.like_count, post.dislike_count, post.comment_count, post.country_id, 
		post.category_id, post.createtime, post.updatetime 
		FROM post 
		INNER JOIN xuser ON xuser.user_id = post.user_id
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
		if err := rows.Scan(
			&post.PostID,
			&post.User.UserID,
			&post.User.Username,
			&post.User.Name,
			&post.Content,
			&post.Blob.BlobID,
			&post.Blob.OriginWidth,
			&post.Blob.OriginHeight,
			&post.Type,
			&post.LikeCount,
			&post.DislikeCount,
			&post.CommentCount,
			&post.CountryID,
			&post.CategoryID,
			&post.Createtime,
			&post.Updatetime,
		); err != nil {
			log.Println("errrr", err)
			return posts, err
		}
		post.Blob.BlobID = makeBlobURL(post)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return posts, err
	}
	return posts, err
}
func getPostsByFollowingUsers(userID int64, page int) (posts []postAPI, err error) { // not done yet
	numPerRequest := 10
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return posts, errors.New("db connection")
	}
	sqlStr := `
		SELECT 
		post.post_id, 
		post.user_id, xuser.username, xuser.name,
		post.content, post.blob_id, post.origin_width, post.origin_height, post.type,
		post.like_count, post.dislike_count, post.comment_count,
		post.country_id, post.category_id, post.createtime, post.updatetime 
		FROM post
		INNER JOIN xuser ON xuser.user_id = post.user_id
		INNER JOIN follow ON follow.following_user_id = post.user_id
		WHERE follow.follower_user_id= $1
		ORDER BY post.createtime
		DESC OFFSET $2 LIMIT $3;
	`
	rows, err := c.db.Query(sqlStr, userID, page*numPerRequest, numPerRequest)
	if err != nil {
		log.Println("getPostsFollowing", err)
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
			&post.Content,
			&post.Blob.BlobID,
			&post.Blob.OriginWidth,
			&post.Blob.OriginHeight,
			&post.Type,
			&post.LikeCount,
			&post.DislikeCount,
			&post.CommentCount,
			&post.CountryID,
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
func getPostsByUser(userID int64, page int) (posts []postAPI, err error) { // not done yet
	numPerRequest := 10
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return posts, errors.New("db connection")
	}
	sqlStr := `
		SELECT 
		post.post_id, 
		post.user_id, xuser.username, xuser.name,
		post.content, post.blob_id, post.origin_width, post.origin_height, post.type,
		post.like_count, post.dislike_count, post.comment_count,
		post.country_id, post.category_id, post.createtime, post.updatetime 
		FROM post
		INNER JOIN xuser ON xuser.user_id = post.user_id
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
		if err := rows.Scan(
			&post.PostID,
			&post.User.UserID,
			&post.User.Username,
			&post.User.Name,
			&post.Content,
			&post.Blob.BlobID,
			&post.Blob.OriginWidth,
			&post.Blob.OriginHeight,
			&post.Type,
			&post.LikeCount,
			&post.DislikeCount,
			&post.CommentCount,
			&post.CountryID,
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
func getPostsByPopular(userID int64, categoryID, page int) (posts []postAPI, err error) { // not done yet
	numPerRequest := 10
	c, err := connectMongoDB(globalConfig[env].MongoConStr)
	defer c.session.Close()
	if err != nil {
		log.Println("mongo session", err)
		return posts, err
	}
	collection := c.session.DB(mongoDBXociety).C(mongoCollectionPostPopular)
	u := userPostPopular{}
	if err := collection.Find(bson.M{"user_id": userID, "category_id": categoryID}).One(&u); err != nil {
		log.Println("getPostsByPopular", err)
		return posts, err
	}
	// fake page because there's no appropriate query in mongo
	for i := 0; i < len(u.Posts); i++ {
		if i >= numPerRequest {
			break
		}
		posts = append(posts, u.Posts[i])
	}
	return posts, nil
}
func getPostByPostIDUserID(postID, userID int64) (count int, err error) {
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return count, errors.New("db connection")
	}
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

// hashtags
func getHashtags(value string, page int) (hashtags []hashtagAPI, err error) {
	numPerRequest := 10
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return hashtags, errors.New("db connection")
	}
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
func getPostsByHashtag(hashtagID int64, page int) (posts []postAPI, err error) { // not done yet
	numPerRequest := 10
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return posts, errors.New("db connection")
	}
	sqlStr := `
		SELECT 
		post.post_id, 
		post.user_id, xuser.username, xuser.name,
		post.content, post.blob_id, post.origin_width, post.origin_height, post.type,
		post.like_count, post.dislike_count, post.comment_count,
		post.country_id, post.category_id, post.createtime, post.updatetime 
		FROM post
		INNER JOIN xuser ON xuser.user_id = post.user_id
		INNER JOIN post_hashtag ON post_hashtag.post_id = post.post_id
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
			&post.Content,
			&post.Blob.BlobID,
			&post.Blob.OriginWidth,
			&post.Blob.OriginHeight,
			&post.Type,
			&post.LikeCount,
			&post.DislikeCount,
			&post.CommentCount,
			&post.CountryID,
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

// tags
func getPostsByTag(userID int64, page int) (posts []postAPI, err error) {
	numPerRequest := 10
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return posts, errors.New("db connection")
	}
	sqlStr := `
		SELECT 
		post.post_id, 
		post.user_id, xuser.username, xuser.name,
		post.content, post.blob_id, post.origin_width, post.origin_height, post.type,
		post.like_count, post.dislike_count, post.comment_count,
		post.country_id, post.category_id, post.createtime, post.updatetime 
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
			&post.Content,
			&post.Blob.BlobID,
			&post.Blob.OriginWidth,
			&post.Blob.OriginHeight,
			&post.Type,
			&post.LikeCount,
			&post.DislikeCount,
			&post.CommentCount,
			&post.CountryID,
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
func getAllTagsByPost(postID int64) (tags []tagOnPostAPI, err error) {
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return tags, errors.New("db connection")
	}
	sqlStr := `
		SELECT 
		post_tag_xuser.post_id,
		post_tag_xuser.user_id, xuser.username, xuser.name,
		post_tag_xuser.x, post_tag_xuser.y,
		post_tag_xuser.valid,
		post_tag_xuser.createtime, post_tag_xuser.updatetime
		FROM post_tag_xuser
		INNER JOIN xuser ON xuser.user_id = post_tag_xuser.user_id
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
	numPerRequest := 10
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return comments, errors.New("db connection")
	}
	sqlStr := `
		SELECT 
		comment.comment_id, 
		comment.user_id, xuser.username, xuser.name, 
		comment.comment, 
		comment.like_count, comment.dislike_count, comment.reply_count,
		comment.createtime, comment.updatetime 
		FROM comment INNER JOIN xuser ON comment.user_id = xuser.user_id 
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
	numPerRequest := 10
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return replies, errors.New("db connection")
	}
	sqlStr := `
		SELECT 
		reply.reply_id, 
		reply.user_id, xuser.username, xuser.name, 
		reply.reply, 
		reply.like_count, reply.dislike_count,
		reply.createtime, reply.updatetime 
		FROM reply INNER JOIN xuser ON reply.user_id = xuser.user_id 
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
	numPerRequest := 10
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return reactionsOnPost, errors.New("db connection")
	}
	sqlStr := `
		SELECT 
		post_reaction.post_id, 
		post_reaction.user_id, xuser.username, xuser.name, 
		post_reaction.reaction_id, post_reaction.createtime 
		FROM public.post_reaction INNER JOIN xuser on post_reaction.user_id = xuser.user_id
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
	numPerRequest := 10
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return reactionsOnComment, errors.New("db connection")
	}
	sqlStr := `
		SELECT 
		comment_reaction.comment_id, 
		comment_reaction.user_id, xuser.username, xuser.name, 
		comment_reaction.reaction_id, comment_reaction.createtime 
		FROM public.comment_reaction INNER JOIN xuser on comment_reaction.user_id = xuser.user_id
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
	numPerRequest := 10
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return reactionsOnReply, errors.New("db connection")
	}
	sqlStr := `
		SELECT 
		reply_reaction.reply_id, 
		reply_reaction.user_id, xuser.username, xuser.name, 
		reply_reaction.reaction_id, reply_reaction.createtime 
		FROM public.reply_reaction INNER JOIN xuser on reply_reaction.user_id = xuser.user_id
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
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return userID, errors.New("db connection")
	}
	sqlStr := `
		INSERT INTO xuser 
		(
			username, email, password, name, phone, 
			gender, bio, credit, photo_url, 
			language_id, country_id, 
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
		user.CountryID,
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
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return us, errors.New("db connection")
	}
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
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return us, errors.New("db connection")
	}
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

// post
func postInsert(post postAPI) (postID int64, err error) {
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return postID, errors.New("db connection")
	}
	sqlStr := `
		INSERT INTO post 
		(
			user_id, content, blob_id, origin_width, origin_height, type, 
			like_count, dislike_count, comment_count,
			country_id, category_id, public, createtime, updatetime
		) VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING post_id;
	`
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
		post.CountryID,
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
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return us, errors.New("db connection")
	}
	sqlStr := `
		UPDATE post 
		SET content=$1, country_id=$2, category_id=$3, updatetime=$4
		WHERE post_id=$5 AND user_id=$6;
	`
	res, err := c.db.Exec(sqlStr,
		post.Content, post.CountryID, post.CategoryID, post.Updatetime,
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
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return us, errors.New("db connection")
	}
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
func postPopularRead(categoryID, indexRead int, userID int64) (posts []postAPI, err error) {
	// you can wrap func as a transaction
	numPerRequest := 10
	lastIndexNextList := indexRead + numPerRequest
	c, err := connectMongoDB(globalConfig[env].MongoConStr)
	defer c.session.Close()
	if err != nil {
		log.Println("mongo session", err)
		return posts, err
	}
	collection := c.session.DB(mongoDBXociety).C(mongoCollectionPostPopular)
	u := userPostPopular{}
	if err := collection.Find(bson.M{"user_id": userID, "category_id": categoryID}).One(&u); err != nil {
		log.Println("postPopularRead", err)
		return posts, err
	}
	weekTimestamp := getNowUnixWeekTimestamp()
	timestamp := getNowUnixTimestamp()
	postsPopular := []postAPI{} // popular order
	postsRead := make(map[int64]int)
	for i := 0; i < len(u.Posts); i++ {
		if i <= indexRead {
			postsRead[u.Posts[i].PostID] = timestamp // for record read post
		} else {
			// fake page because there's no appropriate query in mongo
			if i <= lastIndexNextList {
				posts = append(posts, u.Posts[i]) // for current query popular post list
			}
			postsPopular = append(postsPopular, u.Posts[i]) // for update post list in db
		}
	}
	collection2 := c.session.DB(mongoDBXociety).C(mongoCollectionPostRead) // you can make this part using another channel queue
	if _, err := collection2.Upsert(
		bson.M{"user_id": userID, "category_id": categoryID, "week_timestamp": weekTimestamp},
		bson.M{"$set": parsePopularPostReadObjectMongo(postsRead)}); err != nil {
		log.Println("upsert", err)
	}
	if _, err := collection.Upsert( // update popular post list
		bson.M{"user_id": userID, "category_id": categoryID},
		bson.M{"$set": bson.M{"posts": postsPopular}}); err != nil {
		log.Println("upsert", err)
	}
	return posts, nil
}

// hashtags
func hashtagInsert(hashtags []string) (hashtagsID []int64, err error) {
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return hashtagsID, errors.New("db connection")
	}
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
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return us, errors.New("db connection")
	}
	sqlStrInsert, sqlStrDelete, args := parsehashtagOnPostSQL(postID, hashtagsID)
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

// post_tags
func tagOnPostUpdate(postID int64, tag tagOnPostSetAPI) (us updateStatusAPI, err error) {
	timestamp := getNowUnixTimestamp()
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return us, errors.New("db connection")
	}
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
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return us, errors.New("db connection")
	}
	sqlStr, args := parseTagOnPostInserSQL(postID, tags)
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
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return us, errors.New("db connection")
	}
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
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return us, errors.New("db connection")
	}
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
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return commentID, errors.New("db connection")
	}
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
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return us, errors.New("db connection")
	}
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
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return us, errors.New("db connection")
	}
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
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return replyID, errors.New("db connection")
	}
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
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return us, errors.New("db connection")
	}
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
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return us, errors.New("db connection")
	}
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
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return us, errors.New("db connection")
	}
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
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return us, errors.New("db connection")
	}
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
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return us, errors.New("db connection")
	}
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
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return us, errors.New("db connection")
	}
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
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return us, errors.New("db connection")
	}
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
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return us, errors.New("db connection")
	}
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
func getPostsByRecentNum(categoryID, numPost int) (posts []postAPI, err error) {
	// combine this with func getPostsByRecentPage
	timestamp := int(time.Now().Unix()) - twoMonthsInSecond // sixHoursInSecond
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return posts, errors.New("db connection")
	}
	sqlStr := `
		SELECT 
		post.post_id,
		post.user_id, xuser.username, xuser.name,
		post.content, post.blob_id, post.origin_width, post.origin_height, post.type, 
		post.like_count, post.dislike_count, post.comment_count, post.country_id, 
		post.category_id, post.createtime, post.updatetime 
		FROM post 
		INNER JOIN xuser ON xuser.user_id = post.user_id
		WHERE post.category_id=$1 AND post.createtime>=$2 
		ORDER BY post.createtime DESC OFFSET $3 LIMIT $4;
	`
	rows, err := c.db.Query(sqlStr, categoryID, timestamp, 0, numPost)
	if err != nil {
		log.Println("getPostsRecent", err)
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
			&post.Content,
			&post.Blob.BlobID,
			&post.Blob.OriginWidth,
			&post.Blob.OriginHeight,
			&post.Type,
			&post.LikeCount,
			&post.DislikeCount,
			&post.CommentCount,
			&post.CountryID,
			&post.CategoryID,
			&post.Createtime,
			&post.Updatetime,
		); err != nil {
			log.Println("errrr", err)
			return posts, err
		}
		post.Blob.BlobID = makeBlobURL(post)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return posts, err
	}
	return posts, nil
}
func getAllUserID() (users []xuserAPI, err error) {
	c, err := connectPostgres(globalConfig[env].PostgresConStr)
	defer c.db.Close()
	if err != nil {
		return users, errors.New("db connection")
	}
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
	c, err := connectMongoDB(globalConfig[env].MongoConStr)
	defer c.session.Close()
	if err != nil {
		log.Println("mongo session", err)
		return posts, err
	}
	collection := c.session.DB(mongoDBXociety).C(mongoCollectionPostRead)
	for t := weekTimestamp; t > weekTimestamp-2*sevenDaysInSecond; t -= sevenDaysInSecond {
		u := userPostPopularRead{}
		if err := collection.Find(bson.M{"user_id": userID, "category_id": categoryID, "week_timestamp": t}).One(&u); err == nil {
			posts = u.Posts
		}
	}
	return posts, nil
}
func upsertPostPopular(categoryID int, userID int64, posts []postAPI) (err error) {
	c, err := connectMongoDB(globalConfig[env].MongoConStr)
	defer c.session.Close()
	if err != nil {
		log.Println("mongo session", err)
		return err
	}
	collection := c.session.DB(mongoDBXociety).C(mongoCollectionPostPopular)
	selector := bson.M{"user_id": userID, "category_id": categoryID}
	if _, err := collection.Upsert(selector, bson.M{"$set": bson.M{"posts": posts}}); err != nil {
		log.Println("upsertPostPopular", err)
		return err
	}
	return nil
}
