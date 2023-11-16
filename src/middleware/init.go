package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitMiddleware(e *echo.Echo) error {
	//cors 미들웨어 설정
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	//RestLogger : 로깅 미들웨어
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	//API sever timeout 24s
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultTimeoutConfig.Skipper,
		ErrorMessage: "timeout",
		Timeout:      180 * time.Second,
	}))

	return nil
}
