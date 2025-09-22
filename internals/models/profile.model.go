package models

import "mime/multipart"

type Profile struct {
	UserID    int    `db:"user_id" json:"id"`
	Email     string `db:"email" json:"email,omitempty"`
	Password  string `db:"password" json:"password,omitempty"`
	Image     *string `db:"image" json:"image,omitempty"`
	FirstName string `db:"firstname" json:"first_name"`
	LastName  string `db:"lastname" json:"last_name"`
	Phone     string `db:"phonenumber" json:"phone"`
	Point     string `db:"point" json:"point"`
}

type ProfileBody struct {
	FirstName *string               `form:"first_name"`
	LastName  *string               `form:"last_name"`
	Phone     *string               `form:"phone"`
	Image     *multipart.FileHeader `form:"image"`
}
