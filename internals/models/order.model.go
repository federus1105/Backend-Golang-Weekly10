package models

type Order struct {
	Id       int     `json:"id,omitempty"`
	Schedule *int    `json:"schedule" binding:"required"`
	User     *int    `json:"user" binding:"required"`
	Payment  *int    `json:"payment" binding:"required"`
	Total    float32 `json:"total"`
	Fullname string  `json:"fullname" binding:"required"`
	Email    string  `json:"email" binding:"required,email"`
	Phone    string  `json:"phone" binding:"required"`
	Paid     bool    `json:"paid"`
}
