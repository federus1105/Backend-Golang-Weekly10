package models

import (
	"mime/multipart"
	"time"

	"github.com/federus1105/weekly/internals/utils"
)

type Movie struct {
	Id          int       `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Image       string    `db:"image" json:"poster_path,omitempty"`
	ReleaseDate time.Time `db:"release_date" json:"release_date,omitzero"`
	Genres      string    `db:"genres" json:"genres"`
	Backdrop    string    `db:"backdrop" json:"backdrop_path,omitempty"`
	Duration    string    `db:"duration" json:"duration,omitempty"`
	Synopsis    string    `db:"synopsis" json:"synopsis,omitempty"`
	Director    string    `db:"id_director" json:"director,omitempty"`
	// Rating      float64   `db:"rating" json:"rating,omitempty"`
	Actor string `db:"actor" json:"actor,omitempty"`
}

type MovieBody struct {
	Id           int                   `form:"id"`
	Title        string                `form:"title" binding:"required"`
	ReleaseDate  utils.DateOnly        `form:"release_date" binding:"required"`
	Duration     string                `form:"duration" binding:"required"`
	Synopsis     string                `form:"synopsis" binding:"required"`
	Director     int                   `form:"id_director" binding:"required"`
	ActorIDs     []int                 `form:"actor_ids" binding:"required"`
	GenreIDs     []int                 `form:"genre_ids" binding:"required"`
	Rating       float64               `form:"rating" binding:"required"`
	Image        *multipart.FileHeader `form:"poster_path" binding:"required"`
	PosterPath   string                `json:"-"`
	Backdrop     *multipart.FileHeader `form:"backdrop_path" binding:"required"`
	BackdropPath string                `json:"-"`
}

// type MovieBody struct {
// 	Title    string                `form:"title"`
// 	Duration string                `form:"duration"`
// 	Synopsis string                `form:"synopsis"`
// 	Image    *multipart.FileHeader `form:"image"`
// }
