package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/federus1105/weekly/internals/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MoviesRepository struct {
	db *pgxpool.Pool
}

func NewMoviesRepository(db *pgxpool.Pool) *MoviesRepository {
	return &MoviesRepository{db: db}
}

func (mr *MoviesRepository) GetUpcomingMovies(rctx context.Context, limit, offset int) ([]models.Movie, error) {
	// ambil data movie
	sql := `SELECT 
    m.id,
    m.image,
    m.title,
    m.release_date,
    STRING_AGG(g.name, ', ') AS genres
FROM movies m
JOIN movies_genre mg ON m.id = mg.id_movies
JOIN genres g ON mg.id_genre = g.id
WHERE m.release_date > CURRENT_DATE
GROUP BY m.id, m.image, m.title, m.release_date
ORDER BY m.release_date ASC
LIMIT $1 OFFSET $2
`
	rows, err := mr.db.Query(rctx, sql, limit, offset)
	if err != nil {
		log.Println("Internal Server Error: ", err.Error())
		return nil, err
	}
	defer rows.Close()
	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.Id, &movie.Image, &movie.Title, &movie.ReleaseDate, &movie.Genres); err != nil {
			log.Println("Internal Server Error: ", err.Error())
			return nil, err
		}
		movies = append(movies, movie)
	}
	return movies, nil
}

func (mr *MoviesRepository) GetPopularMovies(rctx context.Context, limit, offset int) ([]models.Movie, error) {
	sql := `
SELECT 
  m.id,
  m.title,
  m.image,
  STRING_AGG(g.name, ', ') AS genres
FROM movies m
JOIN movies_genre mg ON m.id = mg.id_movies
JOIN genres g ON mg.id_genre = g.id
WHERE m.rating >= 7.0
GROUP BY m.id, m.image, m.title
ORDER BY m.rating ASC
LIMIT $1 OFFSET $2
`
	rows, err := mr.db.Query(rctx, sql, limit, offset)
	if err != nil {
		log.Println("Internal Server Error: ", err.Error())
		return nil, err
	}
	defer rows.Close()
	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.Id, &movie.Title, &movie.Image, &movie.Genres); err != nil {
			log.Println("Internal Server Error: ", err.Error())
			return nil, err
		}
		movies = append(movies, movie)
	}
	return movies, nil
}

func (mr *MoviesRepository) GetFilterMovie(rctx context.Context, limit, offset int) ([]models.Movie, error) {
	sql := `
	SELECT
m.id,	
  m.title,
  m.image,
  STRING_AGG(g.name, ', ') AS genres
FROM movies m
JOIN movies_genre mg ON m.id = mg.id_movies
JOIN genres g ON mg.id_genre = g.id
WHERE 
  LOWER(m.title) ILIKE LOWER('%s%')
  AND LOWER(g.name) = LOWER('Action')
  GROUP BY m.id, m.image, m.title
ORDER BY m.title asc
LIMIT $1 OFFSET $2
`
	rows, err := mr.db.Query(rctx, sql, limit, offset)
	if err != nil {
		log.Println("Internal Server Error: ", err.Error())
		return nil, err
	}
	defer rows.Close()
	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.Id, &movie.Title, &movie.Image, &movie.Genres); err != nil {
			log.Println("Internal Server Error: ", err.Error())
			return nil, err
		}
		movies = append(movies, movie)
	}
	return movies, nil
}

func (mr *MoviesRepository) GetDetailMovie(rctx context.Context, movieID int) ([]models.Movie, error) {
	sql := `SELECT id, image, backdrop, title, release_date, duration, id_director, synopsis
FROM movies where id = $1`

	rows, err := mr.db.Query(rctx, sql, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.Id, &movie.Image, &movie.Backdrop, &movie.Title, &movie.ReleaseDate, &movie.Duration, &movie.Director, &movie.Synopsis); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	return movies, nil
}

func (mr *MoviesRepository) GetAllMovie(rctx context.Context, limit, offset int) ([]models.Movie, error) {
	sql := `SELECT
  m.id,
  m.title,
  m.image,
  STRING_AGG(g.name, ', ') AS genres
FROM movies m
JOIN director d ON m.id_director = d.id
LEFT JOIN movies_genre mg ON m.id = mg.id_movies
LEFT JOIN genres g ON mg.id_genre = g.id
GROUP BY m.id, d.name
ORDER BY m.title ASC
LIMIT $1 OFFSET $2`

	rows, err := mr.db.Query(rctx, sql, limit, offset)
	if err != nil {
		log.Println("Internal Server Error: ", err.Error())
		return nil, err
	}
	defer rows.Close()
	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.Id, &movie.Title, &movie.Image, &movie.Genres); err != nil {
			log.Println("Internal Server Error: ", err.Error())
			return nil, err
		}
		movies = append(movies, movie)
	}
	return movies, nil
}

func (mr *MoviesRepository) DeleteMovie(ctx context.Context, id int) error {
	query := `DELETE FROM movies WHERE id = $1`
	log.Printf("Executing query: %s with id: %d", query, id)

	result, err := mr.db.Exec(ctx, query, id)
	if err != nil {
		log.Printf("failed to execute delete query: %v", err)
		if ctxErr := ctx.Err(); ctxErr != nil {
			log.Printf("context error: %v", ctxErr)
		}
		return err
	}

	rows := result.RowsAffected()
	log.Printf("Rows affected: %d", rows)

	if rows == 0 {
		return fmt.Errorf("movie with id %d not found", id)
	}

	log.Printf("movie with id %d successfully deleted", id)
	return nil
}