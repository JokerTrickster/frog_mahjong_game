package handler

import (
	_interface "main/features/board_game/model/interface"
	mw "main/middleware"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type FindItRankBoardGameHandler struct {
	UseCase _interface.IFindItRankBoardGameUseCase
}

func NewFindItRankBoardGameHandler(c *echo.Echo, useCase _interface.IFindItRankBoardGameUseCase) _interface.IFindItRankBoardGameHandler {
	handler := &FindItRankBoardGameHandler{
		UseCase: useCase,
	}
	c.GET("/find-it/v0.1/game/rank", handler.FindItRank, mw.TokenChecker)
	return handler
}

// 틀린그림찾기 랭킹 가져오기
// @Router /find-it/v0.1/game/rank [get]
// @Summary 틀린그림찾기 랭킹 가져오기
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
// @Success 200 {object} response.ResFindItRankBoardGame
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags app/find-it/game
func (d *FindItRankBoardGameHandler) FindItRank(c echo.Context) error {
	ctx, _, _ := utils.CtxGenerate(c)

	res, err := d.UseCase.FindItRank(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
