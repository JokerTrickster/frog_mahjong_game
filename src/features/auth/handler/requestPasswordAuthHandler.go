package handler

import (
	"context"
	"main/features/auth/model/entity"
	_interface "main/features/auth/model/interface"
	"main/features/auth/model/request"
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
	c.POST("/v0.1/auth/password/request", handler.RequestPassword)
	return handler
}

// 비밀번호 변경 요청
// @Router /v0.1/auth/password/request [post]
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
// @Param json body request.ReqRequestPassword true "email"
// @Produce json
// @Success 200 {object} string
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags auth
func (d *RequestPasswordAuthHandler) RequestPassword(c echo.Context) error {
	ctx := context.Background()
	req := &request.ReqRequestPassword{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}

	entity := entity.RequestPasswordAuthEntity{
		Email: req.Email,
	}
	code, err := d.UseCase.RequestPassword(ctx, entity)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, code)
}
