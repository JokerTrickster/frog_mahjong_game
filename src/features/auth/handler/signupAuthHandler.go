package handler

import (
	"context"
	_interface "main/features/auth/model/interface"
	"main/features/auth/model/request"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SignupAuthHandler struct {
	UseCase _interface.ISignupAuthUseCase
}

func NewSignupAuthHandler(c *echo.Echo, useCase _interface.ISignupAuthUseCase) _interface.ISignupAuthHandler {
	handler := &SignupAuthHandler{
		UseCase: useCase,
	}
	c.POST("/v0.1/auth/signup", handler.Signup)
	return handler
}

// 회원 가입
// @Router /v0.1/auth/signup [post]
// @Summary 회원 가입
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description
// @Description ■ errCode with 400
// @Description USER_ALREADY_EXISTED : 유저가 이미 존재
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Param json body request.ReqSignup true "이름, 이메일, 비밀번호"
// @Produce json
// @Success 200 {object} int
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags auth
func (d *SignupAuthHandler) Signup(c echo.Context) error {
	ctx := context.Background()
	req := &request.ReqSignup{}
	if err := utils.ValidateReq(c, req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err := d.UseCase.Signup(ctx, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, true)
}
