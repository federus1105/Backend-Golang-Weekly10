package models

import (
	"time"
)

type Schedule struct {
	Id        int    `db:"id" json:"id"`
	Idmovie   int    `db:"id_movie" json:"id_movie"`
	Date      string `db:"date" json:"date"`
	Title     string `db:"title" json:"title"`
	Image     string `db:"image" json:"image_cinema"`
	Id_Cinema int    `db:"id" json:"id_cinema"`
	Cinema    string `db:"cinema" json:"cinema"`
	Time      string `db:"id_time" json:"time"`
	Location  string `db:"id_location" json:"tocation"`
}

type BodyScheduleInput struct {
	Id        int    `db:"id" json:"id"`
	Date      string `db:"date" json:"date"`
	Id_movie  int    `json:"id_movie"`
	Id_Cinema []int  `json:"id_cinema"`
	Time      []int  `json:"id_time"`
	Location  []int  `json:"id_location"`
}
type BodySchedule struct {
	Id          int       `db:"id" json:"id"`
	Id_movie    int       `db:"id_movie" json:"id_movie"`
	Date        time.Time `db:"date" json:"date"`
	Id_Cinema   int       `db:"id_cinema" json:"id_cinema"`
	Id_Time     int       `db:"id_time" json:"id_time"`
	Id_Location int       `db:"id_location" json:"id_location"`
}
