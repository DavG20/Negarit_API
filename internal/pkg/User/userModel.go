package user

import "time"

type User struct {
	UserId      string    `bson:"_id,omitempty"  json:"id,omitempty"`
	Email       string    `json:"email,omitempty"`
	Username    string    `json:"username,omitempty"`
	Password    string    `json:"password,omitempty"`
	Userprofile string    `json:"userprofile,omitempty"`
	Bio         string    `json:"bio,omitempty"`
	CreatedAt   time.Time `json:"createdat,omitempty"`
}

type SignUpInput struct {
	Email       string `json:"email,omitempty"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	Userprofile string `json:"userprofile,omitempty"`
	Bio         string `json:"bio,omitempty"`
}

type SignInInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DBResponse struct{
	Email string `json:"email"`
	Username string `json:"username"`
	Userprofile string `json:"userprofile"`
	Bio string `json:"bio"`

}