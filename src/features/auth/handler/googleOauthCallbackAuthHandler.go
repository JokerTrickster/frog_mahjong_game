package handler

import (
	"context"
	_interface "main/features/auth/model/interface"
	"net/http"

	"github.com/labstack/echo/v4"
)

type GoogleOauthCallbackAuthHandler struct {
	UseCase _interface.IGoogleOauthCallbackAuthUseCase
}

func NewGoogleOauthCallbackAuthHandler(c *echo.Echo, useCase _interface.IGoogleOauthCallbackAuthUseCase) _interface.IGoogleOauthCallbackAuthHandler {
	handler := &GoogleOauthCallbackAuthHandler{
		UseCase: useCase,
	}
	c.GET("/v0.1/auth/google/callback", handler.GoogleOauthCallback)
	return handler
}

// google oauth 로그인 콜백
// @Router /v0.1/auth/google/callback [get]
// @Summary google oauth 로그인 콜백
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description USER_NOT_EXIST : 유저가 존재하지 않음
// @Description USER_ALREADY_EXISTED : 유저가 이미 존재
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Produce json
// @Success 200 {object} response.GoogleOauthCallbackRes
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags auth
func (d *GoogleOauthCallbackAuthHandler) GoogleOauthCallback(c echo.Context) error {
	ctx := context.Background()
	oauthState, _ := c.Request().Cookie("state")
	if oauthState.Value != c.FormValue("state") {
		return c.JSON(http.StatusBadRequest, "invalid oauth state")
	}
	res, err := d.UseCase.GoogleOauthCallback(ctx, c.FormValue("code"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}
