package handler

import (
	"context"
	"main/features/game_auth/model/entity"
	_interface "main/features/game_auth/model/interface"
	"main/features/game_auth/model/request"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RequestSignupAuthHandler struct {
	UseCase _interface.IRequestSignupAuthUseCase
}

func NewRequestSignupAuthHandler(c *echo.Echo, useCase _interface.IRequestSignupAuthUseCase) _interface.IRequestSignupAuthHandler {
	handler := &RequestSignupAuthHandler{
		UseCase: useCase,
	}
	c.POST("/v0.1/game/auth/signup/request", handler.RequestSignup)
	return handler
}

// 이메일 인증 요청
// @Router /v0.1/game/auth/signup/request [post]
// @Summary 이메일 인증 요청
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description USER_NOT_EXIST : 유저가 존재하지 않음
// @Description USER_ALREADY_EXISTED : 유저가 이미 존재
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Param json body request.ReqGameRequestSignup true "email"
// @Produce json
// @Success 200 {object} bool
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags game/auth
func (d *RequestSignupAuthHandler) RequestSignup(c echo.Context) error {
	ctx := context.Background()
	req := &request.ReqGameRequestSignup{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}

	entity := entity.RequestSignupAuthEntity{
		Email: req.Email,
	}
	_, err := d.UseCase.RequestSignup(ctx, entity)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, true)
}
