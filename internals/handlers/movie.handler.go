package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/federus1105/weekly/internals/repositories"
	"github.com/gin-gonic/gin"
)

type movieHandler struct {
	mr *repositories.MoviesRepository
}

func NewMovieHandler(mr *repositories.MoviesRepository) *movieHandler {
	return &movieHandler{mr: mr}
}

// GetUpcomingMovies godoc
// @Summary Get upcoming movies
// @Tags Movies
// @Produce json
// @Param page query int false "Page"
// @Success 200 {object} map[string]interface{}
// @Router /movies/upcoming [get]
func (mh *movieHandler) GetUpcomingMovies(ctx *gin.Context) {
	// pagination
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 1
	}
	limit := 5
	offset := (page - 1) * limit

	movies, err := mh.mr.GetUpcomingMovies(ctx.Request.Context(), limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"succes": false,
			"data":   movies,
		})
		return
	}
	if len(movies) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    []string{},
			"message": "Tidak ada data movie",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"succes": true,
		"data":   movies,
	})
}

// GetPopularMovies godoc
// @Summary Get Popular movies
// @Tags Movies
// @Produce json
// @Param page query int false "Page"
// @Success 200 {object} map[string]interface{}
// @Router /movies/popular [get]
func (mh *movieHandler) GetPopularMovies(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))

	if err != nil {
		page = 1
	}
	limit := 5
	offset := (page - 1) * limit
	movies, err := mh.mr.GetPopularMovies(ctx.Request.Context(), limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"data":    movies,
		})
		return
	}
	if len(movies) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    []string{},
			"message": "Tidak ada data movie",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"succes": true,
		"data":   movies,
	})
}

// GetFilterMovies godoc
// @Summary Get Filter movies
// @Tags Movies
// @Produce json
// @Param page query int false "Page"
// @Success 200 {object} map[string]interface{}
// @Router /movies/filter [get]
func (mh *movieHandler) GetFilterMovie(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))

	if err != nil {
		page = 1
	}
	limit := 5
	offset := (page - 1) * limit
	movies, err := mh.mr.GetFilterMovie(ctx.Request.Context(), limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"data":    movies,
		})
		return
	}
	if len(movies) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    []string{},
			"message": "Tidak ada data movie",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"succes": true,
		"data":   movies,
	})

}

// GetDetailMovie godoc
// @Summary Get Detail Movie
// @Tags Movies
// @Produce json
// @Param id path int true "Movie Detail"
// @Success 200 {object} map[string]interface{}
// @Router /movies/{id} [get]
func (mh *movieHandler) GetDetailMovie(ctx *gin.Context) {
	movieIDStr := ctx.Param("id")
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID Movie tidak valid",
		})
		return
	}
	movies, err := mh.mr.GetDetailMovie(ctx.Request.Context(), movieID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Gagal mengambil data Movie",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    movies,
	})
}

// GetAllMovie godoc
// @Summary Get All Movie
// @Tags Movies
// @Produce json
// @Param page query int false "Page"
// @Success 200 {object} map[string]interface{}
// @Router /movies/allmovie [get]
func (mh *movieHandler) GetAllMovie(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 1
	}
	limit := 10
	offset := (page - 1) * limit
	movies, err := mh.mr.GetAllMovie(ctx.Request.Context(), limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"data":    movies,
		})
		return
	}
	if len(movies) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    []string{},
			"message": "Tidak ada data movie",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"succes": true,
		"data":   movies,
	})

}

func (mh *movieHandler) DeleteMovie(ctx *gin.Context) {
	// Ambil param ID
	movieIDStr := ctx.Param("movie_id")
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "movie id tidak valid",
		})
		return
	}

	// Panggil method dari repository
	err = mh.mr.DeleteMovie(ctx, movieID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Sukses
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("movie dengan id %d berhasil dihapus", movieID),
	})
}
