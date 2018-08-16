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

// city
type geometryCityAPI struct {
	Type        string        `bson:"type"`
	Coordinates [][][]float64 `bson:"coordinates"`
}
type propertiesCityAPI struct {
	Name0 string `json:"name_0,omitempty" bson:"NAME_0,omitempty"`
	Name1 string `json:"name_1,omitempty" bson:"NAME_1,omitempty"`
	Name2 string `json:"name_2,omitempty" bson:"NAME_2,omitempty"`
	Name3 string `json:"name_3,omitempty" bson:"NAME_3,omitempty"`
	Name4 string `json:"name_4,omitempty" bson:"NAME_4,omitempty"`
	Name5 string `json:"name_5,omitempty" bson:"NAME_5,omitempty"`
	GID0  string `json:"gid_0,omitempty" bson:"GID_0,omitempty"`
	GID1  string `json:"gid_1,omitempty" bson:"GID_1,omitempty"`
	GID2  string `json:"gid_2,omitempty" bson:"GID_2,omitempty"`
	GID3  string `json:"gid_3,omitempty" bson:"GID_3,omitempty"`
	GID4  string `json:"gid_4,omitempty" bson:"GID_4,omitempty"`
	GID5  string `json:"gid_5,omitempty" bson:"GID_5,omitempty"`
	Type0 string `json:"type_0,omitempty" bson:"TYPE_0,omitempty"`
	Type1 string `json:"type_1,omitempty" bson:"TYPE_1,omitempty"`
	Type2 string `json:"type_2,omitempty" bson:"TYPE_2,omitempty"`
	Type3 string `json:"type_3,omitempty" bson:"TYPE_3,omitempty"`
	Type4 string `json:"type_4,omitempty" bson:"TYPE_4,omitempty"`
	Type5 string `json:"type_5,omitempty" bson:"TYPE_5,omitempty"`
}
type cityAPI struct {
	Properties propertiesCityAPI `json:"properties" bson:"properties"`
}
type cityLevelAPI struct {
	Name      string `json:"name" bson:"NAME"`
	GID       string `json:"git" bson:"GID"`
	Type      string `json:"type" bson:"TYPE"`
	PostCount int    `json:"post_count" bson:"POST_COUNT"`
}

// post
type userBasicAPI struct {
	UserID   int64  `json:"user_id" bson:"user_id"`
	Username string `json:"username" bson:"username"`
	Name     string `json:"name" bson:"name"`
	PhotoURL string `json:"photo_url" bson:"photo_url"`
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
type tagOnPostSetAPI struct {
	PostID     int64 `json:"post_id"`
	UserID     int64 `json:"user_id"`
	X          int   `json:"x"`
	Y          int   `json:"y"`
	Valid      bool  `json:"valid"`
	Createtime int   `json:"createtime"`
	Updatetime int   `json:"updatetime"`
}
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
