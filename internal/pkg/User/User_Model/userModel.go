package usermodel

import "time"

type User struct {
	UserId      string    `bson:"_id,omitempty"  json:"id,omitempty"`
	Email       string    `json:"email,omitempty"`
	Username    string    `json:"username,omitempty"`
	Password    string    `json:"password,omitempty"`
	Userprofile string    `json:"userprofile,omitempty"`
	Bio         string    `json:"bio,omitempty"`
	CreatedOn   time.Time `json:"createdon,omitempty"`
}

type SignUpInput struct {
	Email       string    `json:"email,omitempty"`
	Username    string    `json:"username,omitempty"`
	Password    string    `json:"password,omitempty"`
	Userprofile string    `json:"userprofile,omitempty"`
	Bio         string    `json:"bio,omitempty"`
	CreatedOn   time.Time `json:"createdon,omitempty"`
}

type SignInInput struct { //auth input
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DBResponse struct {
	UserId      string    `bson:"_id,omitempty"  json:"id,omitempty"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	Userprofile string    `json:"userprofile"`
	Bio         string    `json:"bio"`
	CreatedOn   time.Time `json:"createdon"`
}
type DBResponseFailed struct{ 
    Message string `json:"message"`
	
}
