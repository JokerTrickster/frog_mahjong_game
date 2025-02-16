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

type RequestPasswordAuthHandler struct {
	UseCase _interface.IRequestPasswordAuthUseCase
}

func NewRequestPasswordAuthHandler(c *echo.Echo, useCase _interface.IRequestPasswordAuthUseCase) _interface.IRequestPasswordAuthHandler {
	handler := &RequestPasswordAuthHandler{
		UseCase: useCase,
	}
	c.POST("/v0.1/game/auth/password/request", handler.RequestPassword)
	return handler
}

// 비밀번호 변경 요청
// @Router /v0.1/game/auth/password/request [post]
// @Summary 비밀번호 변경 요청
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description USER_NOT_EXIST : 유저가 존재하지 않음
// @Description USER_ALREADY_EXISTED : 유저가 이미 존재
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Param json body request.ReqGameRequestPassword true "email"
// @Produce json
// @Success 200 {object} bool
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags game/auth
func (d *RequestPasswordAuthHandler) RequestPassword(c echo.Context) error {
	ctx := context.Background()
	req := &request.ReqGameRequestPassword{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}

	entity := entity.RequestPasswordAuthEntity{
		Email: req.Email,
	}
	_, err := d.UseCase.RequestPassword(ctx, entity)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, true)
}
