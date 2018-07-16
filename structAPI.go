package main

import "github.com/globalsign/mgo/bson"

// common
type updateStatusAPI struct {
	RowsAffected int `json:"rows_affected,omitempty"`
}

// login
type loginAPI struct {
	Token string `json:"token,omitempty"`
}

// user
type xuserAPI struct {
	UserID     int64  `json:"user_id,omitempty"`
	Username   string `json:"username,omitempty"`
	Email      string `json:"email,omitempty"`
	Name       string `json:"name,omitempty"`
	Phone      string `json:"phone,omitempty"`
	Gender     int    `json:"gender,omitempty"`
	Bio        string `json:"bio,omitempty"`
	Credit     int    `json:"credit,omitempty"`
	PhotoURL   string `json:"photo_url,omitempty"`
	LanguageID int    `json:"language_id,omitempty"`
	CountryID  int    `json:"country_id,omitempty"`
	Timezone   int    `json:"timezone,omitempty"`
	LastIP     string `json:"last_ip,omitempty"`
	Updatetime int    `json:"updatetime,omitempty"`
	Createtime int    `json:"createtime,omitempty"`
}

// follow
type userFollowingAPI struct {
	UserID        int64  `json:"user_id,omitempty"`
	UserName      string `json:"user_name,omitempty"`
	Name          string `json:"name,omitempty"`
	PhotoURL      string `json:"photo_url,omitempty"`
	FollowingTime int    `json:"following_time,omitempty"`
}
type userFollowerAPI struct {
	UserID        int64  `json:"user_id,omitempty"`
	UserName      string `json:"user_name,omitempty"`
	Name          string `json:"name,omitempty"`
	PhotoURL      string `json:"photo_url,omitempty"`
	FollowingTime int    `json:"following_time,omitempty"`
}

// post
type userBasicAPI struct {
	UserID   int64  `json:"user_id,omitempty" bson:"user_id"`
	Username string `json:"username,omitempty" bson:"username"`
	Name     string `json:"name,omitempty" bson:"name"`
}
type blobAPI struct {
	BlobID       string `json:"blob_id,omitempty" bson:"blob_id"`
	OriginWidth  int    `json:"origin_width,omitempty" bson:"origin_width"`
	OriginHeight int    `json:"origin_height,omitempty" bson:"origin_height"`
}
type postAPI struct {
	PostID       int64        `json:"post_id,omitempty" bson:"post_id"`
	User         userBasicAPI `json:"user,omitempty" bson:"user"`
	Content      string       `json:"content,omitempty" bson:"content"`
	Blob         blobAPI      `json:"blob,omitempty" bson:"blob"`
	Type         int          `json:"type,omitempty" bson:"type"`
	LikeCount    int64        `json:"like_count,omitempty" bson:"like_count"`
	DislikeCount int64        `json:"dislike_count,omitempty" bson:"dislike_count"`
	CommentCount int64        `json:"comment_count,omitempty" bson:"comment_count"`
	CountryID    int          `json:"country_id,omitempty" bson:"country_id"`
	CategoryID   int          `json:"category_id,omitempty" bson:"category_id"`
	Public       bool         `json:"public,omitempty" bson:"public"`
	Createtime   int          `json:"createtime,omitempty" bson:"createtime"`
	Updatetime   int          `json:"updatetime,omitempty" bson:"updatetime"`
}

// comment
type commentAPI struct {
	CommentID    int64        `json:"comment_id,omitempty"`
	PostID       int64        `json:"post_id,omitempty"`
	User         userBasicAPI `json:"user,omitempty"`
	Comment      string       `json:"comment,omitempty"`
	LikeCount    int64        `json:"like_count,omitempty"`
	DislikeCount int64        `json:"dislike_count,omitempty"`
	CommentCount int64        `json:"comment_count,omitempty"`
	Createtime   int          `json:"createtime,omitempty"`
	Updatetime   int          `json:"updatetime,omitempty"`
}

// reaction
type reactionOnPostAPI struct {
	PostID     int64        `json:"post_id,omitempty"`
	User       userBasicAPI `json:"user,omitempty"`
	ReactionID int          `json:"reaction_id,omitempty"`
	Createtime int          `json:"createtime,omitempty"`
}
type reactionOnCommentAPI struct {
	CommentID  int64        `json:"comment_id,omitempty"`
	User       userBasicAPI `json:"user,omitempty"`
	ReactionID int          `json:"reaction_id,omitempty"`
	Createtime int          `json:"createtime,omitempty"`
}
type reactionOnReplyAPI struct {
	ReplyID    int64        `json:"reply_id,omitempty"`
	User       userBasicAPI `json:"user,omitempty"`
	ReactionID int          `json:"reaction_id,omitempty"`
	Createtime int          `json:"createtime,omitempty"`
}

// common
type countryAPI struct {
	CountryID   int    `json:"country_id,omitempty"`
	Country     string `json:"country,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
}
type languageAPI struct {
	LanguageID      int    `json:"language_id,omitempty"`
	DisplayLanguage string `json:"display_language,omitempty"`
	Value           string `json:"value,omitempty"`
}
type genderAPI struct {
	GenderID int    `json:"gender_id,omitempty"`
	Value    string `json:"value,omitempty"`
}
type reactionAPI struct {
	ReactionID int    `json:"reaction_id,omitempty"`
	Value      string `json:"value,omitempty"`
}
type postTypeAPI struct {
	PostTypeID int      `json:"post_type_id,omitempty"`
	Value      string   `json:"value,omitempty"`
	FileFormat []string `json:"file_format,omitempty"`
}
type categoryAPI struct {
	CategoryID   int    `json:"category_id,omitempty"`
	CategoryName string `json:"category_name,omitempty"`
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
