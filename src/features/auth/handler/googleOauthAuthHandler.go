package handler

import (
	"context"
	"fmt"
	_interface "main/features/auth/model/interface"
	"main/features/auth/model/request"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type GoogleOauthAuthHandler struct {
	UseCase _interface.IGoogleOauthAuthUseCase
}

func NewGoogleOauthAuthHandler(c *echo.Echo, useCase _interface.IGoogleOauthAuthUseCase) _interface.IGoogleOauthAuthHandler {
	handler := &GoogleOauthAuthHandler{
		UseCase: useCase,
	}
	c.POST("/v0.1/auth/google", handler.GoogleOauth)
	return handler
}

// google oauth 로그인
// @Router /v0.1/auth/google [post]
// @Summary google oauth 로그인
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description USER_NOT_EXIST : 유저가 존재하지 않음
// @Description USER_ALREADY_EXISTED : 유저가 이미 존재
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Param json body request.ReqGoogleOauth true "json body"
// @Produce json
// @Success 200 {object} bool
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags auth
func (d *GoogleOauthAuthHandler) GoogleOauth(c echo.Context) error {
	ctx := context.Background()
	req := &request.ReqGoogleOauth{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}
	fmt.Println(req.Credential)
	err := d.UseCase.GoogleOauth(ctx)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, true)
}
