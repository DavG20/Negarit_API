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
type DBResponseFailed struct {
	Message string `json:"message"`
}

type UpdateUserProfile struct {
	Userprofile string `json:"userprofile"`
	Bio         string `json:"bio"`
}

// custom struct not to expose credentials
type UserProfileUpdateResponse struct {
	Message string       `json:"message"`
	Success bool         `json:"success"`
	User    []DBResponse `json:"user"`
}

type UpdatePassword struct {
	OldPassword string `json:"password"`
	NewPassword string `json:"newpassword"`
}

type UserProfile struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	UserPro []DBResponse
}
