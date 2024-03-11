package main

import (
	"fmt"
	swaggerDocs "main/docs"
	"main/features"
	"main/middleware"
	"main/utils"

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
		fmt.Sprintln("handler 초기화 에러 : %s", err.Error())
		return
	}
	next := nextValue()

	println(next()) // 1
	println(next()) // 2
	println(next()) // 3

	anotherNext := nextValue()
	println(anotherNext()) // 1 다시 시작
	println(anotherNext()) // 2

	println(next()) // 4

	// swagger 초기화
	if utils.Env.IsLocal {
		swaggerDocs.SwaggerInfo.Host = "localhost:8080"
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	} else {
		swaggerDocs.SwaggerInfo.Host = fmt.Sprintf("%s-%s-api.breathings.net", utils.Env.Env, "frog")
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}
	e.HideBanner = true
	e.Logger.Fatal(e.Start(":" + utils.Env.Port))
	return
}

func nextValue() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}
