package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/federus1105/weekly/internals/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type MoviesRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewMoviesRepository(db *pgxpool.Pool, rdb *redis.Client) *MoviesRepository {
	return &MoviesRepository{db: db, rdb: rdb}
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

func (mr *MoviesRepository) GetFilterMovie(rctx context.Context, title string, genre string, limit, offset int) ([]models.Movie, error) {
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
        LOWER(m.title) ILIKE LOWER($1)
        AND LOWER(g.name) = LOWER($2)
    GROUP BY m.id, m.image, m.title
    ORDER BY m.title ASC
    LIMIT $3 OFFSET $4
    `
	titlePattern := "%" + title + "%"

	rows, err := mr.db.Query(rctx, sql, titlePattern, genre, limit, offset)
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
	sql := `SELECT 
  m.id,
  m.image,
  m.backdrop,
  m.title,
  m.release_date,
  m.duration,
  STRING_AGG(DISTINCT d.name, ', ') AS director,
  m.synopsis,
  STRING_AGG(DISTINCT g.name, ', ') AS genres,
  STRING_AGG(DISTINCT a.name, ', ') AS actor
FROM movies m
JOIN movies_genre mg ON m.id = mg.id_movies
JOIN genres g ON mg.id_genre = g.id
JOIN movies_actor ma ON m.id = ma.id_movie
JOIN actor a ON ma.id_actor = a.id
JOIN director d ON m.id_director = d.id
WHERE m.id = $1
GROUP BY 
  m.id, 
  m.image, 
  m.backdrop, 
  m.title, 
  m.release_date, 
  m.duration, 
  m.synopsis;`

	rows, err := mr.db.Query(rctx, sql, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.Id, &movie.Image, &movie.Backdrop, &movie.Title, &movie.ReleaseDate, &movie.Duration, &movie.Director, &movie.Synopsis, &movie.Genres, &movie.Actor); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	return movies, nil
}

func (mr *MoviesRepository) GetAllMovie(rctx context.Context, limit, offset int) ([]models.Movie, error) {
	start := time.Now()
	redisKey := "firdaus:allmovies"
	if offset == 0 {
		cmd := mr.rdb.Get(rctx, redisKey)
		if cmd.Err() != nil {
			if cmd.Err() == redis.Nil {
				log.Printf("Key %s does not exist\n", redisKey)
			} else {
				log.Println("Redis Error. \nCause: ", cmd.Err().Error())
			}
		} else {
			// cache hit
			var cachedSchedules []models.Movie
			cmdByte, err := cmd.Bytes()
			if err != nil {
				log.Println("Internal server error.\nCause: ", err.Error())
			} else {
				if err := json.Unmarshal(cmdByte, &cachedSchedules); err != nil {
					log.Println("Internal Server Error. \nCause: ", err.Error())
				}
			}
			if len(cachedSchedules) > 0 {
				log.Printf("Key %s found in cache ✅", redisKey)
				log.Printf("Served in %s using Redis", time.Since(start))
				return cachedSchedules, nil

			}
		}
	}
	sql := `SELECT
  m.id,
  m.title,
  m.image,
  STRING_AGG(g.name, ', ') AS genres
FROM movies m
JOIN director d ON m.id_director = d.id
LEFT JOIN movies_genre mg ON m.id = mg.id_movies
LEFT JOIN genres g ON mg.id_genre = g.id
WHERE is_deleted = FALSE
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
	// renew cache
	if offset == 0 {
		bt, err := json.Marshal(movies)
		if err != nil {
			log.Println("Internal Server Error.\n Cause: ", err.Error())
		}
		if err := mr.rdb.Set(rctx, redisKey, string(bt), 1*time.Minute).Err(); err != nil {
			log.Println("Redis Error. \nCause: ", err.Error())
		}
	}
	log.Printf("[REDIS TIMING] Served in %s using DB (cache miss)", time.Since(start))
	return movies, nil
}

func (mr *MoviesRepository) DeleteMovie(ctx context.Context, id int) error {
	query := `UPDATE movies 
	SET is_deleted = TRUE 
	WHERE id = $1 AND is_deleted = FALSE`
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

func (mr *MoviesRepository) EditMovie(rctx context.Context, Image, Title, Duration, Synopsis string, id int) (models.Movie, error) {
	sql := `UPDATE movies SET image=$1, title=$2, duration=$3, synopsis=$4 WHERE id=$5 RETURNING id, image, title, duration, synopsis`
	values := []any{Image, Title, Duration, Synopsis, id}
	var movie models.Movie
	err := mr.db.QueryRow(rctx, sql, values...).Scan(&movie.Id, &movie.Image, &movie.Title, &movie.Duration, &movie.Synopsis)
	if err != nil {
		log.Println("Internal server error.\nCause: ", err.Error())
		return models.Movie{}, err
	}
	return movie, nil
}

func (mr *MoviesRepository) CreateMovie(rctx context.Context, body models.MovieCreate) (models.MovieCreate, error) {
	tx, err := mr.db.Begin(rctx)
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		return models.MovieCreate{}, err
	}

	defer tx.Rollback(rctx)
 // format ke string tanggal saja

	// Insert ke tabel movies
	sql := `INSERT INTO movies (title, release_date, duration, synopsis, id_director, rating)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, title, release_date, duration, synopsis, id_director, rating`
	values := []any{body.Title, body.ReleaseDate, body.Duration, body.Synopsis, body.Director, body.Rating}

	var newMovie models.MovieCreate
	if err := tx.QueryRow(rctx, sql, values...).Scan(
		&newMovie.Id,
		&newMovie.Title,
		&newMovie.ReleaseDate,
		&newMovie.Duration,
		&newMovie.Synopsis,
		&newMovie.Director,
		&newMovie.Rating,
	); err != nil {
		log.Println("Failed to insert movie:", err)
		return models.MovieCreate{}, err
	}

	// Insert ke tabel movies_actor
	for _, actorID := range body.ActorIDs {
		actorSQL := `INSERT INTO movies_actor (id_movies, id_actor) VALUES ($1, $2)`
		if _, err := tx.Exec(rctx, actorSQL, newMovie.Id, actorID); err != nil {
			log.Println("Failed to insert actor relation:", err)
			return models.MovieCreate{}, err
		}
	}

	// Insert ke tabel movies_genre
	for _, genreID := range body.GenreIDs {
		genreSQL := `INSERT INTO movies_genre (id_movies, id_genre) VALUES ($1, $2)`
		if _, err := tx.Exec(rctx, genreSQL, newMovie.Id, genreID); err != nil {
			log.Println("Failed to insert genre relation:", err)
			return models.MovieCreate{}, err
		}
	}

	// Commit transaksi
	if err := tx.Commit(rctx); err != nil {
		log.Println("Failed to commit transaction:", err)
		return models.MovieCreate{}, err
	}

	return newMovie, nil
}
