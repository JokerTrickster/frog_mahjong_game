package handler

import (
	_interface "main/features/board_game/model/interface"
	mw "main/middleware"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SequenceRankBoardGameHandler struct {
	UseCase _interface.ISequenceRankBoardGameUseCase
}

func NewSequenceRankBoardGameHandler(c *echo.Echo, useCase _interface.ISequenceRankBoardGameUseCase) _interface.ISequenceRankBoardGameHandler {
	handler := &SequenceRankBoardGameHandler{
		UseCase: useCase,
	}
	c.GET("/sequence/v0.1/game/rank", handler.SequenceRank, mw.TokenChecker)
	return handler
}

// 시퀀스 랭크 가져오기
// @Router /sequence/v0.1/game/rank [get]
// @Summary 시퀀스 랭크 가져오기
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
// @Success 200 {object} response.ResSequenceRank
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags app/sequence/game
func (d *SequenceRankBoardGameHandler) SequenceRank(c echo.Context) error {
	ctx, _, _ := utils.CtxGenerate(c)

	res, err := d.UseCase.SequenceRank(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
