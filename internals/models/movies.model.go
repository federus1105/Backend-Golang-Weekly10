package models

import (
	"mime/multipart"
	"time"
)

type Movie struct {
	Id          int       `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Image       string    `db:"image" json:"poster_path,omitempty"`
	ReleaseDate time.Time `db:"release_date" json:"release_date,omitzero"`
	Genres      []string  `db:"genres" json:"genres"`
	Backdrop    string    `db:"backdrop" json:"backdrop_path,omitempty"`
	Duration    string    `db:"duration" json:"duration,omitempty"`
	Synopsis    string    `db:"synopsis" json:"synopsis,omitempty"`
	Director    string    `db:"director" json:"director,omitempty"`
	Rating      float64   `db:"rating" json:"rating,omitempty"`
	Actor       []string  `db:"actor" json:"actor,omitempty"`
}

type MovieAdmin struct {
	Id          int       `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Image       string    `db:"image" json:"poster_path,omitempty"`
	ReleaseDate time.Time `db:"release_date" json:"release_date,omitzero"`
	Genres      string    `db:"genres" json:"genres"`
	Duration    string    `db:"duration" json:"duration,omitempty"`
}

type BodySchedules struct {
	Idmovie    int      `form:"id_movie"`
	Date       []string `form:"date"`
	IdCinema   []int    `form:"id_cinema"`
	IdTime     []int    `form:"id_time"`
	IdLocation []int    `form:"id_location"`
}
type MovieBody struct {
	Id          int                   `form:"id"`
	Title       string                `form:"title"`
	ReleaseDate time.Time             `form:"release_date"`
	Duration    string                `form:"duration"`
	Synopsis    string                `form:"synopsis"`
	Director    int                   `form:"id_director,omitempty"`
	ActorIDs    []int                 `form:"actor_ids"`
	GenreIDs    []int                 `form:"genre_ids"`
	Schedules   []BodySchedules       `form:"schedule"`
	Rating      float64               `form:"rating"`
	Image       *multipart.FileHeader `form:"poster_path"`
	Backdrop    *multipart.FileHeader `form:"backdrop_path"`
	Imagestr    string                `json:"image"`
	Backdropstr string                `json:"backdrop"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
