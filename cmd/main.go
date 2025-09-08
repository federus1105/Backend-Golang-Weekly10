package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/federus1105/weekly/internals/configs"
	"github.com/federus1105/weekly/internals/repositories"
	"github.com/federus1105/weekly/internals/routers"
	"github.com/federus1105/weekly/pkg"
	"github.com/joho/godotenv"
)

// @title		Weekly 10 Koda
// @version		1.0
// @description	Restful API craeted using gin for Koda Batch 3
// @host		localhost:8080
// @basepath	/
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Failed to load env\nCause:", err.Error())
		return
	}
	log.Println(os.Getenv("DBUSER"))

	// inisialisasi DB
	db, err := configs.InitDB()
	if err != nil {
		log.Println("Failed to connect to database\nCause: ", err.Error())
		return
	}
	defer db.Close()

	log.Println("DB Connected")

	router := routers.InitRouter(db)

	router.Run("localhost:8080")
	hc := pkg.NewHashConfig()
	hc.UseRecommended()

	repo := repositories.NewAuthRepository(db)
	users, err := repo.GetAllUsers(context.Background())
	if err != nil {
		log.Fatal("Error GetAllUsers:", err)
	}

	for _, u := range users {
		fmt.Println(u.Email, u.Role)
	}

	// users, _ := hc.Login()
	// for _, user := range users {
	// 	hashed, err := hc.GenHash(user.Password)
	// 	if err != nil {
	// 		log.Println("Gagal hash password:", user.Email)
	// 		continue
	// 	}

	// 	err = hc.UpdateUserPassword(user.Id, hashed)
	// 	if err != nil {
	// 		log.Println("Gagal update password user:", user.Email)
	// 	} else {
	// 		log.Println("Password berhasil di-hash untuk:", user.Email)
	// 	}
	// }
	// password := "koda"
	// hash, err := hc.GenHash(password)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// isMatch, err := hc.CompareHashAndPassword(password, hash)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(isMatch)
}
