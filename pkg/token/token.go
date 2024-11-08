package token

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"

	"kbox-api/shared/httperror"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func GenerateToken(jwtKey []byte, userID, role string, duration time.Duration) (string, error) {
	expirationTime := time.Now().Add(duration)
	claims := &Claims{
		UserID: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "",
			httperror.New(
				http.StatusInternalServerError,
				err.Error(),
			)
	}
	return tokenString, nil
}

func ParseToken(tokenString string, jwtKey []byte) (*Claims, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(_ *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, httperror.New(
			http.StatusUnauthorized,
			"Недействительный токен",
		)
	}
	return claims, nil
}
