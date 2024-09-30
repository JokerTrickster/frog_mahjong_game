package utils

import (
	"fmt"
	"time"

	"context"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4/middleware"
)

type JwtCustomClaims struct {
	CreateTime int64  `json:"createTime"`
	UserID     uint   `json:"userID"`
	Email      string `json:"email"`
	jwt.StandardClaims
}

var AccessTokenSecretKey []byte
var RefreshTokenSecretKey []byte
var JwtConfig middleware.JWTConfig

const (
	AccessTokenExpiredTime  = 24         //hours
	RefreshTokenExpiredTime = 1 * 24 * 7 //hours
)

func InitJwt() error {
	secret := "secret"
	AccessTokenSecretKey = []byte(secret)
	RefreshTokenSecretKey = []byte(secret)
	return nil
}

func GenerateToken(email string, userID uint) (string, int64, string, int64, error) {
	now := time.Now()
	accessToken, accessTknExpiredAt, err := GenerateAccessToken(email, now, userID)
	if err != nil {
		return "", 0, "", 0, err
	}
	refreshToken, refreshTknExpiredAt, err := GenerateRefreshToken(email, now, userID)
	if err != nil {
		return "", 0, "", 0, err
	}
	return accessToken, accessTknExpiredAt, refreshToken, refreshTknExpiredAt, nil
}

func GenerateAccessToken(email string, now time.Time, userID uint) (string, int64, error) {
	// Set custom claims
	expiredAt := now.Add(time.Hour * AccessTokenExpiredTime).Unix()
	claims := &JwtCustomClaims{
		TimeToEpochMillis(now),
		userID,
		email,
		jwt.StandardClaims{
			ExpiresAt: expiredAt,
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	accessToken, err := token.SignedString(AccessTokenSecretKey)
	if err != nil {
		return "", 0, err
	}
	return accessToken, expiredAt, nil
}

func GenerateRefreshToken(email string, now time.Time, userID uint) (string, int64, error) {
	expiredAt := now.Add(time.Hour * RefreshTokenExpiredTime).Unix()
	claims := &JwtCustomClaims{
		TimeToEpochMillis(now),
		userID,
		email,
		jwt.StandardClaims{
			ExpiresAt: expiredAt,
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	refreshToken, err := token.SignedString(RefreshTokenSecretKey)
	if err != nil {
		return "", 0, ErrorMsg(context.TODO(), ErrBadToken, Trace(), fmt.Sprintf("failed to generate refresh token - %v", err), ErrFromInternal)
	}
	return refreshToken, expiredAt, nil
}
func VerifyToken(tokenString string) error {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return AccessTokenSecretKey, nil
	})
	if err != nil {
		return ErrorMsg(context.TODO(), ErrBadToken, Trace(), fmt.Sprintf("failed to parse token - %v err - %v", token, err.Error()), ErrFromClient)
	}

	// Check token validity
	if !token.Valid {
		return ErrorMsg(context.TODO(), ErrBadToken, Trace(), fmt.Sprintf("invalid token - %v ", token), ErrFromClient)
	}
	return nil
}
func ParseToken(tokenString string) (uint, string, error) {
	token, _ := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return AccessTokenSecretKey, nil
	})
	// if err != nil {
	// 	return 0, "", ErrorMsg(context.TODO(), ErrBadToken, Trace(), fmt.Sprintf("failed to parse token - %v", token), ErrFromClient)
	// }
	// Extract claims
	claims, ok := token.Claims.(*JwtCustomClaims)
	if !ok {
		return 0, "", ErrorMsg(context.TODO(), ErrBadToken, Trace(), fmt.Sprintf("failed to extract claims - %v", token), ErrFromClient)
	}
	fmt.Println("claims: ", claims)
	// Extract email and userID
	email := claims.Email
	userID := claims.UserID
	return userID, email, nil
}
