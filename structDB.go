package main

type country struct {
	CountryID int
	Name      string
	Code      string
}

type language struct {
	LanguageID      int
	DisplayLanguage string
	Value           string
}

type xuserDB struct {
	UserID     string `json:"user_id,omitempty"`
	Username   string `json:"username,omitempty"`
	Email      string `json:"email,omitempty"`
	Password   string `json:"password,omitempty"`
	Name       string `json:"name,omitempty"`
	Phone      string `json:"phone,omitempty"`
	Gender     int    `json:"gender,omitempty"`
	Bio        string `json:"bio,omitempty"`
	Credit     int    `json:"credit,omitempty"`
	LanguageID int    `json:"language_id,omitempty"`
	CountryID  int    `json:"country_id,omitempty"`
	Timezone   int    `json:"timezone,omitempty"`
	LastIP     string `json:"last_ip,omitempty"`
	Updatetime int    `json:"updatetime,omitempty"`
	Createtime int    `json:"createtime,omitempty"`
}

type postDB struct {
	PostID  string `json:"post_id,omitempty"`
	UserID  string `json:"user_id,omitempty"`
	Content string `json:"content,omitempty"`
	BlobID  string `json:"blob_id,omitempty"`
	// Point
	CountryID  int  `json:"country_id,omitempty"`
	CategoryID int  `json:"category_id,omitempty"`
	Public     bool `json:"public,omitempty"`
	Type       int  `json:"type,omitempty"`
	Createtime int  `json:"createtime,omitempty"`
	Updatetime int  `json:"updatetime,omitempty"`
}

type hashtagDB struct {
	HashtagID string
	Name      string
}

type postHashtagDB struct {
	PostID    string
	HashtagID string
}

type postLikesDB struct {
	PostID     string
	UserID     string
	Type       int
	Createtime int
}

type commentsDB struct {
	CommentID  string
	PostID     string
	UserID     string
	Comment    string
	Createtime int
	Updatetime int
}

type commentLikesDB struct {
	CommentID  string
	UserID     string
	Type       int
	Createtime int
}
