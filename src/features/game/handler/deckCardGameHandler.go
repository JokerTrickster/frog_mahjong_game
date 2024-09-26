package handler

import (
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	mw "main/middleware"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type DeckCardGameHandler struct {
	UseCase _interface.IDeckCardGameUseCase
}

func NewDeckCardGameHandler(c *echo.Echo, useCase _interface.IDeckCardGameUseCase) _interface.IDeckCardGameHandler {
	handler := &DeckCardGameHandler{
		UseCase: useCase,
	}
	c.GET("/v0.1/game/card/deck", handler.DeckCard, mw.TokenChecker)
	return handler
}

// 카드 버리기
// @Router /v0.1/game/card/deck [get]
// @Summary 카드 버리기
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
// @Param query string roomID true "json body"
// @Produce json
// @Success 200 {object} response.ResDeckCard
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags game
func (d *DeckCardGameHandler) DeckCard(c echo.Context) error {
	ctx, userID, _ := utils.CtxGenerate(c)
	req := &request.ReqDeckCard{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}
	res, err := d.UseCase.DeckCard(ctx, int(userID), req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}
