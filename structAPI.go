package main

type loginAPI struct {
	Token string `json:"token,omitempty"`
}

type xuserAPI struct {
	UserID     string `json:"user_id,omitempty"`
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

type xuserFollowingAPI struct {
	UserID        string `json:"user_id,omitempty"`
	UserName      string `json:"user_name,omitempty"`
	Name          string `json:"name,omitempty"`
	PhotoURL      string `json:"photo_url,omitempty"`
	FollowingTime int    `json:"following_time,omitempty"`
}
type xuserFollowerAPI struct {
	UserID        string `json:"user_id,omitempty"`
	UserName      string `json:"user_name,omitempty"`
	Name          string `json:"name,omitempty"`
	PhotoURL      string `json:"photo_url,omitempty"`
	FollowingTime int    `json:"following_time,omitempty"`
}
