package handler

import (
	mw "main/middleware"
	"main/utils"
	"net/http"

	_interface "main/features/board_game/model/interface"
	"main/features/board_game/model/request"

	"github.com/labstack/echo/v4"
)

type SlimeWarResultBoardGameHandler struct {
	UseCase _interface.ISlimeWarResultBoardGameUseCase
}

func NewSlimeWarResultBoardGameHandler(c *echo.Echo, useCase _interface.ISlimeWarResultBoardGameUseCase) _interface.ISlimeWarResultBoardGameHandler {
	handler := &SlimeWarResultBoardGameHandler{
		UseCase: useCase,
	}
	c.POST("/slime-war/v0.1/game/result", handler.SlimeWarResult, mw.TokenChecker)
	return handler
}

// [슬라임 전쟁] 게임 결과 가져오기
// @Router /slime-war/v0.1/game/result [post]
// @Summary [슬라임 전쟁] 게임 결과 가져오기
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
// @Param json body request.ReqSlimeWarResult true "json body"
// @Produce json
// @Success 200 {object} response.ResSlimeWarResult
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags app/slime-war/game
func (d *SlimeWarResultBoardGameHandler) SlimeWarResult(c echo.Context) error {
	ctx, _, _ := utils.CtxGenerate(c)
	req := &request.ReqSlimeWarResult{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}
	res, err := d.UseCase.SlimeWarResult(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
