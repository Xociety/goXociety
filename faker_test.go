package main

import (
	"log"
	"math/rand"
	"sort"
	"testing"
	"time"

	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/manveru/faker"

	_ "github.com/lib/pq"
)

func genXuserFaker(fake *faker.Faker, r *rand.Rand) userDB {
	timeNow := getNowUnixTimestamp()
	return userDB{
		Username:   fake.UserName(),
		Email:      fake.Email(),
		Password:   "salted",
		Name:       fake.Name(),
		Phone:      fake.PhoneNumber(),
		Gender:     r.Intn(3),
		Bio:        fake.Sentence(r.Intn(3), true),
		Credit:     0,
		PhotoURL:   "",
		LanguageID: 13,
		CountryID:  206,
		Timezone:   28800,
		LastIP:     fake.IPv4Address().String(),
		Updatetime: timeNow,
		Createtime: timeNow,
	}
}

func genPostFaker(userID int64, faker *faker.Faker, r *rand.Rand) postAPI {
	timeNow := getNowUnixTimestamp()
	postType := r.Intn(2)
	originWidth := 0
	originHeight := 0
	switch postTypeMapID2Type[postType] {
	case mediaFormatJPG:
		originWidth = 1920
		originHeight = 1280
	case mediaFormatHLS:
		originWidth = 1920
		originHeight = 1080
	}
	return postAPI{
		User: userBasicAPI{
			UserID: userID,
		},
		Content: "test post",
		Blob: blobAPI{
			BlobID:       "sample",
			OriginWidth:  originWidth,
			OriginHeight: originHeight,
		},
		Type:         postType,
		LikeCount:    0,
		DislikeCount: 0,
		CommentCount: 0,
		CountryID:    0,
		CategoryID:   0,
		Public:       true,
		Createtime:   timeNow,
		Updatetime:   timeNow,
	}
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
		CommentCount: 0,
		Createtime:   timestamp,
		Updatetime:   timestamp,
	}
}

func TestStartInsertXuserfaker(t *testing.T) {
	// createtime > 1531187418
	// totalNumUser := 20
	// totalNumPost := 100
	// userIDsNew := []int64{}
	// fake, err := faker.New("en")
	// if err != nil {
	// 	log.Fatalln("faker", err)
	// }
	// r := rand.New(rand.NewSource(99))
	// for i := 0; i < totalNumUser; i++ {
	// 	user := genXuserFaker(fake, r)
	// 	userID, err := userInsert(user)
	// 	if err == nil {
	// 		userIDsNew = append(userIDsNew, userID)
	// 	}
	// 	// log.Println("user", err)
	// }
	// for i := 0; i < totalNumPost; i++ {
	// 	userID := userIDsNew[r.Intn(len(userIDsNew))]
	// 	post := genPostFaker(userID, fake, r)
	// 	postID, err := postInsert(post)
	// 	if err != nil {
	// 		log.Println("post err", err)
	// 	}
	// 	userIDsNewReaction := append([]int64{}, userIDsNew...)
	// 	totalNumPostReaction := r.Intn(len(userIDsNewReaction))
	// 	for j := 0; j < totalNumPostReaction; j++ {
	// 		index := r.Intn(len(userIDsNewReaction))
	// 		reactionOnPost := genPostReactionFaker(postID, userIDsNewReaction[index], r)
	// 		if _, err = reactionOnPostSet(reactionOnPost); err != nil {
	// 			log.Println("reaction err", err)
	// 		}
	// 		userIDsNewReaction = append(userIDsNewReaction[:index], userIDsNewReaction[index+1:]...)
	// 	}
	// 	userIDsNewComment := append([]int64{}, userIDsNew...)
	// 	totalNumPostComment := r.Intn(len(userIDsNewComment))
	// 	for j := 0; j < totalNumPostComment; j++ {
	// 		index := r.Intn(len(userIDsNewComment))
	// 		commentOnPost := genPostCommentFaker(postID, userIDsNewComment[index], fake, r)
	// 		if _, err = commentOnPostInsert(commentOnPost); err != nil {
	// 			log.Println("comment err", err)
	// 		}
	// 		userIDsNewComment = append(userIDsNewComment[:index], userIDsNewComment[index+1:]...)
	// 	}
	// }
}

type PostByPopular []postAPI

func (p PostByPopular) Len() int      { return len(p) }
func (p PostByPopular) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p PostByPopular) Less(i, j int) bool {
	return p[i].LikeCount+p[i].DislikeCount+p[i].CommentCount > p[j].LikeCount+p[j].DislikeCount+p[j].CommentCount
}
func TestSortPostPopular(t *testing.T) {
	posts := getPostsByRecentNum(0, 15)
	sort.Sort(PostByPopular(posts))
	session, err := mgo.Dial(mongoConStr)
	if err != nil {
		log.Println("session", err)
		return
	}
	c := session.DB(mongoDBXociety).C(mongoCollectionPostPopular)
	for i := 0; i < len(posts); i++ {
		data, err := bson.Marshal(posts[i])
		if err != nil {
			log.Println("data", err)
		}
		update := bson.M{}
		err = bson.Unmarshal(data, update)
		if err != nil {
			log.Println("update", err)
		}
		selector := bson.M{"post_id": posts[i].PostID}
		if _, err := c.Upsert(selector, bson.M{"$set": update}); err != nil {
			log.Println("upsert", err)
		}
	}
}

func TestQueryPostPopular(t *testing.T) {
	session, err := mgo.Dial(mongoConStr)
	defer session.Close()
	if err != nil {
		log.Printf("mongo connection lost ..\n")
		// return err
	}
	session.SetMode(mgo.Monotonic, true)
	session.SetSocketTimeout(3 * time.Minute)
	c := session.DB(mongoDBXociety).C(mongoCollectionPostPopular)
	posts := []postAPI{}
	if err := c.Find(bson.M{}).All(&posts); err != nil {
		log.Println(err)
		// return err
	}
	log.Println(posts)
	// return err
}
