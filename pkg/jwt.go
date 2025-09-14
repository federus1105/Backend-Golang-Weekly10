package pkg

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId int    `json:"id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTClaims(userid int, role string) *Claims {
	return &Claims{
		UserId: userid,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			Issuer:    os.Getenv("JWT_ISSUER"),
		},
	}
}

func (c *Claims) GenToken() (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("no secret found")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(jwtSecret))
}

func (c *Claims) VerifyToken(token string) error {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return errors.New("no secret found")
	}
	parsedToken, err := jwt.ParseWithClaims(token, c, func(t *jwt.Token) (any, error) { return []byte(jwtSecret), nil })
	if err != nil {
		return err
	}
	if !parsedToken.Valid {
		return jwt.ErrTokenExpired
	}
	iss, err := parsedToken.Claims.GetIssuer()
	if err != nil {
		return err
	}
	if iss != os.Getenv("JWT_ISSUER") {
		return jwt.ErrTokenInvalidIssuer
	}
	return nil
}

// func VerifyToken(ctx *gin.Context) {
// 	bearerToken := ctx.GetHeader("Authorization")
// 	if bearerToken == "" {
// 		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 			"success": false,
// 			"message": "Authorization header required",
// 		})
// 		return
// 	}

// 	// Pastikan token ada dan sesuai format "Bearer <token>"
// 	parts := strings.Split(bearerToken, " ")
// 	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
// 		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 			"success": false,
// 			"message": "Authorization header format must be Bearer {token}",
// 		})
// 		return
// 	}

// 	tokenString := parts[1]

// 	// lanjutkan validasi JWT (decode, verify claims, dsb)
// 	claims, err := your_jwt_library.VerifyToken(tokenString)
// 	if err != nil {
// 		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 			"success": false,
// 			"message": "Invalid token",
// 		})
// 		return
// 	}

// 	// Simpan claims ke context
// 	ctx.Set("claims", claims)
// 	ctx.Next()
// }
