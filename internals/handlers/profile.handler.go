package handlers

import (
	"net/http"
	"strconv"

	"github.com/federus1105/weekly/internals/repositories"
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
// func (ph *ProfileHandler) EditProfile(ctx *gin.Context) {
// 	profileIDStr := ctx.Param("id")
// 	ProfileID, err := strconv.Atoi(profileIDStr)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "Profile ID tidak valid",
// 		})
// 		return
// 	}
// 	profiles, err := ph.pr.EditProfile(ctx.Request.Context(), ProfileID)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"success": false,
// 			"message": "Gagal mengambil data Profile",
// 		})
// 		return
// 	}
// 	profiles, err := ph.pr.EditProfile(ctx.Request.Context(), input.Image, input.FirstName, input.LastName, input.Phone, id)
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update profile")
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"data":    profiles,
// 	})
// }
