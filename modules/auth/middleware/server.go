package middleware

import (
	"diploma/modules/auth/jwt"
)

type AuthMiddleware struct {
	jwt *jwt.JSONWebToken
}

func NewAuthMiddleware(jwt *jwt.JSONWebToken) *AuthMiddleware {
	return &AuthMiddleware{
		jwt: jwt,
	}
}
