package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JSONWebToken struct {
	jwtKey []byte
}

func NewJSONWebToken(jwtSecret string) *JSONWebToken {
	return &JSONWebToken{
		jwtKey: []byte(jwtSecret),
	}
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (sa *JSONWebToken) GetJWTKey() []byte {
	return sa.jwtKey
}

func (sa *JSONWebToken) GenerateJSONWebTokens(id int64, username string, role int) (string, string, error) {
	accessToken, err := sa.generateShortLivedJSONWebToken(username)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := sa.generateLongLivedJSONWebToken(username)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (sa *JSONWebToken) generateShortLivedJSONWebToken(username string) (string, error) {
	expiration := time.Now().Add(5 * time.Minute)
	return sa.generateJSONWebToken(username, expiration)
}

func (sa *JSONWebToken) generateLongLivedJSONWebToken(username string) (string, error) {
	expiration := time.Now().Add(24 * time.Hour)
	return sa.generateJSONWebToken(username, expiration)
}

func (sa *JSONWebToken) generateJSONWebToken(username string, expirationTime time.Time) (string, error) {
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(sa.jwtKey)
}

func (sa *JSONWebToken) RefreshAccessToken(refreshToken string) (string, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return sa.jwtKey, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", errors.New("invalid refresh token")
	}

	newAccessToken, err := sa.generateShortLivedJSONWebToken(claims.Username)
	if err != nil {
		return "", err
	}
	return newAccessToken, nil
}
