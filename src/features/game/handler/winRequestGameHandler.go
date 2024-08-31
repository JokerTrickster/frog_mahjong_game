package handler

import (
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	mw "main/middleware"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type WinRequestGameHandler struct {
	UseCase _interface.IWinRequestGameUseCase
}

func NewWinRequestGameHandler(c *echo.Echo, useCase _interface.IWinRequestGameUseCase) _interface.IWinRequestGameHandler {
	handler := &WinRequestGameHandler{
		UseCase: useCase,
	}
	c.POST("/v0.1/game/win-request", handler.WinRequest, mw.TokenChecker)
	return handler
}

// 게임 승리 요청
// @Router /v0.1/game/win-request [post]
// @Summary 게임 승리 요청
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description NOT_ALL_USERS_READY : 모든 유저가 준비되지 않음
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Param tkn header string true "accessToken"
// @Param json body request.ReqWinRequest true "json body"
// @Produce json
// @Success 200 {object} boolean
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags game
func (d *WinRequestGameHandler) WinRequest(c echo.Context) error {
	ctx, _, _ := utils.CtxGenerate(c)
	req := &request.ReqWinRequest{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}
	result, err := d.UseCase.WinRequest(ctx, req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}
