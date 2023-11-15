package main

import (
	swaggerDocs "main/docs"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000" // 기본 포트 번호
	}
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!!!!!!!")
	})
	// swagger 초기화
	swaggerDocs.SwaggerInfo.Host = "localhost:" + port
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.HideBanner = true
	e.Logger.Fatal(e.Start(":" + port))

	return
}
