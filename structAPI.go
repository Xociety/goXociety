package main

import "github.com/globalsign/mgo/bson"

// common
type updateStatusAPI struct {
	RowsAffected int `json:"rows_affected"`
}

// login
type loginAPI struct {
	Token string `json:"token"`
}

// user
type xuserAPI struct {
	UserID     int64  `json:"user_id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Gender     int    `json:"gender"`
	Bio        string `json:"bio"`
	Credit     int    `json:"credit"`
	PhotoURL   string `json:"photo_url"`
	LanguageID int    `json:"language_id"`
	CountryID  int    `json:"country_id"`
	Timezone   int    `json:"timezone"`
	LastIP     string `json:"last_ip"`
	Updatetime int    `json:"updatetime"`
	Createtime int    `json:"createtime"`
}

// follow
type userFollowingAPI struct {
	UserID        int64  `json:"user_id"`
	UserName      string `json:"user_name"`
	Name          string `json:"name"`
	PhotoURL      string `json:"photo_url"`
	FollowingTime int    `json:"following_time"`
}
type userFollowerAPI struct {
	UserID        int64  `json:"user_id"`
	UserName      string `json:"user_name"`
	Name          string `json:"name"`
	PhotoURL      string `json:"photo_url"`
	FollowingTime int    `json:"following_time"`
}

// post
type userBasicAPI struct {
	UserID   int64  `json:"user_id" bson:"user_id"`
	Username string `json:"username" bson:"username"`
	Name     string `json:"name" bson:"name"`
}
type blobAPI struct {
	BlobID       string `json:"blob_id" bson:"blob_id"`
	OriginWidth  int    `json:"origin_width" bson:"origin_width"`
	OriginHeight int    `json:"origin_height" bson:"origin_height"`
}
type postAPI struct {
	PostID       int64        `json:"post_id" bson:"post_id"`
	User         userBasicAPI `json:"user" bson:"user"`
	Content      string       `json:"content" bson:"content"`
	Blob         blobAPI      `json:"blob" bson:"blob"`
	Type         int          `json:"type" bson:"type"`
	LikeCount    int64        `json:"like_count" bson:"like_count"`
	DislikeCount int64        `json:"dislike_count" bson:"dislike_count"`
	CommentCount int64        `json:"comment_count" bson:"comment_count"`
	CountryID    int          `json:"country_id" bson:"country_id"`
	CategoryID   int          `json:"category_id" bson:"category_id"`
	Public       bool         `json:"public" bson:"public"`
	Createtime   int          `json:"createtime" bson:"createtime"`
	Updatetime   int          `json:"updatetime" bson:"updatetime"`
}

// hashtag
type hashtagAPI struct {
	HashtagID int64  `json:"hashtag_id"`
	Value     string `json:"value"`
	Count     int64  `json:"count"`
}
type hashtagOnPostAPI struct {
	PostID    int64 `json:"post_id"`
	HashtagID int64 `json:"hashtag_id"`
}

// tag
type tagOnPostAPI struct {
	PostID     int64        `json:"post_id"`
	User       userBasicAPI `json:"user"`
	X          int          `json:"x"`
	Y          int          `json:"y"`
	Valid      bool         `json:"valid"`
	Createtime int          `json:"createtime"`
	Updatetime int          `json:"updatetime"`
}

// comment
type commentAPI struct {
	CommentID    int64        `json:"comment_id"`
	PostID       int64        `json:"post_id"`
	User         userBasicAPI `json:"user"`
	Comment      string       `json:"comment"`
	LikeCount    int64        `json:"like_count"`
	DislikeCount int64        `json:"dislike_count"`
	ReplyCount   int64        `json:"reply_count"`
	Createtime   int          `json:"createtime"`
	Updatetime   int          `json:"updatetime"`
}

// reply
type replyAPI struct {
	ReplyID      int64        `json:"reply_id"`
	CommentID    int64        `json:"comment_id"`
	User         userBasicAPI `json:"user"`
	Reply        string       `json:"reply"`
	LikeCount    int64        `json:"like_count"`
	DislikeCount int64        `json:"dislike_count"`
	Createtime   int          `json:"createtime"`
	Updatetime   int          `json:"updatetime"`
}

// reaction
type reactionOnPostAPI struct {
	PostID     int64        `json:"post_id"`
	User       userBasicAPI `json:"user"`
	ReactionID int          `json:"reaction_id"`
	Createtime int          `json:"createtime"`
}
type reactionOnCommentAPI struct {
	CommentID  int64        `json:"comment_id"`
	User       userBasicAPI `json:"user"`
	ReactionID int          `json:"reaction_id"`
	Createtime int          `json:"createtime"`
}
type reactionOnReplyAPI struct {
	ReplyID    int64        `json:"reply_id"`
	User       userBasicAPI `json:"user"`
	ReactionID int          `json:"reaction_id"`
	Createtime int          `json:"createtime"`
}

// common
type countryAPI struct {
	CountryID   int    `json:"country_id"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
}
type languageAPI struct {
	LanguageID      int    `json:"language_id"`
	DisplayLanguage string `json:"display_language"`
	Value           string `json:"value"`
}
type genderAPI struct {
	GenderID int    `json:"gender_id"`
	Value    string `json:"value"`
}
type reactionAPI struct {
	ReactionID int    `json:"reaction_id"`
	Value      string `json:"value"`
}
type postTypeAPI struct {
	PostTypeID int      `json:"post_type_id"`
	Value      string   `json:"value"`
	FileFormat []string `json:"file_format"`
}
type categoryAPI struct {
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
}

// cronjob
type userPostPopular struct {
	ID         bson.ObjectId `bson:"_id"`
	UserID     int64         `bson:"user_id"`
	CategoryID int           `bson:"category_id"`
	Posts      []postAPI     `bson:"posts"`
}

type userPostPopularRead struct {
	ID            bson.ObjectId `bson:"_id"`
	UserID        int64         `bson:"user_id"`
	CategoryID    int           `bson:"category_id"`
	WeekTimestamp int           `bson:"week_timestamp"`
	Posts         map[int64]int `bson:"posts"` // k: PostID, v: timestamp
}
