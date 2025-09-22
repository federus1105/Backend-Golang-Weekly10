package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/federus1105/weekly/internals/models"
	"github.com/federus1105/weekly/internals/repositories"
	"github.com/federus1105/weekly/internals/utils"
	"github.com/federus1105/weekly/pkg"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type AuthHandler struct {
	ar          *repositories.AuthRepository
	redisClient *redis.Client
}

func NewAuthHandler(ar *repositories.AuthRepository, rdb *redis.Client) *AuthHandler {
	return &AuthHandler{ar: ar, redisClient: rdb}
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

func (a *AuthHandler) ResetPassword(c *gin.Context) {
	var body models.ChangePassword
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// mengambil user yang lagi login sekarang
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	err := a.ar.ResetPassword(c.Request.Context(), userID, body.OldPassword, body.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil direset"})
}

func (a *AuthHandler) Logout(ctx *gin.Context) {
	token := utils.GetTokenFromHeader(ctx)
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Token tidak ditemukan"})
		return
	}

	expiry := utils.GetTokenExpiry(token)
	if expiry.IsZero() {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Token tidak valid atau tidak memiliki expiry"})
		return
	}

	duration := time.Until(expiry)
	if duration <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Token sudah kedaluwarsa"})
		return
	}

	// Simpan token ke Redis blacklist dengan TTL sesuai sisa masa berlaku token
	err := a.redisClient.Set(ctx, "blacklist:"+token, true, duration).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal logout"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Berhasil logout"})
}
