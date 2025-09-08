package routers

import (
	"github.com/federus1105/weekly/internals/handlers"
	"github.com/federus1105/weekly/internals/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitMoviesRouter(router *gin.Engine, db *pgxpool.Pool) {
	movieRouter := router.Group("/movies")
	sr := repositories.NewMoviesRepository(db)
	sh := handlers.NewMovieHandler(sr)

	movieRouter.GET("/", sh.GetUpcomingMovies)
	movieRouter.GET("/upcoming", sh.GetUpcomingMovies)
	movieRouter.GET("/popular", sh.GetPopularMovies)
	movieRouter.GET("/filter", sh.GetFilterMovie)
	movieRouter.GET("/:id", sh.GetDetailMovie)
	movieRouter.GET("/allmovie", sh.GetAllMovie)
	movieRouter.DELETE("/:movie_id", sh.DeleteMovie)
}
