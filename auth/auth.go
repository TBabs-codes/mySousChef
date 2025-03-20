package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Generates a hash from user's password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hash), err
}

// Compares password and stored hash securely
func CheckPasswordHash(password, hash string) error {
	password_byte := []byte(password)
	hash_byte := []byte(hash)
	return bcrypt.CompareHashAndPassword(hash_byte, password_byte)
}

// Creates JSON web token
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	signingMethod := jwt.GetSigningMethod("HS256")
	claim := jwt.RegisteredClaims{
		Issuer:    "MealsMe",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(signingMethod, claim)

	return token.SignedString([]byte(tokenSecret))
}

//Validates JSON web token
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claim := &jwt.RegisteredClaims{}
	keyFunc := func(*jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, claim, keyFunc)
	if err != nil {
		//401 Unauthorized
		return uuid.UUID{}, fmt.Errorf("parsing failed, action unauthorized: %v", err)
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("invalid token claims")
	}

	return uuid.MustParse(claims.Subject), nil
}

//Obtains token string from header
func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if len(authHeader) < 5 {
		return "", fmt.Errorf("no authorization header included")
	} else if authHeader[0:6] != "Bearer" {
		return "", fmt.Errorf("malformed authorization header")
	} else {
		return strings.TrimPrefix(authHeader, "Bearer "), nil
	}
}

