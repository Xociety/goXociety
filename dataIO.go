package main

import (
	"database/sql"
	"log"
)

type country struct {
	CountryID int
	Name      string
	Code      string
}

type language struct {
	LanguageID      int
	DisplayLanguage string
	Value           string
}

type xuserDB struct {
	UserID     int64  `json:"user_id,omitempty"`
	Username   string `json:"username,omitempty"`
	Email      string `json:"email,omitempty"`
	Password   string `json:"password,omitempty"`
	Name       string `json:"name,omitempty"`
	Phone      string `json:"phone,omitempty"`
	Gender     int    `json:"gender,omitempty"`
	Bio        string `json:"bio,omitempty"`
	Credit     int    `json:"credit,omitempty"`
	LanguageID int    `json:"language_id,omitempty"`
	CountryID  int    `json:"country_id,omitempty"`
	Timezone   int    `json:"timezone,omitempty"`
	LastIP     string `json:"last_ip,omitempty"`
	Updatetime int    `json:"updatetime,omitempty"`
	Createtime int    `json:"createtime,omitempty"`
}

type postDB struct {
	PostID  int64  `json:"post_id,omitempty"`
	UserID  string `json:"user_id,omitempty"`
	Content string `json:"content,omitempty"`
	BlobID  string `json:"blob_id,omitempty"`
	// Point
	CountryID  string `json:"country_id,omitempty"`
	CategoryID string `json:"category_id,omitempty"`
	Public     bool   `json:"public,omitempty"`
	Createtime int    `json:"createtime,omitempty"`
	Updatetime int    `json:"updatetime,omitempty"`
}

type hashtagDB struct {
	HashtagID int64
	Name      string
}

type postHashtagDB struct {
	PostID    int64
	HashtagID int64
}

type postLikesDB struct {
	PostID     int64
	UserID     int64
	Type       int
	Createtime int
}

type commentsDB struct {
	CommentID  int64
	PostID     int64
	UserID     int64
	Comment    string
	Createtime int
	Updatetime int
}

type commentLikesDB struct {
	CommentID  int64
	UserID     int64
	Type       int
	Createtime int
}

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

func login(email, password string) (user xuserDB) {
	c := connectDB(postgresConStr, "PgSQL")
	defer c.db.Close()
	row := c.db.QueryRow("SELECT user_id, username, email, password, name, phone, gender, bio, credit, language_id, country_id, timezone, last_ip, createtime, updatetime FROM xuser WHERE email=$1 AND password=$2", email, password)
	if err := row.Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Phone,
		&user.Gender,
		&user.Bio,
		&user.Credit,
		&user.LanguageID,
		&user.CountryID,
		&user.Timezone,
		&user.LastIP,
		&user.Createtime,
		&user.Updatetime,
	); err != nil {
		log.Println("login", err)
	}
	return user
}

func getPost(categoryID int) (posts []postDB) {
	c := connectDB(postgresConStr, "PgSQL")
	defer c.db.Close()
	rows, err := c.db.Query("SELECT post_id, user_id, content, blob_id, country_id, category_id, public, createtime, updatetime FROM post WHERE category_id=$1;", categoryID)
	if err != nil {
		log.Println("getPost", err)
	}
	defer rows.Close()
	for rows.Next() {
		post := postDB{}
		if err := rows.Scan(
			&post.PostID,
			&post.UserID,
			&post.Content,
			&post.BlobID,
			&post.CountryID,
			&post.CategoryID,
			&post.Public,
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

func getXuserByID(userID string) (user xuserDB) {
	c := connectDB(postgresConStr, "PgSQL")
	defer c.db.Close()
	row := c.db.QueryRow("SELECT user_id, username, email, password, name, phone, gender, bio, credit, language_id, country_id, timezone, last_ip, createtime, updatetime FROM xuser WHERE user_id=$1;", userID)
	if err := row.Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Phone,
		&user.Gender,
		&user.Bio,
		&user.Credit,
		&user.LanguageID,
		&user.CountryID,
		&user.Timezone,
		&user.LastIP,
		&user.Createtime,
		&user.Updatetime,
	); err != nil {
		log.Println("getXuserByID", err)
	}
	return user
}

func getXuserByUsername(username string) (user xuserDB) {
	c := connectDB(postgresConStr, "PgSQL")
	defer c.db.Close()
	row := c.db.QueryRow("SELECT user_id, username, email, password, name, phone, gender, bio, credit, language_id, country_id, timezone, last_ip, createtime, updatetime FROM xuser WHERE username=$1;", username)
	if err := row.Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Phone,
		&user.Gender,
		&user.Bio,
		&user.Credit,
		&user.LanguageID,
		&user.CountryID,
		&user.Timezone,
		&user.LastIP,
		&user.Createtime,
		&user.Updatetime,
	); err != nil {
		log.Println("getXuserByUsername", err)
	}
	return user
}

func postNow(post postDB) (postID string) {
	// c := connectDB(postgresConStr, "PgSQL")
	// defer c.db.Close()
	// var userid int
	// err := c.db.QueryRow(`INSERT INTO users(name, favorite_fruit, age)
	// VALUES('beatrice', 'starfruit', 93) RETURNING id`).Scan(&userid)
	// log.Println(err)
	return postID
}
