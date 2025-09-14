package models

type User struct {
	Id       int    `db:"id" json:"id"`
	Email    string `db:"email" json:"email" binding:"required"`
	Password string `db:"password" json:"password" binding:"required, min=4"`
	Role     string `db:"role" json:"role"`
	// Image    string `db:"image" json:"image"`
}

type UserRegister struct {
	Id       int    `json:"id,omitempty"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role" binding:"required"`
}

// type UserBody struct {
// 	User
// 	Images *multipart.FileHeader `form:"image"`
// }

type UserAuth struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChangePassword struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}
