package models

type Profile struct {
	ID        int    `db:"id" json:"id"`
	Email     string `db:"email" json:"email"`
	Image     string `db:"image" json:"image"`
	FisrtName string `db:"firstname" json:"first_name"`
	LastName  string `db:"lastname" json:"Last_name"`
	Phone     string `db:"phonenumber" json:"phone"`
}
