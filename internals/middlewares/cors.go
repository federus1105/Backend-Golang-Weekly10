package middlewares

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware(ctx *gin.Context) {
	whitelist := []string{"http://localhost:5173", "http://192.168.100.152:5173"}
	origin := ctx.GetHeader("Origin")
	if slices.Contains(whitelist, origin) {
		ctx.Header("Access-Control-Allow-Origin", origin)
	}

	ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
	ctx.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")

	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}
	ctx.Next()
}
