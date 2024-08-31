package main

import (
	"fmt"
	"main/features"
	"main/middleware"
	"main/utils"

	swaggerDocs "main/docs"

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

	//핸드러 초기화

	if err := features.InitHandler(e); err != nil {
		fmt.Printf("handler 초기화 에러 : %v", err.Error())
		return
	}

	// swagger 초기화

	if utils.Env.IsLocal {
		swaggerDocs.SwaggerInfo.Host = "localhost:8080"
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	} else {

		swaggerDocs.SwaggerInfo.Host = fmt.Sprintf("%s-%s-api.jokertrickster.com", utils.Env.Env, "frog")
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}
	utils.LogInfo("frog Server started")
	e.HideBanner = true
	e.Logger.Fatal(e.Start(":" + utils.Env.Port))
	// e.Logger.Fatal(e.Start(":8080"))

	return
}
