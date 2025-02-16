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

type ValidatePasswordAuthHandler struct {
	UseCase _interface.IValidatePasswordAuthUseCase
}

func NewValidatePasswordAuthHandler(c *echo.Echo, useCase _interface.IValidatePasswordAuthUseCase) _interface.IValidatePasswordAuthHandler {
	handler := &ValidatePasswordAuthHandler{
		UseCase: useCase,
	}
	c.POST("/v0.1/game/auth/password/validate", handler.ValidatePassword)
	return handler
}

// 비밀번호 변경 검증
// @Router /v0.1/game/auth/password/validate [post]
// @Summary 비밀번호 변경 검증
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description USER_NOT_EXIST : 유저가 존재하지 않음
// @Description USER_ALREADY_EXISTED : 유저가 이미 존재
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Param json body request.ReqGameValidatePassword true "email, password, code"
// @Produce json
// @Success 200 {object} bool
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags game/auth
func (d *ValidatePasswordAuthHandler) ValidatePassword(c echo.Context) error {
	ctx := context.Background()
	req := &request.ReqGameValidatePassword{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}

	entity := entity.ValidatePasswordAuthEntity{
		Email:    req.Email,
		Password: req.Password,
		Code:     req.Code,
	}
	err := d.UseCase.ValidatePassword(ctx, entity)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, true)
}
