package models

import "mime/multipart"

type Profile struct {
	UserID    int `db:"user_id" json:"id"`
	Email     string `db:"email" json:"email,omitempty"`
	Password  string `db:"password" json:"password,omitempty"`
	Image     string `db:"image" json:"image"`
	FirstName string `db:"firstname" json:"first_name"`
	LastName  string `db:"lastname" json:"last_name"`
	Phone     string `db:"phonenumber" json:"phone"`
}

type ProfileBody struct {
	FirstName string                `form:"first_name" binding:"required"`
	LastName  string                `form:"last_name" binding:"required"`
	Phone     string                `form:"phone" binding:"required"`
	Image     *multipart.FileHeader `form:"image" binding:"required"`
}
