package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"main/utils"
	"net/http"
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
		requestID := random.String(32)
		c.Set("rID", requestID)
		c.Set("startTime", startTime)
		req := c.Request()
		url := req.URL.Path
		if req.Method == "GET" && url == "/health" {
			return next(c)
		}

		// 요청으로부터 JSON, 쿼리, 패스 파라미터 값을 추출하여 출력
		// JSON Body
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil && err != io.EOF {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
		}
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 원래 요청에 복사한 바디를 재설정

		// JSON Body를 읽어서 출력
		var requestBody map[string]interface{}
		queryParams := map[string][]string{}
		pathValues := make(map[string]string)
		if req.Method == "POST" || req.Method == "PUT" || req.Method == "PATCH" {
			if err := json.Unmarshal(bodyBytes, &requestBody); err != nil {
				fmt.Println("Failed to unmarshal JSON body:", err)
			}
		} else {
			// Query Parameters
			queryParams = c.QueryParams()

			// Path Parameters
			pathParams := c.ParamNames()
			for _, param := range pathParams {
				pathValues[param] = c.Param(param)
			}
		}

		err = next(c)
		//에러 파싱
		resError := utils.Err{}
		var resCode int
		if c.Response().Status == 404 {
			err = utils.ErrorMsg(context.TODO(), utils.ErrNotFound, "", fmt.Sprintf("Invalid url call : %s", url), utils.ErrFromClient)
		}
		if err != nil {
			resError = ErrorParsing(err.Error())
			resCode = resError.HttpCode
		} else {
			resCode = c.Response().Status
		}
		// 로깅
		logging := utils.Log{}
		logging.MakeLog("", url, req.Method, startTime, resCode, requestID, requestBody, queryParams, pathValues)
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
