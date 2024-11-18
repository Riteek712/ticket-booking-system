package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePassword(p string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(hash)
}

func ComparePassword(hashedPassword, password string) bool {
	log.Println("Comparing password")
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
