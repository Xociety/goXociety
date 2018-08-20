package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"testing"

	"cloud.google.com/go/storage"
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
	postType := 0             // r.Intn(2)
	originWidth := 0          // plz check current sample size
	originHeight := 0         // plz check current sample size
	categoryID := r.Intn(16)  // plz check current category
	sampleID := r.Intn(7) + 1 // plz check current sample file
	if categoryID == 0 {
		sampleID = 1
	}
	switch postTypeMapID2Type[postType] {
	case mediaFormatJPG:
		originWidth = 1920
		originHeight = 1280
		if categoryID != 0 {
			originWidth = 1242
			originHeight = 2004
		}
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
			BlobID:       "sample/" + strconv.Itoa(categoryID) + "/" + strconv.Itoa(sampleID),
			OriginWidth:  originWidth,
			OriginHeight: originHeight,
		},
		Type:         postType,
		LikeCount:    0,
		DislikeCount: 0,
		CommentCount: 0,
		CountryID:    0,
		CategoryID:   categoryID,
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

func genPopularPostUpsert() {
	users, err := getAllUserID()
	if err != nil {
		log.Println("user", err)
	}
	usersID := []int64{}
	for i := 0; i < len(users); i++ {
		usersID = append(usersID, users[i].UserID)
	}
	for categoryID := range categoryMapID2Name {
		posts, err := getPostsByRecentNum(categoryID, numPopularPostPerRefresh)
		if err != nil {
			log.Println("post", err)
		}
		sort.Sort(postByPopular(posts))
		// old
		for i := 0; i < len(usersID); i++ {
			weekTimestamp := getNowUnixWeekTimestamp()
			if postsRead, err := getPostsReadByUser(categoryID, weekTimestamp, usersID[i]); err == nil {
				filteredPosts := filterReadedPost(postsRead, posts)
				if err := upsertPostPopular(categoryID, usersID[i], filteredPosts); err != nil {
					log.Println(err)
				}
			}
		}
		// new
		if err := upsertInitPostPopularIndex(categoryID, usersID); err != nil {
			log.Println(err)
		}
		// common
		if err := upsertPostPopularCommon(categoryID, posts); err != nil {
			log.Println(err)
		}
	}
}
func TestUploadImagesSample(t *testing.T) { // delete this in the future
	type sampleImageCategory struct {
		category map[string]map[string]int
	}
	if os.Getenv("test_fakedata") != "true" || os.Getenv("test_fakedata") == "" {
		t.Skip("skip TestUploadSampleMedia")
	}
	sampleCategoryAppendIndex := func(c sampleImageCategory, cn string, scn []string, isc []int) {
		if len(scn) != len(isc) {
			return
		}
		if c.category[cn] == nil {
			c.category[cn] = make(map[string]int)
		}
		for i := range scn {
			c.category[cn][scn[i]] = isc[i]
		}
	}
	categorySample := sampleImageCategory{category: make(map[string]map[string]int)}

	list := []string{"Fragrance", "Makeup"}
	index := []int{15, 8}
	sampleCategoryAppendIndex(categorySample, "Beauty", list, index)
	list = []string{"Art", "Car", "Design", "Fitness", "Food", "Party", "Pets", "Shopping", "Sports", "Tech", "Travel"}
	index = []int{13, 9, 12, 10, 6, 5, 4, 11, 3, 2, 1}
	sampleCategoryAppendIndex(categorySample, "Life", list, index)
	list = []string{"Formal", "Street"}
	index = []int{14, 7}
	sampleCategoryAppendIndex(categorySample, "Style", list, index)
	list = []string{"Sup"}
	index = []int{0}
	sampleCategoryAppendIndex(categorySample, "Sup", list, index)
	// log.Println(categorySample)

	client, err := storage.NewClient(context.Background(), clientOptionGoogleAPI)
	defer client.Close()
	if err != nil {
		log.Println("gcp sdk client", err)
	}
	for k0, v0 := range categorySample.category {
		for k1, v1 := range v0 {
			for s := 1; s <= 7; s++ {
				// local
				ir, err := ioReaderFromFile("./development/upload/Pics/" + k0 + "/" + k1 + "/" + k1 + "/" + strconv.Itoa(s) + ".jpg")
				if err != nil {
					log.Println("err", err)
					continue
				}
				// cloud storage
				foldername := bucketImagesCloudStorage + "/sample/" + strconv.Itoa(v1) + "/" + strconv.Itoa(s) + "/"
				filename := "0.jpg"
				log.Println(globalConfig[env].GCPBucketRootCloudStorage, foldername+filename)
				if err := writeAndMakePublicCloudStorageGCP(client, globalConfig[env].GCPBucketRootCloudStorage, foldername+filename, ir); err != nil {
					log.Println("upload failed: ", err)
					continue
				}
			}
		}
	}
}
func TestFakeDataSample(t *testing.T) {
	// createtime > 1533631289
	if os.Getenv("test_fakedata") != "true" || os.Getenv("test_fakedata") == "" {
		t.Skip("skip TestFakeData")
	}
	log.Println("start")
	totalNumUser := 30
	totalNumPost := 200
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
			followingUserID = usersID[r.Intn(totalNumUser)]
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
	log.Printf("start random total %d post", totalNumPost)
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
	genPopularPostUpsert()
}
func TestPopularPostUpsertSample(t *testing.T) {
	if os.Getenv("test_fakedata") != "true" || os.Getenv("test_fakedata") == "" {
		t.Skip("skip TestSortPostPopular")
	}
	genPopularPostUpsert()
}
