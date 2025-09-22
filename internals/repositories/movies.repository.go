package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/federus1105/weekly/internals/models"
	"github.com/jackc/pgx/v5"
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
	start := time.Now()
	redisKey := "firdaus:popular-movies"
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
	// ambil data movie
	sql := `SELECT 
    m.id,
    m.image,
    m.title,
    m.release_date,
    ARRAY_AGG(g.name) AS genres
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

func (mr *MoviesRepository) GetPopularMovies(rctx context.Context, limit, offset int) ([]models.Movie, error) {
	start := time.Now()
	redisKey := "firdaus:popular-movies"
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
	sql := `
SELECT 
  m.id,
  m.image,
  m.title,
  ARRAY_AGG(g.name, ', ') AS genres
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
		if err := rows.Scan(&movie.Id, &movie.Image, &movie.Title, &movie.Genres); err != nil {
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
	ARRAY_AGG(DISTINCT g.name) AS genres,
	ARRAY_AGG(DISTINCT a.name) AS actor
	FROM movies m
	JOIN movies_genre mg ON m.id = mg.id_movies
	JOIN genres g ON mg.id_genre = g.id
	JOIN movies_actor ma ON m.id = ma.id_movie
	JOIN actor a ON ma.id_actor = a.id
	JOIN director d ON m.id_director = d.id
	WHERE m.id = $1
	AND m.is_deleted = false
	GROUP BY 
	m.id, 
	m.image, 
	m.backdrop, 
	m.title, 
	m.release_date, 
	m.duration, 
	m.synopsis
	`
	rows, err := mr.db.Query(rctx, sql, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.Id, &movie.Image, &movie.Backdrop, &movie.Title, &movie.ReleaseDate, &movie.Duration, &movie.Director, &movie.Synopsis, &movie.Genres, &movie.Actor); err != nil {
			log.Println("Error saat scan rows:", err)
			return nil, err
		}
		movies = append(movies, movie)
	}
	return movies, nil
}

func (mr *MoviesRepository) GetAllOrFilteredMovies(rctx context.Context, title string, genre []string, limit, offset int) ([]models.Movie, error) {
	start := time.Now()
	isFiltering := title != "" || len(genre) > 0

	// Redis cache hanya tanpa filter dan pagination
	redisKey := "firdaus:allmovies"
	if !isFiltering && offset == 0 {
		cmd := mr.rdb.Get(rctx, redisKey)
		if cmd.Err() == nil {
			// Cache hit
			var cachedMovies []models.Movie
			cmdBytes, err := cmd.Bytes()
			if err == nil {
				if err := json.Unmarshal(cmdBytes, &cachedMovies); err == nil {
					log.Printf("Key %s found in cache ✅", redisKey)
					log.Printf("Served in %s using Redis", time.Since(start))
					return cachedMovies, nil
				}
			}
		}
	}

	baseSQL := `
	SELECT
    m.id,
    m.title,
    m.image,
    COALESCE(ARRAY_AGG(DISTINCT g.name) FILTER (WHERE g.name IS NOT NULL)) AS genres
FROM movies m
INNER JOIN movies_genre mg ON m.id = mg.id_movies
INNER JOIN genres g ON mg.id_genre = g.id`

	// Dynamic conditions
	var conditions []string
	var args []interface{}
	argIdx := 1

	// filter title
	if title != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(m.title) ILIKE LOWER($%d)", argIdx))
		args = append(args, "%"+title+"%")
		argIdx++
	}

	if len(conditions) > 0 {
		baseSQL += " AND " + strings.Join(conditions, " AND ")
	}

	// Tambahkan join dengan unnest genre jika ada
	if len(genre) > 0 {
		baseSQL += `
	LEFT JOIN (
		SELECT unnest($` + fmt.Sprint(argIdx) + `::text[]) AS genre_name
	) AS selected_genres ON LOWER(g.name) = LOWER(selected_genres.genre_name)
	`
		args = append(args, genre)
		argIdx++
	}
	// Akhir query
	baseSQL += `
GROUP BY m.id, m.title, m.image`

	if len(genre) > 0 {
		baseSQL += fmt.Sprintf(`
HAVING COUNT(DISTINCT selected_genres.genre_name) = $%d
`, argIdx)
		args = append(args, len(genre))
		argIdx++
	}

	// Order, limit, offset
	baseSQL += fmt.Sprintf(`
ORDER BY m.title ASC
LIMIT $%d OFFSET $%d
`, argIdx, argIdx+1)
	args = append(args, limit, offset)

	// Query
	rows, err := mr.db.Query(rctx, baseSQL, args...)
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

	// Cache jika tidak filter
	if !isFiltering && offset == 0 {
		bt, err := json.Marshal(movies)
		if err == nil {
			if err := mr.rdb.Set(rctx, redisKey, string(bt), 1*time.Minute).Err(); err != nil {
				log.Println("Redis Set Error:", err.Error())
			}
		}
	}

	log.Printf("[MOVIES] Served in %s using %s", time.Since(start), func() string {
		if !isFiltering && offset == 0 {
			return "Redis or DB (cached)"
		}
		return "DB (filtered)"
	}())

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

func (r *MoviesRepository) EditMovie(ctx context.Context, body models.MovieBody, image *string,
	backdrop *string) (models.Movie, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return models.Movie{}, err
	}

	defer tx.Rollback(ctx)

	setClauses := []string{}
	args := []any{}
	argID := 1

	if body.Image != nil {
		// Simpan file dan dapatkan pathnya dulu di layer service/controller
		setClauses = append(setClauses, fmt.Sprintf("image = $%d", argID))
		args = append(args, *image)
		argID++
	}
	if body.Backdrop != nil {
		setClauses = append(setClauses, fmt.Sprintf("backdrop = $%d", argID))
		args = append(args, *backdrop)
		argID++
	}
	if body.Title != "" {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", argID))
		args = append(args, body.Title)
		argID++
	}
	if !body.ReleaseDate.IsZero() {
		setClauses = append(setClauses, fmt.Sprintf("release_date = $%d", argID))
		args = append(args, body.ReleaseDate)
		argID++
	}
	if body.Duration != "" {
		setClauses = append(setClauses, fmt.Sprintf("duration = $%d", argID))
		args = append(args, body.Duration)
		argID++
	}
	if body.Director != 0 {
		setClauses = append(setClauses, fmt.Sprintf("id_director = $%d", argID))
		args = append(args, body.Director)
		argID++
	}
	if body.Synopsis != "" {
		setClauses = append(setClauses, fmt.Sprintf("synopsis = $%d", argID))
		args = append(args, body.Synopsis)
		argID++
	}
	if body.Rating != 0 {
		setClauses = append(setClauses, fmt.Sprintf("rating = $%d", argID))
		args = append(args, body.Rating)
		argID++
	}

	if len(setClauses) == 0 {
		return models.Movie{}, fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf(`
        UPDATE movies
        SET %s
        WHERE id = $%d
        RETURNING id, image, backdrop, title, release_date, duration, id_director, synopsis, rating
    `, strings.Join(setClauses, ", "), argID)

	args = append(args, body.Id)

	var movie models.Movie
	err = tx.QueryRow(ctx, query, args...).Scan(
		&movie.Id,
		&movie.Image,
		&movie.Backdrop,
		&movie.Title,
		&movie.ReleaseDate,
		&movie.Duration,
		&movie.Director,
		&movie.Synopsis,
		&movie.Rating,
	)
	if err != nil {
		tx.Rollback(ctx)
		return models.Movie{}, err
	}

	// Update movies_actor relasi
	// _, err = tx.Exec(ctx, "DELETE FROM movies_actor WHERE id_movie = $1", body.Id)
	// if err != nil {
	// 	tx.Rollback(ctx)
	// 	return models.Movie{}, err
	// }
	for _, actorID := range body.ActorIDs {
		_, err = tx.Exec(ctx, "INSERT INTO movies_actor (id_movie, id_actor) VALUES ($1, $2)", body.Id, actorID)
		if err != nil {
			tx.Rollback(ctx)
			return models.Movie{}, err
		}
	}

	// Update movies_genre relasi
	// _, err = tx.Exec(ctx, "DELETE FROM movies_genre WHERE id_movies = $1", body.Id)
	// if err != nil {
	// 	tx.Rollback(ctx)
	// 	return models.Movie{}, err
	// }
	for _, genreID := range body.GenreIDs {
		_, err = tx.Exec(ctx, "INSERT INTO movies_genre (id_movies, id_genre) VALUES ($1, $2)", body.Id, genreID)
		if err != nil {
			tx.Rollback(ctx)
			return models.Movie{}, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return models.Movie{}, err
	}

	return movie, nil
}

func (mr *MoviesRepository) CreateMovie(rctx context.Context, body models.MovieBody) (models.MovieBody, error) {
	tx, err := mr.db.Begin(rctx)
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		return models.MovieBody{}, err
	}

	defer tx.Rollback(rctx)
	// Insert ke tabel movies
	sql := `INSERT INTO movies (title, release_date, duration, synopsis, id_director, rating, image, backdrop)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, title, release_date, duration, synopsis, id_director, rating, image, backdrop`
	values := []any{body.Title, body.ReleaseDate, body.Duration, body.Synopsis, body.Director, body.Rating, body.Image,
		body.Backdrop}
	var newMovie models.MovieBody
	if err := tx.QueryRow(rctx, sql, values...).Scan(
		&newMovie.Id,
		&newMovie.Title,
		&newMovie.ReleaseDate,
		&newMovie.Duration,
		&newMovie.Synopsis,
		&newMovie.Director,
		&newMovie.Rating,
		&newMovie.Image,
		&newMovie.Backdrop,
	); err != nil {
		log.Println("Failed to insert movie:", err)
		return models.MovieBody{}, err
	}

	// Insert ke tabel movies_actor
	for _, actorID := range body.ActorIDs {
		actorSQL := `INSERT INTO movies_actor (id_movie, id_actor) VALUES ($1, $2)`
		if _, err := tx.Exec(rctx, actorSQL, newMovie.Id, actorID); err != nil {
			log.Println("Failed to insert actor relation:", err)
			return models.MovieBody{}, err
		}
	}

	// Insert ke tabel movies_genre
	for _, genreID := range body.GenreIDs {
		genreSQL := `INSERT INTO movies_genre (id_movies, id_genre) VALUES ($1, $2)`
		if _, err := tx.Exec(rctx, genreSQL, newMovie.Id, genreID); err != nil {
			log.Println("Failed to insert genre relation:", err)
			return models.MovieBody{}, err
		}

	}

	// Commit transaksi
	if err := tx.Commit(rctx); err != nil {
		log.Println("Failed to commit transaction:", err)
		return models.MovieBody{}, err
	}

	return newMovie, nil
}

func (mr *MoviesRepository) GetMovieAdmin(rctx context.Context, limit, offset int) ([]models.MovieAdmin, error) {
	sql := `SELECT 
		m.id,
		m.image,
		m.title,
		m.release_date,
		m.duration,
		STRING_AGG(DISTINCT g.name, ', ') AS genres
		FROM movies m
		JOIN movies_genre mg ON m.id = mg.id_movies
		JOIN genres g ON mg.id_genre = g.id
		WHERE m.is_deleted = false
		GROUP BY 
			m.id, 
			m.image, 
			m.title, 
			m.release_date, 
			m.duration
		ORDER BY m.id
		LIMIT $1 OFFSET $2;`
	rows, err := mr.db.Query(rctx, sql, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var movies []models.MovieAdmin
	for rows.Next() {
		var movie models.MovieAdmin
		if err := rows.Scan(&movie.Id, &movie.Image, &movie.Title, &movie.ReleaseDate, &movie.Duration, &movie.Genres); err != nil {
			log.Println("Error saat scan rows", err)
			return nil, err
		}
		movies = append(movies, movie)
	}
	return movies, nil
}

func (sr *MoviesRepository) GetMoviesByAllGenres(ctx context.Context, genreIDs []string) ([]models.Movie, error) {
	if len(genreIDs) == 0 {
		return nil, errors.New("genreIDs is empty")
	}
	placeholders := make([]string, len(genreIDs))
	args := make([]interface{}, len(genreIDs))
	for i, id := range genreIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	query := fmt.Sprintf(`
    SELECT 
    m.id,
    m.title,
    ARRAY_AGG(DISTINCT g.name) AS genres
FROM 
    movies m
JOIN 
    movies_genre mg ON m.id = mg.id_movies
JOIN 
    genres g ON g.id = mg.id_genre
WHERE 
    mg.id_genre IN (%s)
GROUP BY 
    m.id
HAVING 
    COUNT(DISTINCT mg.id_genre) = %d;
    `, strings.Join(placeholders, ", "), len(genreIDs))

	rows, err := sr.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie

	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.Id, &movie.Title, &movie.Genres); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	return movies, nil
}

func (r *MoviesRepository) GetAllGenres(ctx context.Context) ([]models.Genre, error) {
	rows, err := r.db.Query(ctx, "SELECT id, name FROM genres")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []models.Genre
	for rows.Next() {
		var g models.Genre
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			return nil, err
		}
		genres = append(genres, g)
	}

	return genres, nil
}
