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

func getUserByUserID(userID int64) (user xuserAPI, err error) {
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

func getUserByUsername(username string) (user xuserAPI, err error) {
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

func getUsersByFollowing(followerUserID int64, page int) (users []userFollowingAPI, err error) {
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

func makeBlobURL(post postAPI) string {
	// restore url
	url := "http://" + bucketRootCloudStorage + "/" + makeBucketFolderName(post.Type, post.Blob.BlobID)
	switch postTypeMapID2Type[post.Type] {
	case mediaFormatJPG:
		url += "0." + mediaFormatJPG
	case mediaFormatHLS:
		url += "0." + mediaFormatM3U8
	}
	return url
}
func getPostsByRecent(categoryID, page int) (posts []postAPI) {
	numPerRequest := 10
	timestamp := int(time.Now().Unix()) - twoMonthsInSecond // sixHoursInSecond
	c := connectDB(postgresConStr, "PgSQL")
	defer c.db.Close()
	sqlStr := `
		SELECT 
		post.post_id,
		post.user_id, xuser.username, xuser.name,
		post.content, post.blob_id, post.origin_width, post.origin_height, post.type, 
		post.like_count, post.dislike_count, post.comment_count, post.country_id, 
		post.category_id, post.createtime, post.updatetime 
		FROM post 
		FULL OUTER JOIN xuser ON xuser.user_id = post.user_id
		WHERE post.category_id=$1 AND post.createtime>=$2 
		ORDER BY post.createtime DESC OFFSET $3 LIMIT $4;
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
		}
		post.Blob.BlobID = makeBlobURL(post)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
	}
	return posts
}
func getPostsByFollowingUsers(userID int64, page int) (posts []postAPI) { // not done yet
	numPerRequest := 10
	c := connectDB(postgresConStr, "PgSQL")
	defer c.db.Close()
	sqlStr := `
		SELECT 
		post.post_id, 
		post.user_id, xuser.username, xuser.name,
		post.content, post.blob_id, post.origin_width, post.origin_height, post.type,
		post.like_count, post.dislike_count, post.comment_count,
		post.country_id, post.category_id, post.createtime, post.updatetime 
		FROM post
		FULL OUTER JOIN xuser ON xuser.user_id = post.user_id
		FULL OUTER JOIN follow ON follow.following_user_id = post.user_id
		WHERE follow.follower_user_id= $1
		ORDER BY post.createtime
		DESC OFFSET $2 LIMIT $3;
	`
	rows, err := c.db.Query(sqlStr, userID, page*numPerRequest, numPerRequest)
	if err != nil {
		log.Println("getPostsFollowing", err)
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
		}
		post.Blob.BlobID = makeBlobURL(post)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
	}
	// log.Println("posts", posts)
	return posts
}
func getPostsByUser(userID int64, page int) (posts []postAPI) { // not done yet
	numPerRequest := 10
	c := connectDB(postgresConStr, "PgSQL")
	defer c.db.Close()
	sqlStr := `
		SELECT 
		post.post_id, 
		post.user_id, xuser.username, xuser.name,
		post.content, post.blob_id, post.origin_width, post.origin_height, post.type,
		post.like_count, post.dislike_count, post.comment_count,
		post.country_id, post.category_id, post.createtime, post.updatetime 
		FROM post
		FULL OUTER JOIN xuser ON xuser.user_id = post.user_id
		WHERE post.user_id= $1
		ORDER BY post.createtime
		DESC OFFSET $2 LIMIT $3;
	`
	rows, err := c.db.Query(sqlStr, userID, page*numPerRequest, numPerRequest)
	if err != nil {
		log.Println("getPostsUser", err)
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
		}
		post.Blob.BlobID = makeBlobURL(post)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
	}
	// log.Println("posts", posts)
	return posts
}

func getCommentsByPost(postID int64, page int) (comments []commentAPI, err error) {
	numPerRequest := 10
	c := connectDB(postgresConStr, "PgSQL")
	defer c.db.Close()
	sqlStr := `
		SELECT 
		comment.comment_id, 
		comment.user_id, xuser.username, xuser.name, 
		comment.comment, comment.createtime, comment.updatetime 
		FROM comment FULL OUTER JOIN xuser ON comment.user_id = xuser.user_id 
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
		if err := rows.Scan(
			&comment.CommentID,
			&comment.User.UserID,
			&comment.User.Username,
			&comment.User.Name,
			&comment.Comment,
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

func getReactionsByPost(postID int64, page int) (reactionsOnPost []reactionOnPostAPI, err error) {
	numPerRequest := 10
	c := connectDB(postgresConStr, "PgSQL")
	defer c.db.Close()
	sqlStr := `
		SELECT 
		post_reaction.post_id, 
		post_reaction.user_id, xuser.username, xuser.name, 
		post_reaction.reaction_id, post_reaction.createtime 
		FROM public.post_reaction FULL OUTER JOIN xuser on post_reaction.user_id = xuser.user_id
		WHERE post_id=$1
		ORDER BY createtime DESC OFFSET $2 LIMIT $3;
	`
	rows, err := c.db.Query(sqlStr, postID, page*numPerRequest, numPerRequest)
	if err != nil {
		log.Println("getPostReactions", err)
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
func getReactionsByComment(commentID int64, page int) (reactionsOnComment []reactionOnCommentAPI, err error) {
	numPerRequest := 10
	c := connectDB(postgresConStr, "PgSQL")
	defer c.db.Close()
	sqlStr := `
		SELECT 
		comment_reaction.comment_id, 
		comment_reaction.user_id, xuser.username, xuser.name, 
		comment_reaction.reaction_id, comment_reaction.createtime 
		FROM public.comment_reaction FULL OUTER JOIN xuser on comment_reaction.user_id = xuser.user_id
		WHERE comment_id=$1
		ORDER BY createtime DESC OFFSET $2 LIMIT $3;
	`
	rows, err := c.db.Query(sqlStr, commentID, page*numPerRequest, numPerRequest)
	if err != nil {
		log.Println("getReactionsByComment", err)
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

// mutation
func follow(followingUserID, followerUserID int64) (us updateStatusAPI, err error) {
	c := connectDB(postgresConStr, "PgSQL")
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
	c := connectDB(postgresConStr, "PgSQL")
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

func postInsert(post postAPI) (postID int64, err error) {
	c := connectDB(postgresConStr, "PgSQL")
	defer c.db.Close()
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
	// tag
	// hashtag
	// post_hashtag
	if err != nil {
		// log.Println(err)
		return postID, err
	}
	return postID, err
}
func postUpdate(post postAPI) (us updateStatusAPI, err error) {
	c := connectDB(postgresConStr, "PgSQL")
	defer c.db.Close()
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
	c := connectDB(postgresConStr, "PgSQL")
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

func commentOnPostInsert(comment commentAPI) (commentID int64, err error) {
	c := connectDB(postgresConStr, "PgSQL")
	defer c.db.Close()
	sqlStr := `
		INSERT INTO comment 
		(
			post_id, user_id, comment, 
			like_count, dislike_count, comment_count,
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
		comment.CommentCount,
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
	c := connectDB(postgresConStr, "PgSQL")
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
	c := connectDB(postgresConStr, "PgSQL")
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

func reactionOnPostSet(reactionOnPost reactionOnPostAPI) (us updateStatusAPI, err error) {
	c := connectDB(postgresConStr, "PgSQL")
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
	sqlStr = `
		UPDATE post
		SET ` + reactionsTypeMapID2Description[reactionOnPost.ReactionID] + `_count =
		(SELECT COUNT(*) FROM post_reaction
		WHERE post_reaction.post_id = post.post_id AND post_reaction.post_id = $1) WHERE post_id = $1;
	`
	_, err = c.db.Exec(sqlStr, reactionOnPost.PostID)
	if err != nil {
		return us, err
	}
	return us, nil
}
func reactionOnPostDelete(reactionOnPost reactionOnPostAPI) (us updateStatusAPI, err error) {
	c := connectDB(postgresConStr, "PgSQL")
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
	sqlStr = `
		UPDATE post
		SET ` + reactionsTypeMapID2Description[reactionOnPost.ReactionID] + `_count =
		(SELECT COUNT(*) FROM post_reaction
		WHERE post_reaction.post_id = post.post_id AND post_reaction.post_id = $1) WHERE post_id = $1;
	`
	_, err = c.db.Exec(sqlStr, reactionOnPost.PostID)
	if err != nil {
		return us, err
	}
	return us, nil
}
func reactionOnCommentSet(reactionOnComment reactionOnCommentAPI) (us updateStatusAPI, err error) {
	c := connectDB(postgresConStr, "PgSQL")
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
	sqlStr = `
		UPDATE comment
		SET ` + reactionsTypeMapID2Description[reactionOnComment.ReactionID] + `_count =
		(SELECT COUNT(*) FROM comment_reaction
		WHERE comment_reaction.comment_id = comment.comment_id AND comment_reaction.comment_id = $1) WHERE comment_id = $1;
	`
	_, err = c.db.Exec(sqlStr, reactionOnComment.CommentID)
	if err != nil {
		return us, err
	}
	return us, nil
}
func reactionOnCommentDelete(reactionOnComment reactionOnCommentAPI) (us updateStatusAPI, err error) {
	c := connectDB(postgresConStr, "PgSQL")
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
	sqlStr = `
		UPDATE comment
		SET ` + reactionsTypeMapID2Description[reactionOnComment.ReactionID] + `_count =
		(SELECT COUNT(*) FROM comment_reaction
		WHERE comment_reaction.comment_id = comment.comment_id AND comment_reaction.comment_id = $1) WHERE comment_id = $1;
	`
	_, err = c.db.Exec(sqlStr, reactionOnComment.CommentID)
	if err != nil {
		return us, err
	}
	return us, nil
}
