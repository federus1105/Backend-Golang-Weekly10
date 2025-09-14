package models

type Order struct {
	Id       int     `json:"id,omitempty"`
	Schedule int     `json:"schedule" binding:"required"`
	User     int     `json:"user,omitempty"`
	Payment  int     `json:"payment" binding:"required"`
	Total    float32 `json:"total,omitempty"`
	Fullname string  `json:"fullname" binding:"required"`
	Email    string  `json:"email" binding:"required,email"`
	Phone    string  `json:"phone" binding:"required"`
	Paid     bool    `json:"paid" binding:"required"`
	Seats    []int   `json:"seats" binding:"required"`
}

// type Order struct {
// 	Schedule int     `json:"id_schedule"`
// 	Payment  int     `json:"id_payment_method"`
// 	Total    float64 `json:"total"`
// 	Fullname string  `json:"fullname"`
// 	Email    string  `json:"email"`
// 	Phone    string  `json:"phone_number"`
// 	Paid     bool    `json:"paid"`
// 	Seats    []int   `json:"seats"`
// }
