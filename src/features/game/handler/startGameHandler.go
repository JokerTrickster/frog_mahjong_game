package handler

import (
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	mw "main/middleware"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type StartGameHandler struct {
	UseCase _interface.IStartGameUseCase
}

func NewStartGameHandler(c *echo.Echo, useCase _interface.IStartGameUseCase) _interface.IStartGameHandler {
	handler := &StartGameHandler{
		UseCase: useCase,
	}
	c.POST("/v0.1/game/start", handler.Start, mw.TokenChecker)
	return handler
}

// 게임 시작
// @Router /v0.1/game/start [post]
// @Summary 게임 시작
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Param tkn header string true "accessToken"
// @Param json body request.ReqStart true "json body"
// @Produce json
// @Success 200 {object} boolean
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags game
func (d *StartGameHandler) Start(c echo.Context) error {
	ctx, _, email := utils.CtxGenerate(c)
	req := &request.ReqStart{}
	if err := utils.ValidateReq(c, req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err := d.UseCase.Start(ctx, email, req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, true)
}
