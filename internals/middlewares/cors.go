package middlewares

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware(ctx *gin.Context) {
	whitelist := []string{"http://localhost:5500", "http://127.0.0.1:3001"}
	origin := ctx.GetHeader("Origin")
	if slices.Contains(whitelist, origin) {
		ctx.Header("Access-Control-Allow-Origin", origin)
	}

	ctx.Header("Access-Control-Allow-Methods", "GET")
	ctx.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")

	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}
	ctx.Next()
}
