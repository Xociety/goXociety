package main

import (
	"log"
	"testing"
)

func TestMongoConnection(t *testing.T) {
	log.Println("popular post")
	posts, err := getPostsByPopular(1, 0, 0)
	log.Println(len(posts), err)
}

func TestPostgresConnection(t *testing.T) {
	log.Println("search user")
	user, err := getUserByUserID(1)
	log.Println(user, err)
}
