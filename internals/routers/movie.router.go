package routers

import (
	"github.com/federus1105/weekly/internals/handlers"
	"github.com/federus1105/weekly/internals/middlewares"
	"github.com/federus1105/weekly/internals/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitMoviesRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	movieRouter := router.Group("/movies")
	sr := repositories.NewMoviesRepository(db, rdb)
	sh := handlers.NewMovieHandler(sr)

	movieRouter.GET("/genres/list", sh.GetAllGenres)
	// movieRouter.GET("/genres", sh.GetMoviesByGenres)
	movieRouter.GET("/admin", sh.GetMovieAdmin)
	movieRouter.GET("/", sh.GetAllMovie)
	movieRouter.GET("/upcoming", sh.GetUpcomingMovies)
	movieRouter.GET("/popular", sh.GetPopularMovies)
	movieRouter.GET("/:id", middlewares.VerifyToken, middlewares.Access("User", "Admin"), sh.GetDetailMovie)
	movieRouter.GET("/allmovie", middlewares.VerifyToken, middlewares.Access("Admin"), sh.GetAllMovie)
	movieRouter.DELETE("/:movie_id", middlewares.VerifyToken, middlewares.Access("Admin"), sh.DeleteMovie)
	movieRouter.PATCH("/:id", middlewares.VerifyToken, middlewares.Access("Admin"), sh.EditMovie)
	// movieRouter.POST("/create", middlewares.VerifyToken, middlewares.Access("Admin"), sh.CreateMovie)
}
