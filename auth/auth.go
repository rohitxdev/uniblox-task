package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserID int `json:"userId"`
}

// GenerateLoginToken generates a login token for the claims.
func GenerateLoginToken(tokenClaims TokenClaims, jwtSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"userId": tokenClaims.UserID,
		"nbf":    time.Now().Unix(),
		"exp":    time.Now().Add(expiresIn).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("Failed to generate login token: %w", err)
	}
	return tokenStr, nil
}

// ValidateLoginToken validates the login token and returns the claims.
func ValidateLoginToken(tokenStr string, jwtSecret string) (*TokenClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to validate login token: %w", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Failed to validate login token: invalid claims")
	}
	tokenClaims := TokenClaims{
		UserID: int(claims["userId"].(float64)),
	}

	return &tokenClaims, nil
}
