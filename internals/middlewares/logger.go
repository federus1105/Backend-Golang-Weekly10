package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func MyLogger(ctx *gin.Context) {
	log.Println("Start")
	start := time.Now()
	ctx.Next() // next digunakan untuk lanjut ke middleware selanjutnya
	duration := time.Since(start)
	log.Printf("Durasi Request: %dus\n", duration.Microseconds())
}
