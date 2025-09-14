package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/federus1105/weekly/internals/models"
	"github.com/federus1105/weekly/internals/repositories"
	"github.com/federus1105/weekly/pkg"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	ar *repositories.AuthRepository
}

func NewAuthHandler(ar *repositories.AuthRepository) *AuthHandler {
	return &AuthHandler{ar: ar}
}

// Login godoc
// @Summary Login
// @Tags Authentication
// @Accept json
// @Produce json
// @Param order body models.UserAuth true "Login"
// @Success 201 {object} models.UserAuth
// @Router /auth/login [post]
func (a *AuthHandler) Login(ctx *gin.Context) {
	// menerima body
	var body models.UserAuth
	if err := ctx.ShouldBind(&body); err != nil {
		if strings.Contains(err.Error(), "required") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Nama dan Password harus diisi",
			})
			return
		}

		if strings.Contains(err.Error(), "min") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Password minimum 4 karakter",
			})
			return
		}
		log.Println("Internal Server Error.\nCause: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "internal server error",
		})
		return
	}
	// ambil data user
	user, err := a.ar.GetUserWithPasswordAndRole(ctx.Request.Context(), body.Email)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Nama atau Password salah",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "internal server errorr",
		})
		return
	}

	// bandingkan password
	hc := pkg.NewHashConfig()
	isMatched, err := hc.CompareHashAndPassword(body.Password, user.Password)
	if err != nil {
		log.Println("Internal Server Error.\nCause: ", err)
		re := regexp.MustCompile("hash|crypto|argon2id|format")
		if re.Match([]byte(err.Error())) {
			log.Println("Error during Hashing")
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "internal server error",
		})
		return
	}
	if !isMatched {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Nama atau Password salah",
		})
		return
	}
	// jika match, maka buatkan jwt dan kirim via response
	claims := pkg.NewJWTClaims(user.Id, user.Role)
	jwtToken, err := claims.GenToken()
	if err != nil {
		log.Println("Internal Server Error.\nCause: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "internal server errorrr",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   jwtToken,
	})

}

// Register godoc
// @Summary Register
// @Tags Authentication
// @Accept json
// @Produce json
// @Param order body models.UserRegister true "Register"
// @Success 201 {object} models.UserRegister
// @Router /auth/register [post]
func (a *AuthHandler) Register(ctx *gin.Context) {
	var body models.UserRegister
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}
	newOrder, err := a.ar.Register(ctx.Request.Context(), body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	// You can add further logic to process the order here
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"order":   newOrder,
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var body models.ChangePassword
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDInterface.(int) // pastikan tipe sesuai
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	err := h.ar.ResetPassword(c.Request.Context(), userID, body.OldPassword, body.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil direset"})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Hanya memberitahu user untuk hapus token dari client
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout berhasil. Silakan hapus token di sisi client.",
	})
}
