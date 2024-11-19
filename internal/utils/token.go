package utils

import (
	"errors"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v4"
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
		log.Println(claims)
		if userID, ok := claims["user_id"].(string); ok {
			return userID, nil
		}
		return "", errors.New("user_id claim not found or invalid")
	}

	return "", errors.New("invalid token")
}

// ParseToken parses and validates the JWT token
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	// Retrieve the secret key from the environment variable
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return nil, errors.New("JWT_SECRET is not set in environment variables")
	}

	jwtSecret := []byte(secretKey)

	// Parse the token string
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure that the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	// Error handling for token parsing
	if err != nil {
		// Handle specific JWT parsing errors (invalid signature, expired token, etc.)
		if ve, ok := err.(*jwt.ValidationError); ok {
			// Depending on the error type, you can return different messages
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("malformed token")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("token has expired")
			} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				return nil, errors.New("invalid signature")
			} else {
				return nil, errors.New("invalid token")
			}
		}
		return nil, err
	}

	// If the token claims are valid, return them
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
