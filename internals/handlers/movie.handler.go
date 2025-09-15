package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/federus1105/weekly/internals/models"
	"github.com/federus1105/weekly/internals/repositories"
	"github.com/federus1105/weekly/internals/utils"
	"github.com/federus1105/weekly/pkg"
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

func (mh *movieHandler) GetAllMovie(ctx *gin.Context) {
	// Ambil query parameter
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 1
	}
	limit := 10
	offset := (page - 1) * limit

	// Ambil filter jika ada
	title := ctx.Query("title") // kosong jika tidak ada
	genre := ctx.Query("genre") // kosong jika tidak ada

	// Ambil data dari repository
	movies, err := mh.mr.GetAllOrFilteredMovies(ctx.Request.Context(), title, genre, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data movie",
			"error":   err.Error(),
		})
		return
	}

	// Jika tidak ada data ditemukan
	if len(movies) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    []string{},
			"message": "Tidak ada data movie ditemukan",
		})
		return
	}

	// Sukses
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    movies,
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
	var filename string
	if body.Image != nil {
		savePath, generatedFilename, err := utils.UploadImageFile(ctx, body.Image, "public", fmt.Sprintf("user_%d", user.UserId))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		if err := ctx.SaveUploadedFile(body.Image, savePath); err != nil {
			log.Println("Gagal menyimpan file.\nSebab:", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Gagal menyimpan file gambar",
			})
			return
		}

		filename = generatedFilename
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

func (mh *movieHandler) CreateMovie(ctx *gin.Context) {
	var body models.MovieBody
	fmt.Println("Content-Type:", ctx.ContentType())
	fmt.Println("release_date raw:", ctx.PostForm("release_date"))

	// Ambil form data
	if err := ctx.ShouldBind(&body); err != nil {
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
	var fileImage string
	if body.Image != nil {
		savePath, generatedFilename, err := utils.UploadImageFile(ctx, body.Image, "public", fmt.Sprintf("poster_path%d", user.UserId))
		if err != nil {
			log.Println("Upload poster gagal:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		if err := ctx.SaveUploadedFile(body.Image, savePath); err != nil {
			log.Println("Gagal menyimpan file.\nSebab:", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Gagal menyimpan file gambar",
			})
			return
		}
		fileImage = generatedFilename
	}
	fmt.Println(fileImage)

	// Upload Backdrop
	var filebackdrop string
	if body.Backdrop != nil {
		savePath, generatedFilename, err := utils.UploadImageFile(ctx, body.Backdrop, "public", fmt.Sprintf("backdrop_path%d", user.UserId))
		if err != nil {
			log.Println("Upload backdrop gagal:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		if err := ctx.SaveUploadedFile(body.Backdrop, savePath); err != nil {
			log.Println("Gagal menyimpan file.\nSebab:", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Gagal menyimpan file gambar",
			})
			return
		}

		filebackdrop = generatedFilename
	}

	// Simpan ke database
	movieEntity := models.MovieBody{
		Title:        body.Title,
		ReleaseDate:  body.ReleaseDate,
		Duration:     body.Duration,
		Synopsis:     body.Synopsis,
		Director:     body.Director,
		ActorIDs:     body.ActorIDs,
		GenreIDs:     body.GenreIDs,
		Rating:       body.Rating,
		PosterPath:   fileImage,
		BackdropPath: filebackdrop,
	}

	movie, err := mh.mr.CreateMovie(ctx.Request.Context(), movieEntity)
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
