package handler

import (
	"context"
	_interface "main/features/game_auth/model/interface"
	"main/features/game_auth/model/request"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SigninAuthHandler struct {
	UseCase _interface.ISigninAuthUseCase
}

func NewSigninAuthHandler(c *echo.Echo, useCase _interface.ISigninAuthUseCase) _interface.ISigninAuthHandler {
	handler := &SigninAuthHandler{
		UseCase: useCase,
	}
	c.POST("/v0.1/game/auth/signin", handler.Signin)
	return handler
}

// 로그인
// @Router /v0.1/game/auth/signin [post]
// @Summary 로그인
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description USER_NOT_EXIST : 유저가 존재하지 않음
// @Description USER_ALREADY_EXISTED : 유저가 이미 존재
// @Description USER_GOOGLE_ALREADY_EXISTED : 구글 계정이 이미 존재
// @Description PASSWORD_NOT_MATCH : 비밀번호가 일치하지 않음
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Param json body request.ReqGameSignin true "이메일, 비밀번호"
// @Produce json
// @Success 200 {object} response.ResGameSignin
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags game/auth
func (d *SigninAuthHandler) Signin(c echo.Context) error {
	ctx := context.Background()
	req := &request.ReqGameSignin{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}
	res, err := d.UseCase.Signin(ctx, req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}
