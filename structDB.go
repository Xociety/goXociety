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
	UserID     int64
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
	PostID       int64
	UserID       int64
	Content      string
	BlobID       string
	Type         int
	LikeCount    int64
	DislikeCount int64
	CommentCount int64
	// Point
	CountryID  int
	CategoryID int
	Public     bool
	Createtime int
	Updatetime int
}
