package models

import "time"

type Movie struct {
	Id          int       `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Image       string    `db:"image" json:"poster_path"`
	ReleaseDate time.Time `db:"release_date" json:"release_date"`
	Genres      string    `db:"genres" json:"genres"`
	Backdrop    string    `db:"backdrop" json:"backdrop_path"`
	Duration    string    `db:"duration" json:"duration"`
	Synopsis    string    `db:"synopsis" json:"synopsis"`
	Director    string    `db:"id_director" json:"director"`
}
	