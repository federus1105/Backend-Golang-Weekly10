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
			"error":   "internal server error",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   jwtToken,
	})

}

func (a *AuthHandler) CreateUser(ctx *gin.Context) {
	var body models.UserRegister
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}
	// You can add further logic to process the order here
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"order":   body,
	})
}

func (a *AuthHandler) MigrateHashPasswords(ctx *gin.Context) {
	log.Println("[DEBUG] Mulai migrasi password")

	users, err := a.ar.GetAllUsers(ctx.Request.Context())
	if err != nil {
		log.Println("[ERROR] Gagal GetAllUsers:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Gagal mengambil data user",
		})
		return
	}

	log.Printf("[DEBUG] Jumlah user: %d", len(users))

	hc := pkg.NewHashConfig()
	hc.UseRecommended()

	var successCount int
	var failed []string

	for _, user := range users {
		log.Printf("[DEBUG] Memproses user: %s", user.Email)

		if user.Password == "" {
			log.Printf("[WARN] Password kosong untuk user: %s", user.Email)
			failed = append(failed, user.Email)
			continue
		}

		// Gunakan GenHash untuk format yang sesuai
		hashed, err := hc.GenHash(user.Password)
		if err != nil {
			log.Printf("[ERROR] Hash gagal untuk %s: %v", user.Email, err)
			failed = append(failed, user.Email)
			continue
		}

		err = a.ar.UpdateUserPassword(ctx.Request.Context(), uint(user.Id), hashed)
		if err != nil {
			log.Printf("[ERROR] Update gagal untuk %s: %v", user.Email, err)
			failed = append(failed, user.Email)
			continue
		}

		successCount++
	}

	log.Println("[DEBUG] Migrasi selesai")
	ctx.JSON(http.StatusOK, gin.H{
		"success":       true,
		"message":       "Migrasi password selesai",
		"success_count": successCount,
		"failed_users":  failed,
	})
}
