package main

import (
	"log"
	"math/rand"
	"os"
	"sort"
	"testing"

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
		ReplyCount:   0,
		Createtime:   timestamp,
		Updatetime:   timestamp,
	}
}

func TestFakeData(t *testing.T) {
	// createtime > 1533631289
	if os.Getenv("test_fakedata") != "true" || os.Getenv("test_fakedata") == "" {
		t.Skip("skip TestSortPostPopular")
	}
	log.Println("start")
	totalNumUser := 20
	totalNumPost := 100
	usersID := []int64{}
	fake, err := faker.New("en")
	if err != nil {
		log.Fatalln("faker", err)
	}
	r := rand.New(rand.NewSource(99))
	// random people
	log.Printf("start random %d people", totalNumUser)
	for i := 0; i < totalNumUser; i++ {
		user := genXuserFaker(fake, r)
		userID, err := userInsert(user)
		if err == nil {
			usersID = append(usersID, userID)
		}
		// log.Println("user", err)
	}
	// random follow
	log.Println("start random follow")
	for i := 0; i < totalNumUser; i++ {
		totalNumFollow := r.Int63n(int64(totalNumUser) - 1)
		followingUserIDCheck := make(map[int64]bool)
		count := int64(0)
		follwerUserID := usersID[i]
		followingUserID := usersID[i]
		for count < totalNumFollow {
			followingUserID = usersID[r.Int63n(int64(totalNumUser))]
			if followingUserID == follwerUserID || followingUserIDCheck[followingUserID] {
				continue
			}
			if _, err := follow(followingUserID, follwerUserID); err != nil {
				log.Println("follow", err)
				continue
			}
			followingUserIDCheck[followingUserID] = true
			count++
		}
	}
	// random post
	log.Printf("start random %d post per user", totalNumPost)
	for i := 0; i < totalNumPost; i++ {
		userID := usersID[r.Intn(len(usersID))]
		post := genPostFaker(userID, fake, r)
		postID, err := postInsert(post)
		if err != nil {
			log.Println("post err", err)
		}
		userIDsNewReaction := append([]int64{}, usersID...) // prevent same person
		totalNumPostReaction := r.Intn(len(userIDsNewReaction))
		for j := 0; j < totalNumPostReaction; j++ {
			index := r.Intn(len(userIDsNewReaction))
			reactionOnPost := genPostReactionFaker(postID, userIDsNewReaction[index], r)
			if _, err = reactionOnPostSet(reactionOnPost); err != nil {
				log.Println("reaction err", err)
			}
			userIDsNewReaction = append(userIDsNewReaction[:index], userIDsNewReaction[index+1:]...)
		}
		totalNumPostComment := r.Intn(len(usersID))
		for j := 0; j < totalNumPostComment; j++ {
			index := r.Intn(len(usersID))
			commentOnPost := genPostCommentFaker(postID, usersID[index], fake, r)
			if _, err = commentOnPostInsert(commentOnPost); err != nil {
				log.Println("comment err", err)
			}
		}
	}
	// gen popular post
	log.Println("start popular post")
	genPopularPostUpsert(usersID)
}
func genPopularPostUpsert(usersID []int64) {
	for categoryID := range categoryMapID2Name {
		posts, err := getPostsByRecentNum(categoryID, numPopularPostPerRefresh)
		if err != nil {
			log.Println("post", err)
		}
		sort.Sort(postByPopular(posts))
		for i := 0; i < len(usersID); i++ {
			weekTimestamp := getNowUnixWeekTimestamp()
			if postsRead, err := getPostsReadByUser(categoryID, weekTimestamp, usersID[i]); err == nil {
				filteredPosts := filterReadedPost(postsRead, posts)
				if err := upsertPostPopular(categoryID, usersID[i], filteredPosts); err != nil {
					log.Println(err)
				}
			}
		}
	}
}
func TestPopularPostUpsert(t *testing.T) {
	if os.Getenv("test_fakedata") != "true" || os.Getenv("test_fakedata") == "" {
		t.Skip("skip TestSortPostPopular")
	}
	users, err := getAllUserID()
	if err != nil {
		log.Println("user", err)
	}
	usersID := []int64{}
	for i := 0; i < len(users); i++ {
		usersID = append(usersID, users[i].UserID)
	}
	genPopularPostUpsert(usersID)
}
