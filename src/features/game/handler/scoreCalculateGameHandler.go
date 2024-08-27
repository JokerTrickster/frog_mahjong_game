package handler

import (
	mw "main/middleware"
	"main/utils"
	"net/http"

	_interface "main/features/game/model/interface"
	"main/features/game/model/response"

	"main/features/game/model/request"

	"github.com/labstack/echo/v4"
)

type ScoreCalculateGameHandler struct {
	UseCase _interface.IScoreCalculateGameUseCase
}

func NewScoreCalculateGameHandler(c *echo.Echo, useCase _interface.IScoreCalculateGameUseCase) _interface.IScoreCalculateGameHandler {
	handler := &ScoreCalculateGameHandler{
		UseCase: useCase,
	}
	c.POST("/v0.1/game/score/calculate", handler.ScoreCalculate, mw.TokenChecker)
	return handler
}

// 점수 계산하기
// @Router /v0.1/game/score/calculate [post]
// @Summary 점수 계산하기
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description NOT_ALL_USERS_READY : 모든 유저가 준비되지 않음
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Description
// @Description ■ bonus list
// @Description same : 같은 패 (2점)
// @Description continuous : 연속 패 (1점)
// @Description allGreen : 올 그린 (10점)
// @Description allGreen : 도라 (하나당 1점)
// @Description allGreen : 적패 (하나당 1점)
// @Description allGreen : 올 그린 (10점)
// @Description superRed : 슈퍼 레드 (20점)
// @Description tangYao :  탕야오 (1점)
// @Description chanTa : 찬타 (2점)
// @Description chinYao : 칭야오 (15점)
// @Param tkn header string true "accessToken"
// @Param json body request.ReqScoreCalculate true "json body"
// @Produce json
// @Success 200 {object} response.ResScoreCalculate
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags game
func (d *ScoreCalculateGameHandler) ScoreCalculate(c echo.Context) error {
	ctx, userID, _ := utils.CtxGenerate(c)
	req := &request.ReqScoreCalculate{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}
	score, bonuses, err := d.UseCase.ScoreCalculate(ctx, userID, req)
	if err != nil {
		return err
	}
	res := response.ResScoreCalculate{
		Score:   score,
		Bonuses: bonuses,
	}
	return c.JSON(http.StatusOK, res)
}
