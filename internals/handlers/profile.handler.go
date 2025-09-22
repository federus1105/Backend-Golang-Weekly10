package handlers

import (
	"fmt"
	"net/http"

	"github.com/federus1105/weekly/internals/models"
	"github.com/federus1105/weekly/internals/repositories"
	"github.com/federus1105/weekly/internals/utils"
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
	// Ambil user_id dari Gin Context
	userIDRaw, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized: user ID tidak ditemukan di context",
		})
		return
	}

	profileID, ok := userIDRaw.(int)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "User dalam context tidak valid",
		})
		return
	}
	profiles, err := ph.pr.GetProfile(ctx.Request.Context(), profileID)
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
	userIDRaw, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	userID, ok := userIDRaw.(int)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Invalid user ID"})
		return
	}

	var body models.ProfileBody
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Format data tidak valid"})
		return
	}

	var filename *string
	if body.Image != nil {
		savePath, generatedFilename, err := utils.UploadImageFile(ctx, body.Image, "public", fmt.Sprintf("user_%d", userID))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
			return
		}

		if err := ctx.SaveUploadedFile(body.Image, savePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Gagal menyimpan file gambar"})
			return
		}

		filename = &generatedFilename
	}

	// Panggil fungsi PATCH di repository
	profile, err := s.pr.EditProfile(
		ctx.Request.Context(),
		filename,
		body.FirstName,
		body.LastName,
		body.Phone,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": profile})
}
