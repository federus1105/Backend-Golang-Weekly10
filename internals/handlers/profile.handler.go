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
)

type ProfileHandler struct {
	pr *repositories.ProfileRepository
}

func NewProfileHandler(pr *repositories.ProfileRepository) *ProfileHandler {
	return &ProfileHandler{pr: pr}
}

// GetProfile godoc
// @Summary Get Profile
// @Tags Profile
// @Produce json
// @Param id path int true "ID Profile"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /profile/{id} [get]
func (ph *ProfileHandler) GetProfile(ctx *gin.Context) {
	profileIDStr := ctx.Param("id")
	ProfileID, err := strconv.Atoi(profileIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Profile ID tidak valid",
		})
		return
	}
	profiles, err := ph.pr.GetProfile(ctx.Request.Context(), ProfileID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data Profile",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    profiles,
	})

}

func (s *ProfileHandler) EditProfile(ctx *gin.Context) {
	profileIDStr := ctx.Param("id")
	profileID, err := strconv.Atoi(profileIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Profile ID tidak valid",
		})
		return
	}

	// Ambil data dari form
	var body models.ProfileBody
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
	profile, err := s.pr.EditProfile(
		ctx.Request.Context(),
		filename,
		body.FirstName,
		body.LastName,
		body.Phone,
		profileID)
	if err != nil {
		log.Println("Gagal update profil.\nSebab:", err.Error())
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
