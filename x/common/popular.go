package common

import (
	"strconv"

	"github.com/chienfuchen32/goXociety/x/config"
	"github.com/globalsign/mgo/bson"
)

type PostByPopular []config.PostAPI

func (p PostByPopular) Len() int      { return len(p) }
func (p PostByPopular) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p PostByPopular) Less(i, j int) bool {
	return p[i].LikeCount+p[i].DislikeCount+p[i].CommentCount > p[j].LikeCount+p[j].DislikeCount+p[j].CommentCount
}

func ParsePopularPostReadObjectMongo(posts map[int64]int) bson.M {
	/* for update mutiple object in one submmit, to generate bson.M like:
	bson.M{"post.1": 1,"post.2": 2}
	*/
	m := make(bson.M)
	for k, v := range posts {
		m["popular_posts."+strconv.FormatInt(k, 10)] = v
	}
	return m
}
func FilterReadedPost(postsRead map[int64]int, posts []config.PostAPI) (filteredPost []config.PostAPI) {
	for i := 0; i < len(posts); i++ {
		if postsRead[posts[i].PostID] == 0 {
			filteredPost = append(filteredPost, posts[i])
		}
	}
	return filteredPost
}

func ParseHashKeyFieldCommonPopularPostUserReadIndex(categoryID int, userID int64) (hashKey, hashfield string) {
	hashKey = config.RedisHashCommonPopularPostUserReadIndex + ":" + strconv.Itoa(categoryID)
	hashfield = "" + strconv.FormatInt(userID, 10)
	return
}
func ParseHashKeyFieldCountryPopularPostUserReadIndex(countryCode string, categoryID int, userID int64) (hashKey, hashfield string) {
	hashKey = config.RedisHashCountryPopularPostUserReadIndex + ":" + countryCode + ":" + strconv.Itoa(categoryID)
	hashfield = "" + strconv.FormatInt(userID, 10)
	return
}
func ParseHashKeyFieldCityPopularPostUserReadIndex(cityID string, categoryID int, userID int64) (hashKey, hashfield string) {
	hashKey = config.RedisHashCityPopularPostUserReadIndex + ":" + cityID + ":" + strconv.Itoa(categoryID)
	hashfield = strconv.FormatInt(userID, 10)
	return
}
