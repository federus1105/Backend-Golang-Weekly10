package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/federus1105/weekly/pkg"
	"github.com/gin-gonic/gin"
)

type ctxKey string

const UserIDKey ctxKey = "user_id"

// const roleKey ctxKey = "role"

// func getTokenFromHeader(c *gin.Context) string {
// 	authHeader := c.GetHeader("Authorization")
// 	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
// 		return ""
// 	}
// 	return strings.TrimPrefix(authHeader, "Bearer ")
// }

// func AuthMiddleware(rdb *redis.Client) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Ambil token dari header
// 		token := getTokenFromHeader(c)
// 		if token == "" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 				"success": false,
// 				"error":   "Authorization header missing or invalid",
// 			})
// 			return
// 		}

// 		// Cek apakah token sudah diblacklist
// 		val, err := rdb.Get(c, "blacklist:"+token).Result()
// 		if err == nil && val == "true" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 				"success": false,
// 				"error":   "Token sudah logout",
// 			})
// 			return
// 		} else if err != nil && err != redis.Nil {
// 			// Handle error selain key not found (redis.Nil)
// 			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
// 				"success": false,
// 				"error":   "Internal server error (Redis): " + err.Error(),
// 			})
// 			return
// 		}

// 		// Verifikasi dan parsing token
// 		claims := &pkg.Claims{}
// 		if err := claims.VerifyToken(token); err != nil {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 				"success": false,
// 				"error":   "Invalid token: " + err.Error(),
// 			})
// 			return
// 		}

// 		// Set ke context
// 		c.Set("user_id", claims.UserId)
// 		c.Set("role", claims.Role)

// 		ctx := context.WithValue(c.Request.Context(), UserIDKey, claims.UserId)
// 		c.Request = c.Request.WithContext(ctx)

// 		c.Next()
// 	}
// }

// fungsi untuk mengetahui user yang login sekarang
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Authorization header missing or invalid",
			})
			return
		}

		// Ambil token-nya
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Buat instance Claims
		claims := &pkg.Claims{}

		// Verifikasi token
		if err := claims.VerifyToken(tokenString); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid token: " + err.Error(),
			})
			return
		}

		// Simpan user_id dan role ke context
		c.Set("user_id", claims.UserId)
		c.Set("role", claims.Role)

		ctx := context.WithValue(c.Request.Context(), UserIDKey, claims.UserId)
		c.Request = c.Request.WithContext(ctx)
		// Lanjut ke handler
		c.Next()
	}
}
