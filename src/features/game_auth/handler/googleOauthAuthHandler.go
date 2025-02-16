package handler

import (
	"context"
	_interface "main/features/game_auth/model/interface"
	"main/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type GoogleOauthAuthHandler struct {
	UseCase _interface.IGoogleOauthAuthUseCase
}

func NewGoogleOauthAuthHandler(c *echo.Echo, useCase _interface.IGoogleOauthAuthUseCase) _interface.IGoogleOauthAuthHandler {
	handler := &GoogleOauthAuthHandler{
		UseCase: useCase,
	}
	c.GET("/v0.1/game/auth/google", handler.GoogleOauth)
	return handler
}

// google oauth 로그인
// @Router /v0.1/game/auth/google [get]
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
// @Produce json
// @Success 200 {object} bool
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags game/auth
func (d *GoogleOauthAuthHandler) GoogleOauth(c echo.Context) error {
	ctx := context.Background()
	state, err := d.UseCase.GoogleOauth(ctx)
	if err != nil {
		return err
	}
	cookie := new(http.Cookie)
	cookie.Name = "state"
	cookie.Value = state
	cookie.Expires = time.Now().Add(1 * time.Minute)
	c.SetCookie(cookie)
	return c.Redirect(http.StatusTemporaryRedirect, utils.GoogleConfig.AuthCodeURL(state))
}
