package helper

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserUID string `json:"user_uid"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Exp     int64  `json:"exp"`
	jwt.RegisteredClaims
}

type Token interface {
	Create(i *Claims) (*string, error)
	Extract(tokenString string) (*Claims, error)
}

type JwtToken struct {
	jwtKey []byte
}

func NewJwtToken(jwtKey []byte) Token {
	return &JwtToken{jwtKey: jwtKey}
}

// Create implements Token.
func (j *JwtToken) Create(i *Claims) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_uid": i.UserUID,
		"name":     i.Name,
		"email":    i.Email,
		"exp":      i.Exp,
	})

	tokenString, err := token.SignedString(j.jwtKey)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

// Extract implements Token.
func (j *JwtToken) Extract(tokenString string) (*Claims, error) {
	var claims Claims

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.jwtKey, nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return &claims, nil
}

func CreateToken(i *Claims, jwtKey []byte) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_uid": i.UserUID,
		"name":     i.Name,
		"email":    i.Email,
		"exp":      i.Exp,
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func ExtractToken(tokenString string, jwtKey []byte) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtKey, nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
