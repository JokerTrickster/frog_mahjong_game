package middleware

import (
	"context"
	"fmt"
	"main/utils"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
)

// Logger : log middleware
func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//로깅 초기 세팅
		startTime := time.Now()
		req := c.Request()
		url := req.URL.Path
		requestID := random.String(32)

		c.Set("requestID", requestID)
		if req.Method == "GET" && url == "/health" {
			return next(c)
		}

		err := next(c)
		//에러 파싱
		resError := utils.Err{}
		var resCode int
		if c.Response().Status == 404 {
			err = utils.ErrorMsg(context.TODO(), utils.ErrNotFound, "", fmt.Sprintf("Invalid url call : %s", url), utils.ErrFromClient)
		}
		fmt.Println(err)
		if err != nil {
			resError = ErrorParsing(err.Error())
			resCode = resError.HttpCode
		} else {
			resCode = c.Response().Status
		}

		// 로깅
		logging := utils.Log{}
		logging.MakeLog("", url, req.Method, startTime, resCode, requestID)
		if resCode >= 400 {
			//에러 로깅
			logging.MakeErrorLog(resError)
			utils.LogError(logging)
			//DB 부하를 생각해서 에러만 쌓는걸로
			return echo.NewHTTPError(resError.HttpCode, utils.ErrType(resError.ErrType).New(resError.ErrType, resError.Msg))
		} else {
			utils.LogInfo(logging)
		}
		return err
	}
}

func ErrorParsing(data string) utils.Err {
	slice := strings.Split(data, "|")
	result := utils.Err{
		HttpCode: utils.ErrHttpCode[slice[0]],
		ErrType:  slice[0],
		Trace:    slice[1],
		Msg:      slice[2],
		From:     slice[3],
	}
	return result
}
