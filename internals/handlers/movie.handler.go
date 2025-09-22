package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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
		log.Println("Error GetDetailMovie:", err)
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
	limit := 20
	offset := (page - 1) * limit

	// Ambil filter jika ada
	title := ctx.Query("title") // kosong jika tidak ada
	// ambil genre dari query, pisah dengan koma
	genreParam := ctx.QueryArray("genre")
	genreMap := make(map[string]bool)
	var genre []string
	for _, g := range genreParam {
		name := strings.ToLower(strings.TrimSpace(g))
		if !genreMap[name] {
			genreMap[name] = true
			genre = append(genre, name)
		}
	}
	// var genre []string
	// if len(genreParam) >0 {
	// 	// genre = strings.Split(genreParam, ",")
	// 	for _, g := range (genreParam) {
	// 		genre = append(genre, strings.TrimSpace(g))
	// 	}
	// }
	// genre := genreParam

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
	// Ambil parameter movie ID dari URL
	MovieIDStr := ctx.Param("id")
	movieID, err := strconv.Atoi(MovieIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Movie ID tidak valid",
		})
		return
	}

	// Bind form data ke struct MovieBody
	var body models.MovieBody
	if err := ctx.ShouldBind(&body); err != nil {
		log.Println("Gagal bind data.\nSebab:", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Format data tidak valid",
		})
		return
	}

	// Set ID dari param ke body
	body.Id = movieID

	// Ambil claims JWT (jika perlu otentikasi)
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
	fmt.Println("User claims:", user)

	// Upload gambar Poster (Image)
	var imagePath *string
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
			log.Println("Gagal menyimpan file poster.\nSebab:", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Gagal menyimpan file poster",
			})
			return
		}
		imagePath = &generatedFilename
	}

	// Upload gambar Backdrop
	var backdropPath *string
	if body.Backdrop != nil {
		savePath, generatedFilename, err := utils.UploadImageFile(ctx, body.Backdrop, "public", fmt.Sprintf("user_%d", user.UserId))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		if err := ctx.SaveUploadedFile(body.Backdrop, savePath); err != nil {
			log.Println("Gagal menyimpan file backdrop.\nSebab:", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Gagal menyimpan file backdrop",
			})
			return
		}
		backdropPath = &generatedFilename
	}

	// Panggil repository untuk update data lengkap dengan transaction
	updatedMovie, err := mh.mr.EditMovie(ctx.Request.Context(), body, imagePath, backdropPath)
	if err != nil {
		log.Println("Gagal update movie.\nSebab:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Terjadi kesalahan saat menyimpan data",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updatedMovie,
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
	// var filebackdrop string
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
		log.Println(generatedFilename)
		// filebackdrop = generatedFilename
	}

	// Simpan ke database
	movieEntity := models.MovieBody{
		Title:       body.Title,
		ReleaseDate: body.ReleaseDate,
		Duration:    body.Duration,
		Synopsis:    body.Synopsis,
		Director:    body.Director,
		ActorIDs:    body.ActorIDs,
		GenreIDs:    body.GenreIDs,
		Rating:      body.Rating,
		// PosterPath:   fileImage,
		// BackdropPath: filebackdrop,
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

func (mh *movieHandler) GetMovieAdmin(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 1
	}
	limit := 20
	offset := (page - 1) * limit

	movies, err := mh.mr.GetMovieAdmin(ctx.Request.Context(), limit, offset)
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

func (h *movieHandler) GetMoviesByGenres(w *gin.Context) {
	genreParam := w.Query("genres")
	if genreParam == "" {
		w.JSON(http.StatusBadRequest, gin.H{
			"succes": false,
			"error":  "No genre IDs provided",
		})
		return
	}

	genreIDs := strings.Split(genreParam, ",")
	// ctx := r.Context()

	ctx := w.Request.Context()

	movies, err := h.mr.GetMoviesByAllGenres(ctx, genreIDs)
	if err != nil {
		w.JSON(http.StatusInternalServerError, gin.H{
			"succes":  false,
			"message": "Failed to get movies",
			"error":   err.Error(),
		})
		fmt.Println("Genre IDs:", genreIDs)
		return
	}
	w.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    movies,
	})
}

func (h *movieHandler) GetAllGenres(c *gin.Context) {
	ctx := c.Request.Context()

	genres, err := h.mr.GetAllGenres(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, genres)
}
