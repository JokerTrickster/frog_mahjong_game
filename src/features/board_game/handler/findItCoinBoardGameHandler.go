package handler

import (
	_interface "main/features/board_game/model/interface"
	"main/features/board_game/model/request"
	mw "main/middleware"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type FindItCoinBoardGameHandler struct {
	UseCase _interface.IFindItCoinBoardGameUseCase
}

func NewFindItCoinBoardGameHandler(c *echo.Echo, useCase _interface.IFindItCoinBoardGameUseCase) _interface.IFindItCoinBoardGameHandler {
	handler := &FindItCoinBoardGameHandler{
		UseCase: useCase,
	}
	c.POST("/board-game/v0.1/game/coin", handler.FindItCoin, mw.TokenChecker)
	return handler
}

// 코인 값 변경 api
// @Router /board-game/v0.1/game/coin [post]
// @Summary 코인 값 변경 api
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description USER_NOT_FOUND : 유저를 찾을 수 없음
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Param tkn header string true "accessToken"
// @Param json body request.ReqFindItCoinBoardGame true "플레이 라운드 수"
// @Produce json
// @Success 200 {object} bool
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags app/board-game/game
func (d *FindItCoinBoardGameHandler) FindItCoin(c echo.Context) error {
	ctx, userID, _ := utils.CtxGenerate(c)
	req := &request.ReqFindItCoinBoardGame{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}
	err := d.UseCase.FindItCoin(ctx, int(userID), req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, true)
}
