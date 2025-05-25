package jwt

import (
	"fmt"
	"hmb_fighting/server/entities"
	"hmb_fighting/server/types"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ClientID  string              `json:"clientID"`
	Email     string              `json:"email"`
	Role      types.UserRoleTypes `json:"role"` // "player" или "spectator"
	TokenType string              `json:"tokenType"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

var jwtKey = []byte("your-secret-key") // В реальном проекте используйте безопасный ключ

func GenerateTokenPair(user entities.User, role types.UserRoleTypes) (TokenPair, error) {
	accessClaims := Claims{
		ClientID:  user.ID,
		Email:     user.Email,
		Role:      role,
		TokenType: "access",
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
		ClientID:  accessClaims.ClientID,
		Email:     user.Email,
		Role:      role,
		TokenType: "refresh",
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

func RefreshToken(refreshTokenString string) (TokenPair, error) {
	token, err := jwt.ParseWithClaims(refreshTokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return TokenPair{}, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || claims.ExpiresAt == nil || claims.ExpiresAt.Before(time.Now()) || claims.TokenType != "refresh" {
		return TokenPair{}, fmt.Errorf("invalid or expired refresh token")
	}

	return GenerateTokenPair(entities.User{ID: claims.ClientID, Email: claims.Email, Name: claims.Email}, claims.Role)
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %v", err)
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || claims.ExpiresAt == nil || claims.ExpiresAt.Before(time.Now()) || claims.TokenType != "access" {
		return nil, fmt.Errorf("invalid or expired token")
	}
	return claims, nil
}
