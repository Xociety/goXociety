package config

import "github.com/globalsign/mgo/bson"

// common
type UpdateStatusAPI struct {
	RowsAffected int `json:"rows_affected"`
}

// login
type LoginAPI struct {
	Token string `json:"token"`
}

// user
type XuserAPI struct {
	UserID      int64  `json:"user_id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Gender      int    `json:"gender"`
	Bio         string `json:"bio"`
	Credit      int    `json:"credit"`
	PhotoURL    string `json:"photo_url"`
	LanguageID  int    `json:"language_id"`
	CountryCode string `json:"country_code"`
	Timezone    int    `json:"timezone"`
	LastIP      string `json:"last_ip"`
	Updatetime  int    `json:"updatetime"`
	Createtime  int    `json:"createtime"`
}

// follow
type UserFollowingAPI struct {
	UserID        int64  `json:"user_id"`
	UserName      string `json:"user_name"`
	Name          string `json:"name"`
	PhotoURL      string `json:"photo_url"`
	FollowingTime int    `json:"following_time"`
}
type UserFollowerAPI struct {
	UserID        int64  `json:"user_id"`
	UserName      string `json:"user_name"`
	Name          string `json:"name"`
	PhotoURL      string `json:"photo_url"`
	FollowingTime int    `json:"following_time"`
}

// city
type PropertiesCityAPI struct {
	CountryName string `json:"country_name,omitempty" bson:"country_name,omitempty"`
	Name1       string `json:"name_1,omitempty" bson:"name_1,omitempty"`
	Name2       string `json:"name_2,omitempty" bson:"name_2,omitempty"`
	Name3       string `json:"name_3,omitempty" bson:"name_3,omitempty"`
	Name4       string `json:"name_4,omitempty" bson:"name_4,omitempty"`
	Name5       string `json:"name_5,omitempty" bson:"name_5,omitempty"`
	CountryCode string `json:"country_code,omitempty" bson:"country_code,omitempty"`
	CityID1     string `json:"city_id_1,omitempty" bson:"city_id_1,omitempty"`
	CityID2     string `json:"city_id_2,omitempty" bson:"city_id_2,omitempty"`
	CityID3     string `json:"city_id_3,omitempty" bson:"city_id_3,omitempty"`
	CityID4     string `json:"city_id_4,omitempty" bson:"city_id_4,omitempty"`
	CityID5     string `json:"city_id_5,omitempty" bson:"city_id_5,omitempty"`
	Type1       string `json:"type_1,omitempty" bson:"type_1,omitempty"`
	Type2       string `json:"type_2,omitempty" bson:"type_2,omitempty"`
	Type3       string `json:"type_3,omitempty" bson:"type_3,omitempty"`
	Type4       string `json:"type_4,omitempty" bson:"type_4,omitempty"`
	Type5       string `json:"type_5,omitempty" bson:"type_5,omitempty"`
}
type CityGeometryAPI struct {
	Properties PropertiesCityAPI `json:"properties" bson:"properties"`
}
type CityAPI struct {
	Level           string    `json:"level" bson:"level"`
	CountryCode     string    `json:"country_code" bson:"country_code"`
	CountryName     string    `json:"country_name" bson:"country_name"`
	CityID1         string    `json:"city_id_1" bson:"city_id_1"`
	CityID2         string    `json:"city_id_2" bson:"city_id_2"`
	CityID3         string    `json:"city_id_3" bson:"city_id_3"`
	CityID4         string    `json:"city_id_4" bson:"city_id_4"`
	CityID5         string    `json:"city_id_5" bson:"city_id_5"`
	Name            string    `json:"name" bson:"name"`
	Type            string    `json:"type" bson:"type"`
	PostCount       int       `json:"post_count" bson:"post_count"`
	SupPopularPosts []PostAPI `json:"sup_popular_posts" bson:"sup_popular_posts"`
}

// place
type PlacesLookupAPI struct {
	Place         []PlaceAPI `json:"place"`
	NextPageToken string     `json:"next_page_token"`
}
type PlaceAPI struct {
	PlaceID     int64   `json:"place_id" bson:"place_id"`
	CountryCode string  `json:"country_code" bson:"country_code"`
	CityID1     string  `json:"city_id_1" bson:"city_id_1"`
	CityID2     string  `json:"city_id_2" bson:"city_id_2"`
	CityID3     string  `json:"city_id_3" bson:"city_id_3"`
	CityID4     string  `json:"city_id_4" bson:"city_id_4"`
	CityID5     string  `json:"city_id_5" bson:"city_id_5"`
	Lat         float64 `json:"lat" bson:"lat"`
	Lon         float64 `json:"lon" bson:"lon"`
	Name        string  `json:"name" bson:"name"`
}

// post
type UserBasicAPI struct {
	UserID   int64  `json:"user_id" bson:"user_id"`
	Username string `json:"username" bson:"username"`
	Name     string `json:"name" bson:"name"`
	PhotoURL string `json:"photo_url" bson:"photo_url"`
}
type BlobAPI struct {
	BlobID       string `json:"blob_id" bson:"blob_id"`
	OriginWidth  int    `json:"origin_width" bson:"origin_width"`
	OriginHeight int    `json:"origin_height" bson:"origin_height"`
}
type PostAPI struct {
	PostID              int64        `json:"post_id" bson:"post_id"`
	User                UserBasicAPI `json:"user" bson:"user"`
	Content             string       `json:"content" bson:"content"`
	Blob                BlobAPI      `json:"blob" bson:"blob"`
	Type                int          `json:"type" bson:"type"`
	LikeCount           int64        `json:"like_count" bson:"like_count"`
	DislikeCount        int64        `json:"dislike_count" bson:"dislike_count"`
	ReactionByQueryUser int          `json:"reaction_by_query_user" bson:"reaction_by_query_user"`
	CommentCount        int64        `json:"comment_count" bson:"comment_count"`
	CategoryID          int          `json:"category_id" bson:"category_id"`
	Place               PlaceAPI     `json:"place" bson:"place"`
	Public              bool         `json:"public" bson:"public"`
	Createtime          int          `json:"createtime" bson:"createtime"`
	Updatetime          int          `json:"updatetime" bson:"updatetime"`
}

// hashtag
type HashtagAPI struct {
	HashtagID int64  `json:"hashtag_id"`
	Value     string `json:"value"`
	Count     int64  `json:"count"`
}
type HashtagOnPostAPI struct {
	PostID    int64 `json:"post_id"`
	HashtagID int64 `json:"hashtag_id"`
}

// tag
type TagOnPostSetAPI struct {
	PostID     int64 `json:"post_id"`
	UserID     int64 `json:"user_id"`
	X          int   `json:"x"`
	Y          int   `json:"y"`
	Valid      bool  `json:"valid"`
	Createtime int   `json:"createtime"`
	Updatetime int   `json:"updatetime"`
}
type TagOnPostAPI struct {
	PostID     int64        `json:"post_id"`
	User       UserBasicAPI `json:"user"`
	X          int          `json:"x"`
	Y          int          `json:"y"`
	Valid      bool         `json:"valid"`
	Createtime int          `json:"createtime"`
	Updatetime int          `json:"updatetime"`
}

// comment
type CommentAPI struct {
	CommentID    int64        `json:"comment_id"`
	PostID       int64        `json:"post_id"`
	User         UserBasicAPI `json:"user"`
	Comment      string       `json:"comment"`
	LikeCount    int64        `json:"like_count"`
	DislikeCount int64        `json:"dislike_count"`
	ReplyCount   int64        `json:"reply_count"`
	Createtime   int          `json:"createtime"`
	Updatetime   int          `json:"updatetime"`
}

// reply
type ReplyAPI struct {
	ReplyID      int64        `json:"reply_id"`
	CommentID    int64        `json:"comment_id"`
	User         UserBasicAPI `json:"user"`
	Reply        string       `json:"reply"`
	LikeCount    int64        `json:"like_count"`
	DislikeCount int64        `json:"dislike_count"`
	Createtime   int          `json:"createtime"`
	Updatetime   int          `json:"updatetime"`
}

// reaction
type ReactionOnPostAPI struct {
	PostID     int64        `json:"post_id"`
	User       UserBasicAPI `json:"user"`
	ReactionID int          `json:"reaction_id"`
	Createtime int          `json:"createtime"`
}
type ReactionOnCommentAPI struct {
	CommentID  int64        `json:"comment_id"`
	User       UserBasicAPI `json:"user"`
	ReactionID int          `json:"reaction_id"`
	Createtime int          `json:"createtime"`
}
type ReactionOnReplyAPI struct {
	ReplyID    int64        `json:"reply_id"`
	User       UserBasicAPI `json:"user"`
	ReactionID int          `json:"reaction_id"`
	Createtime int          `json:"createtime"`
}

// common
type LanguageAPI struct {
	LanguageID      int    `json:"language_id"`
	DisplayLanguage string `json:"display_language"`
	Value           string `json:"value"`
}
type GenderAPI struct {
	GenderID int    `json:"gender_id"`
	Value    string `json:"value"`
}
type ReactionAPI struct {
	ReactionID int    `json:"reaction_id"`
	Value      string `json:"value"`
}
type PostTypeAPI struct {
	PostTypeID int      `json:"post_type_id"`
	Value      string   `json:"value"`
	FileFormat []string `json:"file_format"`
}
type CategoryAPI struct {
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
}

// cronjob

type PostCommonAPI struct {
	ID           bson.ObjectId `bson:"_id"`
	CategoryID   int           `bson:"category_id"`
	PopularPosts []PostAPI     `bson:"popular_posts"`
}

type PostUserReadAPI struct {
	ID            bson.ObjectId `bson:"_id"`
	UserID        int64         `bson:"user_id"`
	CategoryID    int           `bson:"category_id"`
	WeekTimestamp int           `bson:"week_timestamp"`
	PopularPosts  map[int64]int `bson:"popular_posts"` // k: PostID, v: timestamp
}

type PopularPostUserReadIndexAPI struct {
	UserID                     int64          `json:"user_id" bson:"user_id"`
	CommonPopularPostIndex     map[int]int    `json:"common_popular_post_index" bson:"common_popular_post_index"`           // category_id: index
	CitySupPopularPostIndex    map[string]int `json:"city_sup_popular_post_index" bson:"city_sup_popular_post_index"`       // city_id: index
	CountrySupPopularPostIndex map[string]int `json:"country_sup_popular_post_index" bson:"country_sup_popular_post_index"` // country_code: index
}
