package handler

import (
	mw "main/middleware"
	"main/utils"
	"net/http"

	_interface "main/features/board_game/model/interface"
	"main/features/board_game/model/request"

	"github.com/labstack/echo/v4"
)

type SequenceResultBoardGameHandler struct {
	UseCase _interface.ISequenceResultBoardGameUseCase
}

func NewSequenceResultBoardGameHandler(c *echo.Echo, useCase _interface.ISequenceResultBoardGameUseCase) _interface.ISequenceResultBoardGameHandler {
	handler := &SequenceResultBoardGameHandler{
		UseCase: useCase,
	}
	c.POST("/sequence/v0.1/game/result", handler.SequenceResult, mw.TokenChecker)
	return handler
}

// [시퀀스] 게임 결과 가져오기
// @Router /sequence/v0.1/game/result [post]
// @Summary [시퀀스] 게임 결과 가져오기
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description NOT_ALL_USERS_READY : 모든 유저가 준비되지 않음
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Description
// @Param tkn header string true "accessToken"
// @Param json body request.ReqSequenceResult true "json body"
// @Produce json
// @Success 200 {object} response.ResSequenceResult
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags app/sequence/game
func (d *SequenceResultBoardGameHandler) SequenceResult(c echo.Context) error {
	ctx, _, _ := utils.CtxGenerate(c)
	req := &request.ReqSequenceResult{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}
	res, err := d.UseCase.SequenceResult(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
