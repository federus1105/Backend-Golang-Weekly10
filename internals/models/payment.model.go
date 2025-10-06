package models

type Payment struct {
	Id int `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Image string `db:"image" json:"image"`
}