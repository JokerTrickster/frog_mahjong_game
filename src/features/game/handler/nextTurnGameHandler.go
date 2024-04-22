package handler

import (
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	mw "main/middleware"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type NextTurnGameHandler struct {
	UseCase _interface.INextTurnGameUseCase
}

func NewNextTurnGameHandler(c *echo.Echo, useCase _interface.INextTurnGameUseCase) _interface.INextTurnGameHandler {
	handler := &NextTurnGameHandler{
		UseCase: useCase,
	}

	c.PUT("/v0.1/game/next-turn", handler.NextTurn, mw.TokenChecker)
	return handler
}

// 다음 턴 넘기기
// @Router /v0.1/game/next-turn [put]
// @Summary 다음 턴 넘기기
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description NOT_OWNER : 방장이 시작 요청을 하지 않음
// @Description NOT_FIRST_PLAYER : 첫 플레이어가 아님
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Param tkn header string true "accessToken"
// @Param json body request.ReqNextTurn true "json body"
// @Produce json
// @Success 200 {object} boolean
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags game
func (d *NextTurnGameHandler) NextTurn(c echo.Context) error {
	ctx, _, _ := utils.CtxGenerate(c)
	req := &request.ReqNextTurn{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}
	err := d.UseCase.NextTurn(ctx, req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, true)
}
