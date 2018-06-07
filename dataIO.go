package main

import (
	"database/sql"
	"log"
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

// query
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

// mutation
func postNow(post postDB) (postID string) {
	c := connectDB(postgresConStr, "PgSQL")
	defer c.db.Close()
	err := c.db.QueryRow(`
		INSERT INTO post 
		(user_id, content, blob_id, country_id, category_id, public, type, createtime, updatetime) 
		VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING post_id;`,
		post.UserID,
		post.Content,
		post.BlobID,
		post.CountryID,
		post.CategoryID,
		post.Public,
		post.Type,
		post.Createtime,
		post.Updatetime,
	).Scan(&postID)
	log.Println(err)
	return postID
}
