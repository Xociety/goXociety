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
	UserID     string
	Username   string
	Email      string
	Password   string
	Name       string
	Phone      string
	Gender     int
	Bio        string
	Credit     int
	PhotoURL   string
	LanguageID int
	CountryID  int
	Timezone   int
	LastIP     string
	Updatetime int
	Createtime int
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
