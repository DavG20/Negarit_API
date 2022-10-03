package user

type User struct {
	UserId      string `bson:"_id,omitempty"  json:"id,omitempty"`
	Email       string `json:"email,omitempty"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	Userprofile string `json:"userprofile,omitempty"`
	Bio         string `json:"bio,omitempty"`
}


