package main

import (
	"fmt"
	swaggerDocs "main/docs"
	"main/middleware"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

//export PATH=$PATH:~/go/bin
func main() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000" // 기본 포트 번호
	}
	e := echo.New()
	if err := middleware.InitMiddleware(e); err != nil {
		fmt.Println(err)
		return
	}

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
