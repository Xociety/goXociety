package main

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
type postAPI struct {
	PostID       int64  `json:"post_id,omitempty"`
	UserID       int64  `json:"user_id,omitempty"`
	Username     string `json:"username,omitempty"`
	Name         string `json:"name,omitempty"`
	Content      string `json:"content,omitempty"`
	BlobID       string `json:"blob_id,omitempty"`
	Type         int    `json:"type,omitempty"`
	LikeCount    int64  `json:"like_count,omitempty"`
	DislikeCount int64  `json:"dislike_count,omitempty"`
	CommentCount int64  `json:"comment_count,omitempty"`
	// Point
	CountryID  int  `json:"country_id,omitempty"`
	CategoryID int  `json:"category_id,omitempty"`
	Public     bool `json:"public,omitempty"`
	Createtime int  `json:"createtime,omitempty"`
	Updatetime int  `json:"updatetime,omitempty"`
}

// comment
type commentAPI struct {
	CommentID    int64  `json:"comment_id,omitempty"`
	PostID       int64  `json:"post_id,omitempty"`
	UserID       int64  `json:"user_id,omitempty"`
	Username     string `json:"username,omitempty"`
	Name         string `json:"name,omitempty"`
	Comment      string `json:"comment,omitempty"`
	LikeCount    int64  `json:"like_count,omitempty"`
	DislikeCount int64  `json:"dislike_count,omitempty"`
	CommentCount int64  `json:"comment_count,omitempty"`
	Createtime   int    `json:"createtime,omitempty"`
	Updatetime   int    `json:"updatetime,omitempty"`
}

// reaction
type reactionOnPostAPI struct {
	PostID     int64  `json:"post_id,omitempty"`
	UserID     int64  `json:"user_id,omitempty"`
	Username   string `json:"username,omitempty"`
	Name       string `json:"name,omitempty"`
	ReactionID int    `json:"reaction_id,omitempty"`
	Createtime int    `json:"createtime,omitempty"`
}
type reactionOnCommentAPI struct {
	CommentID  int64  `json:"comment_id,omitempty"`
	UserID     int64  `json:"user_id,omitempty"`
	Username   string `json:"username,omitempty"`
	Name       string `json:"name,omitempty"`
	ReactionID int    `json:"reaction_id,omitempty"`
	Createtime int    `json:"createtime,omitempty"`
}
type reactionOnThreadAPI struct {
	ThreadID   int64  `json:"thread_id,omitempty"`
	UserID     int64  `json:"user_id,omitempty"`
	Username   string `json:"username,omitempty"`
	Name       string `json:"name,omitempty"`
	ReactionID int    `json:"reaction_id,omitempty"`
	Createtime int    `json:"createtime,omitempty"`
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
	PostTypeID int    `json:"post_type_id,omitempty"`
	Value      string `json:"value,omitempty"`
}
