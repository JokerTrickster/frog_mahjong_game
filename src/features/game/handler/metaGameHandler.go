package handler

import (
	"main/utils"
	"net/http"

	_interface "main/features/game/model/interface"

	"github.com/labstack/echo/v4"
)

type MetaGameHandler struct {
	UseCase _interface.IMetaGameUseCase
}

func NewMetaGameHandler(c *echo.Echo, useCase _interface.IMetaGameUseCase) _interface.IMetaGameHandler {
	handler := &MetaGameHandler{
		UseCase: useCase,
	}
	c.GET("/v2.1/game/report/meta", handler.Meta)
	return handler
}

// 신고하기 카테고리 데이터 가져오기
// @Router /v2.1/game/report/meta [get]
// @Summary 신고하기 카테고리 데이터 가져오기
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description NOT_ALL_USERS_READY : 모든 유저가 준비되지 않음
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Description
// @Produce json
// @Success 200 {object} response.ResMetaGame
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags game
func (d *MetaGameHandler) Meta(c echo.Context) error {
	ctx, _, _ := utils.CtxGenerate(c)
	res, err := d.UseCase.Meta(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
