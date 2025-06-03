package handler

import (
	_interface "main/features/board_game/model/interface"
	"main/features/board_game/model/request"
	mw "main/middleware"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type GameOverBoardGameHandler struct {
	UseCase _interface.IGameOverBoardGameUseCase
}

func NewGameOverBoardGameHandler(c *echo.Echo, useCase _interface.IGameOverBoardGameUseCase) _interface.IGameOverBoardGameHandler {
	handler := &GameOverBoardGameHandler{
		UseCase: useCase,
	}
	c.POST("/v0.1/game-over", handler.GameOver, mw.TokenChecker)
	return handler
}

// 게임 종료 처리
// @Router /board-game/v0.1/game-over [post]
// @Summary 게임 종료 처리
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description USER_NOT_FOUND : 유저를 찾을 수 없음
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Param tkn header string true "accessToken"
// @Param json body request.ReqGameOverBoardGame true "json body"
// @Produce json
// @Success 200 {object} bool
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags app/board-game/game
func (d *GameOverBoardGameHandler) GameOver(c echo.Context) error {
	ctx, userID, _ := utils.CtxGenerate(c)
	req := &request.ReqGameOverBoardGame{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}
	err := d.UseCase.GameOver(ctx, int(userID), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, true)
}
