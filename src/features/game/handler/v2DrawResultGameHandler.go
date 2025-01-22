package handler

import (
	mw "main/middleware"
	"main/utils"
	"net/http"
	"strconv"

	_interface "main/features/game/model/interface"

	"github.com/labstack/echo/v4"
)

type V2DrawResultGameHandler struct {
	UseCase _interface.IV2DrawResultGameUseCase
}

func NewV2DrawResultGameHandler(c *echo.Echo, useCase _interface.IV2DrawResultGameUseCase) _interface.IV2DrawResultGameHandler {
	handler := &V2DrawResultGameHandler{
		UseCase: useCase,
	}
	c.GET("/v2.1/game/:roomID/draw", handler.V2DrawResult, mw.TokenChecker)
	return handler
}

// 무승부 게임 결과 가져오기
// @Router /v2.1/game/{roomID}/draw [get]
// @Summary 무승부 게임 결과 가져오기
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description NOT_ALL_USERS_READY : 모든 유저가 준비되지 않음
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Param tkn header string true "accessToken"
// @Param roomID path string true "roomID"
// @Produce json
// @Success 200 {object} response.ResV2DrawResult
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags game
func (d *V2DrawResultGameHandler) V2DrawResult(c echo.Context) error {
	ctx, _, _ := utils.CtxGenerate(c)
	roomID := c.Param("roomID")
	rID, _ := strconv.Atoi(roomID)

	res, err := d.UseCase.V2DrawResult(ctx, rID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
