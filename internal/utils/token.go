package utils

import (
	"errors"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
	})

	secret := []byte(os.Getenv("JWT_SECRET"))
	t, err := token.SignedString(secret)
	if err != nil {
		log.Println("Signed string issue")
		return "", err
	}

	return t, nil
}

func VerifyToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return false, err
	}

	return token.Valid, nil
}

// ExtractUserID extracts the user_id from the token if it's valid.
func ExtractUserID(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return "", err
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["user_id"].(string); ok {
			return userID, nil
		}
		return "", errors.New("user_id claim not found or invalid")
	}

	return "", errors.New("invalid token")
}