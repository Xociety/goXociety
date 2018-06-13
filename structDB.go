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
	PostID       string
	UserID       string
	Content      string
	BlobID       string
	Type         int
	LikeCount    string
	DislikeCount string
	// Point
	CountryID  int
	CategoryID int
	Public     bool
	Createtime int
	Updatetime int
}
