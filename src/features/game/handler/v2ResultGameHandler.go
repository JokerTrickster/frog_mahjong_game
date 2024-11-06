package handler

import (
	mw "main/middleware"
	"main/utils"
	"net/http"

	_interface "main/features/game/model/interface"

	"main/features/game/model/request"

	"github.com/labstack/echo/v4"
)

type V2ResultGameHandler struct {
	UseCase _interface.IV2ResultGameUseCase
}

func NewV2ResultGameHandler(c *echo.Echo, useCase _interface.IV2ResultGameUseCase) _interface.IV2ResultGameHandler {
	handler := &V2ResultGameHandler{
		UseCase: useCase,
	}
	c.POST("/v2.1/game/result", handler.V2Result, mw.TokenChecker)
	return handler
}

// 게임 결과 가져오기
// @Router /v2.1/game/result [post]
// @Summary 게임 결과 가져오기
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description NOT_ALL_USERS_READY : 모든 유저가 준비되지 않음
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Param tkn header string true "accessToken"
// @Param json body request.ReqV2Result true "json body"
// @Produce json
// @Success 200 {object} response.ResV2Result
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags game
func (d *V2ResultGameHandler) V2Result(c echo.Context) error {
	ctx, _, _ := utils.CtxGenerate(c)
	req := &request.ReqV2Result{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}
	res, err := d.UseCase.V2Result(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
