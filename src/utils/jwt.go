package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4/middleware"
)

type JwtCustomClaims struct {
	createTime int64  `json:"createTime"`
	UserID     uint   `json:"userID"`
	Email      string `json:"email"`
	jwt.StandardClaims
}

var AccessTokenSecretKey []byte
var RefreshTokenSecretKey []byte
var JwtConfig middleware.JWTConfig

const (
	AccessTokenExpiredTime  = 1
	RefreshTokenExpiredTime = 24 * 7
)

func InitJwt() error {
	secret := "secret"
	AccessTokenSecretKey = []byte(secret)
	RefreshTokenSecretKey = []byte(secret)
	return nil
}

func GenerateToken(email string, userID uint) (string, string, error) {
	now := time.Now()
	accessToken, err := GenerateAccessToken(email, now, userID)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := GenerateRefreshToken(email, now, userID)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func GenerateAccessToken(email string, now time.Time, userID uint) (string, error) {
	// Set custom claims
	claims := &JwtCustomClaims{
		TimeToEpochMillis(now),
		userID,
		email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * AccessTokenExpiredTime).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	accessToken, err := token.SignedString(AccessTokenSecretKey)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func GenerateRefreshToken(email string, now time.Time, userID uint) (string, error) {
	claims := &JwtCustomClaims{
		TimeToEpochMillis(now),
		userID,
		email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * RefreshTokenExpiredTime).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	refreshToken, err := token.SignedString(RefreshTokenSecretKey)
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

func ValidateAndParseAccessToken(tokenString string) (uint, string, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return AccessTokenSecretKey, nil
	})
	if err != nil {
		return 0, "", err
	}

	// Check token validity
	if !token.Valid {
		return 0, "", errors.New("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(*JwtCustomClaims)
	if !ok {
		return 0, "", errors.New("failed to parse claims")
	}

	// Extract email and userID
	email := claims.Email
	userID := claims.UserID

	return userID, email, nil
}
