package handler

import (
	_interface "main/features/board_game/model/interface"
	mw "main/middleware"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SlimeWarRankBoardGameHandler struct {
	UseCase _interface.ISlimeWarRankBoardGameUseCase
}

func NewSlimeWarRankBoardGameHandler(c *echo.Echo, useCase _interface.ISlimeWarRankBoardGameUseCase) _interface.ISlimeWarRankBoardGameHandler {
	handler := &SlimeWarRankBoardGameHandler{
		UseCase: useCase,
	}
	c.GET("/slime-war/v0.1/game/rank", handler.SlimeWarRank, mw.TokenChecker)
	return handler
}

// 슬라임 전쟁 랭크 가져오기
// @Router /slime-war/v0.1/game/rank [get]
// @Summary 슬라임 전쟁 랭크 가져오기
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description USER_NOT_FOUND : 유저를 찾을 수 없음
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Param tkn header string true "accessToken"
// @Produce json
// @Success 200 {object} response.ResSlimeWarRankBoardGame
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags app/slime-war/game
func (d *SlimeWarRankBoardGameHandler) SlimeWarRank(c echo.Context) error {
	ctx, _, _ := utils.CtxGenerate(c)

	res, err := d.UseCase.SlimeWarRank(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
