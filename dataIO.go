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

type xuser struct {
	UserID     int64
	UserName   string
	Email      string
	Password   string
	Name       string
	Phone      string
	Gender     int
	Bio        string
	Credit     int
	Language   int
	Country    int
	Timezone   int
	LastIP     string
	Updatetime int
	Createtime int
}

type media struct {
	MediaID int64
	UserID  string
	Content string
	BlobID  string
	// Point
	CountryID  string
	CategoryID string
	Createtime int
	Updatetime int
}

type hashtag struct {
	HashtagID int64
	Name      string
}

type mediaHashtag struct {
	MediaID   int64
	HashtagID int64
}

type mediaLikes struct {
	MediaID    int64
	UserID     int64
	Type       int
	Createtime int
}

type comments struct {
	CommentID  int64
	MediaID    int64
	UserID     int64
	Comment    string
	Createtime int
	Updatetime int
}

type commentLikes struct {
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

func login(email, password string) (userID int64) {
	// c := connectDB(postgresConStr, "PgSQL")
	// defer c.db.Close()
	db, err := sql.Open("postgres", postgresConStr)
	if err != nil {
		log.Println("db", err)
	}
	defer db.Close()
	log.Println(email, password)
	// row := db.QueryRow("select user_id from xuser where email='kyler@gmail.com' and password='salted'")
	row := db.QueryRow("select user_id from xuser where email=$1 and password=$2", email, password)
	if err := row.Scan(&userID); err != nil {
		log.Println("login", err)
	}
	log.Println(userID)
	return userID
}
