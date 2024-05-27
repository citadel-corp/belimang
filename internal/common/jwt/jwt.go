package jwt

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	key = []byte(os.Getenv("JWT_SECRET"))

	ErrUnknownClaims = errors.New("unknown claims type")
	ErrTokenInvalid  = errors.New("invalid token")
)

type UserClaims struct {
	UserUID string `json:"uid"`
	Role    string `json:"role"`
	jwt.RegisteredClaims
}

func Sign(ttl time.Duration, subject, role string) (string, error) {
	now := time.Now()
	expiry := now.Add(ttl)

	claims := UserClaims{
		UserUID: subject,
		Role:    role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiry),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(key)
}

func VerifyAndGetSubject(tokenString string) (*UserClaims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, ErrUnknownClaims
	}
}
