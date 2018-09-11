package io

import (
	"log"
	"testing"
)

func TestMongoConnection(t *testing.T) {
	log.Println("popular post")
	c, err := connectPostgres()
	if err != nil {
		log.Println(err)
	}
	defer c.DB.Close()
	cm, err := connectMongoDB()
	if err != nil {
		log.Println(err)
	}
	defer cm.session.Close()
	posts, err := getPostsByPopular(&c, &cm, 1, 0, 0)
	log.Println(len(posts), err)
}

func TestPostgresConnection(t *testing.T) {
	log.Println("search user")
	c, err := connectPostgres()
	if err != nil {
		log.Println(err)
	}
	defer c.DB.Close()
	user, err := getUserByUserID(&c, 1)
	log.Println(user, err)
}
