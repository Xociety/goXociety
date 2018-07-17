package main

import (
	"strconv"

	"github.com/globalsign/mgo/bson"
)

type postByPopular []postAPI

func (p postByPopular) Len() int      { return len(p) }
func (p postByPopular) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p postByPopular) Less(i, j int) bool {
	return p[i].LikeCount+p[i].DislikeCount+p[i].CommentCount > p[j].LikeCount+p[j].DislikeCount+p[j].CommentCount
}

func parsePopularPostReadObjectMongo(posts map[int64]int) bson.M {
	/* for update mutiple object in one submmit, to generate bson.M like:
	bson.M{"post.1": 1,"post.2": 2}
	*/
	m := make(bson.M)
	for k, v := range posts {
		m["posts."+strconv.FormatInt(k, 10)] = v
	}
	return m
}
func filterReadedPost(postsRead map[int64]int, posts []postAPI) (filteredPost []postAPI) {
	for i := 0; i < len(posts); i++ {
		if postsRead[posts[i].PostID] == 0 {
			filteredPost = append(filteredPost, posts[i])
		}
	}
	return filteredPost
}
