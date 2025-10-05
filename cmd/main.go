package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/federus1105/weekly/internals/configs"
	"github.com/federus1105/weekly/internals/routers"
	"github.com/joho/godotenv"
)

// @title		Weekly 10 Koda
// @version		1.0
// @description	Restful API craeted using gin for Koda Batch 3
// @host		localhost:8080
// @basepath	/
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Hanya load file .env di development
	if os.Getenv("ENV") != "production" {
		_ = godotenv.Load()
	}

	log.Println(os.Getenv("DBUSER"))

	// inisialisasi DB
	db, err := configs.InitDB()
	if err != nil {
		log.Println("❌ Failed to connect to database\nCause: ", err.Error())
		return
	}
	defer db.Close()

	log.Println("✅ DB Connected: ", err)

	// Inisialisasi RDB
	rdb, Rdb, err := configs.InitRDB()
	if err != nil {
		log.Println("❌ Failed to connect to redis\nCause: ", err.Error())
		return
	}
	defer rdb.Close()
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		fmt.Println("Failed Connected Redis : ", err.Error())
		return
	}
	log.Println("✅ REDIS Connected: ", Rdb)

	router := routers.InitRouter(db, rdb)
	//
	// router.Run("0.0.0.0:8080")
	router.Run("localhost:8080")
}
