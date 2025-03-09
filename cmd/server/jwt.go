package main

import (
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ClientID string `json:"clientID"`
	Email    string `json:"email"`
	Role     string `json:"role"` // "player" или "spectator"
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

var jwtKey = []byte("your-secret-key") // В реальном проекте используйте безопасный ключ

func generateTokenPair(user User, role string) (TokenPair, error) {
	accessClaims := Claims{
		ClientID: user.ID,
		Email:    user.Email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		return TokenPair{}, err
	}

	refreshClaims := Claims{
		ClientID: accessClaims.ClientID,
		Email:    user.Email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{AccessToken: accessTokenString, RefreshToken: refreshTokenString}, nil
}

func refreshToken(refreshTokenString string) (TokenPair, error) {
	token, err := jwt.ParseWithClaims(refreshTokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return TokenPair{}, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || claims.ExpiresAt.Before(time.Now()) {
		return TokenPair{}, fmt.Errorf("invalid or expired refresh token")
	}

	return generateTokenPair(User{Email: claims.Email, Name: claims.Email}, claims.Role)
}
