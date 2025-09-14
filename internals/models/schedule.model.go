package models

import (
	"time"
)

type Schedule struct {
	Id       int       `db:"id" json:"id"`
	Date     time.Time `db:"date" json:"date"`
	Title    string    `db:"title" json:"title"`
	Image    string    `db:"image" json:"image"`
	Cinema   string    `db:"cinema" json:"cinema"`
	Time     string    `db:"id_time" json:"time"`
	Location string    `db:"id_location" json:"tocation"`
}
