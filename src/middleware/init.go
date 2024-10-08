package middleware

import (
	"fmt"
	"main/utils"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var Store = sessions.NewCookieStore([]byte("secret"))

func InitMiddleware(e *echo.Echo) error {
	e.Use(middleware.Recover())

	//cors 미들웨어 설정
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	//API sever timeout 24s
	// TODO: websocket 연결안되는 이슈 확인 필요
	// e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
	// 	Skipper:      middleware.DefaultTimeoutConfig.Skipper,
	// 	ErrorMessage: "timeout",
	// 	Timeout:      24 * time.Second,
	// }))

	//Logger : 로깅 미들웨어
	e.Use(Logger)

	//jwt 검증 미들웨어
	signingKey := utils.AccessTokenSecretKey
	utils.JwtConfig = middleware.JWTConfig{
		TokenLookup:   "cookie:accessToken",
		SigningKey:    signingKey,
		SigningMethod: "HS256",
		ParseTokenFunc: func(auth string, c echo.Context) (interface{}, error) {
			keyFunc := func(t *jwt.Token) (interface{}, error) {
				if t.Method.Alg() != "HS256" {
					return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
				}
				return signingKey, nil
			}
			// claims are of type `jwt.MapClaims` when token is created with `jwt.Parse`
			token, err := jwt.Parse(auth, keyFunc)
			if err != nil {
				return nil, err
			}
			if !token.Valid {
				return nil, fmt.Errorf("invalid token")
			}
			return token, nil
		},
	}
	return nil
}
