package handler

import (
	"context"
	_interface "main/features/auth/model/interface"
	"main/utils"
	"net/http"

	"main/features/auth/model/request"

	"github.com/labstack/echo/v4"
)

type V02GoogleOauthCallbackAuthHandler struct {
	UseCase _interface.IV02GoogleOauthCallbackAuthUseCase
}

func NewV02GoogleOauthCallbackAuthHandler(c *echo.Echo, useCase _interface.IV02GoogleOauthCallbackAuthUseCase) _interface.IV02GoogleOauthCallbackAuthHandler {
	handler := &V02GoogleOauthCallbackAuthHandler{
		UseCase: useCase,
	}
	c.GET("/v0.2/auth/google/callback", handler.V02GoogleOauthCallback)
	return handler
}

// google oauth 로그인 콜백
// @Router /v0.2/auth/google/callback [get]
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
// @Param code query string true "code"
// @Produce json
// @Success 200 {object} response.ResV02GoogleOauthCallback
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags auth
func (d *V02GoogleOauthCallbackAuthHandler) V02GoogleOauthCallback(c echo.Context) error {
	ctx := context.Background()
	req := &request.ReqV02GoogleOauthCallback{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}

	res, err := d.UseCase.V02GoogleOauthCallback(ctx, req.Code)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}
