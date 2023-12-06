package main

import (
	"fmt"
	swaggerDocs "main/docs"
	"main/middleware"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

//export PATH=$PATH:~/go/bin
func main() {
	e := echo.New()
	if err := utils.InitServer(); err != nil {
		fmt.Println(err)
		return
	}
	if err := middleware.InitMiddleware(e); err != nil {
		fmt.Println(err)
		return
	}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!!!!!!!")
	})
	// swagger 초기화
	swaggerDocs.SwaggerInfo.Host = "localhost:" + utils.Env.Port
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.HideBanner = true
	e.Logger.Fatal(e.Start(":" + utils.Env.Port))
	return
}
