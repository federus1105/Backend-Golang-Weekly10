package utils

// import (
// 	"errors"
// 	"regexp"

// 	// "github.com/federus1105/weekly/internals/models"
// 	// "github.com/federus1105/weekly/internals/models"
// )

// var (
// 	lowerRegex   = regexp.MustCompile(`[a-z]`)
// 	upperRegex   = regexp.MustCompile(`[A-Z]`)
// 	specialRegex = regexp.MustCompile(`[!@#$%^&*/><]`)
// )

// func Validate(body models.User) error {
// 	password := body.Password

// 	if len(password) < 8 {
// 		return errors.New("Password minimal harus terdiri dari 8 karakter")
// 	}
// 	if !lowerRegex.MatchString(password) {
// 		return errors.New("Password harus mengandung huruf kecil")
// 	}
// 	if !upperRegex.MatchString(password) {
// 		return errors.New("Password harus mengandung huruf besar")
// 	}
// 	if !specialRegex.MatchString(password) {
// 		return errors.New("Password harus mengandung karakter spesial")
// 	}

// 	return nil
// }
