package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/federus1105/weekly/internals/models"
	"github.com/federus1105/weekly/internals/repositories"
	"github.com/federus1105/weekly/pkg"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
	title := ctx.Param("title")
	genre := ctx.Query("genre")

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 1
	}
	limit := 5
	offset := (page - 1) * limit
	movies, err := mh.mr.GetFilterMovie(ctx.Request.Context(), title, genre, limit, offset)
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
// @Security BearerAuth
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
// @Security BearerAuth
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

func (mh *movieHandler) EditMovie(ctx *gin.Context) {
	MovieIDStr := ctx.Param("id")
	movieID, err := strconv.Atoi(MovieIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Movie ID tidak valid",
		})
		return
	}

	// Ambil data dari form
	var body models.MovieBody
	if err := ctx.ShouldBind(&body); err != nil {
		log.Println("Gagal bind data.\nSebab:", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Format data tidak valid",
		})
		return
	}

	// Ambil claims dari JWT
	claims, isExist := ctx.Get("claims")
	if !isExist {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Silakan login kembali",
		})
		return
	}
	user, ok := claims.(pkg.Claims)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Terjadi kesalahan internal",
		})
		return
	}

	fmt.Println(user)
	// Upload gambar
	file := body.Image
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d_images_%d%s", time.Now().UnixNano(), user.UserId, ext)
	location := filepath.Join("public", filename)

	if err := ctx.SaveUploadedFile(file, location); err != nil {
		log.Println("Gagal upload gambar.\nSebab:", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Gagal upload gambar",
		})
		return
	}

	// Simpan ke database
	profile, err := mh.mr.EditMovie(
		ctx.Request.Context(),
		filename,
		body.Title,
		body.Duration,
		body.Synopsis,
		movieID)
	if err != nil {
		log.Println("Gagal update Movie.\nSebab:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Terjadi kesalahan saat menyimpan data",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    profile,
	})
}

// func (mh *movieHandler) CreateMovie(ctx *gin.Context) {
// 	var body models.MovieCreate
// 	if err := ctx.ShouldBind(&body); err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error":   err.Error(),
// 			"success": false,
// 		})
// 		return
// 	}
// 	newMovie, err := mh.mr.CreateMovie(ctx.Request.Context(), body)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"success": false,
// 			"error":   err.Error(),
// 		})
// 		return
// 	}
// 	ctx.JSON(http.StatusCreated, gin.H{
// 		"success": true,
// 		"data":    newMovie,
// 	})
// }

func (mh *movieHandler) CreateMovie(ctx *gin.Context) {
	var body models.MovieBody
	fmt.Println("Content-Type:", ctx.ContentType())
	fmt.Println("release_date raw:", ctx.PostForm("release_date"))

	// Ambil form data
	if err := ctx.ShouldBindWith(&body, binding.FormMultipart); err != nil {
		log.Println("Gagal bind data:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,

			"error": "Format data tidak valid",
		})
		return
	}

	// Ambil JWT claims
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Silakan login kembali",
		})
		return
	}

	user, ok := claims.(pkg.Claims)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Terjadi kesalahan internal",
		})
		return
	}

	// Upload Poster (Image)
	file := body.Image
	if file != nil {
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%d_image_%d%s", time.Now().UnixNano(), user.UserId, ext)
		path := filepath.Join("public", filename)

		if err := ctx.SaveUploadedFile(file, path); err != nil {
			log.Println("Gagal upload poster:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Gagal upload poster",
			})
			return
		}

		body.PosterPath = filename
	}

	// Upload Backdrop
	backdrop := body.Backdrop
	if backdrop != nil {
		ext := filepath.Ext(backdrop.Filename)
		filename := fmt.Sprintf("%d_backdrop_%d%s", time.Now().UnixNano(), user.UserId, ext)
		path := filepath.Join("public", filename)

		if err := ctx.SaveUploadedFile(backdrop, path); err != nil {
			log.Println("Gagal upload backdrop:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Gagal upload backdrop",
			})
			return
		}

		body.BackdropPath = filename
	}

	// Simpan ke database
	movie, err := mh.mr.CreateMovie(ctx.Request.Context(), body)
	if err != nil {
		log.Println("Gagal simpan movie:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Terjadi kesalahan saat menyimpan data",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    movie,
	})
}

//  Bangun model movie untuk disimpan
// movie := models.MovieBody{
// 	Title:       body.Title,
// 	ReleaseDate: body.ReleaseDate,
// 	Duration:    body.Duration,
// 	Synopsis:    body.Synopsis,
// 	Director:    body.Director,
// 	ActorIDs:    body.ActorIDs,
// 	GenreIDs:    body.GenreIDs,
// 	Rating:      body.Rating,
// 	Image:       filename,
// 	Backdrop:    filebackdrop,
// }
// Simpan ke DB
// 	result, err := mh.mr.CreateMovie(ctx, movie)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"success": false,
// 			"message": "Gagal membuat movie",
// 			"error":   err.Error(),
// 		})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"data":    result,
// 	})
// }

// func (h *movieHandler) CreateMovie(ctx *gin.Context) {
// 	// Bind form-data biasa (bukan JSON)
// 	var body models.MovieCreate
// 	if err := ctx.ShouldBind(&body); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "Data form tidak valid",
// 			"error":   err.Error(),
// 		})
// 		return
// 	}

// 	// Ambil file gambar (poster)
// 	imageFile, err := ctx.FormFile("image")
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "Gambar poster tidak ditemukan",
// 		})
// 		return
// 	}

// 	// Ambil file backdrop
// 	backdropFile, err := ctx.FormFile("backdrop")
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "Gambar backdrop tidak ditemukan",
// 		})
// 		return
// 	}

// 	// Simpan file image
// 	imageFilename, err := utils.SaveUploadedImage(ctx, imageFile, "image", 0)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"success": false,
// 			"message": "Gagal menyimpan gambar",
// 			"error":   err.Error(),
// 		})
// 		return
// 	}

// 	// Simpan file backdrop
// 	backdropFilename, err := utils.SaveUploadedImage(ctx, backdropFile, "backdrop", 0)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"success": false,
// 			"message": "Gagal menyimpan backdrop",
// 			"error":   err.Error(),
// 		})
// 		return
// 	}

// 	// Bangun model movie untuk disimpan
// 	movie := models.MovieCreate{
// 		Title:       body.Title,
// 		ReleaseDate: body.ReleaseDate,
// 		Duration:    body.Duration,
// 		Synopsis:    body.Synopsis,
// 		Director:    body.Director,
// 		ActorIDs:    body.ActorIDs,
// 		GenreIDs:    body.GenreIDs,
// 		Rating:      body.Rating,
// 		Image:       imageFilename,
// 		Backdrop:    backdropFilename,
// 	}

// 	// Simpan ke DB
// 	result, err := h.mr.CreateMovie(ctx, movie)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"success": false,
// 			"message": "Gagal membuat movie",
// 			"error":   err.Error(),
// 		})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"data":    result,
// 	})
// }
