package handler

import (
	"context"
	_interface "main/features/users/model/interface"
	"net/http"

	"github.com/labstack/echo/v4"
)

type OneCoinUsersHandler struct {
	UseCase _interface.IOneCoinUsersUseCase
}

func NewOneCoinUsersHandler(c *echo.Echo, useCase _interface.IOneCoinUsersUseCase) _interface.IOneCoinUsersHandler {
	handler := &OneCoinUsersHandler{
		UseCase: useCase,
	}
	c.POST("/v2.1/users/batch/coins/one", handler.OneCoin)
	return handler
}

// [배치용] 유저 코인 1씩 회복
// @Router /v2.1/users/batch/coins/one [post]
// @Summary [배치용] 유저 코인 1씩 회복
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description USER_NOT_FOUND : 유저를 찾을 수 없음
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Produce json
// @Success 200 {object} bool
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags user
func (d *OneCoinUsersHandler) OneCoin(c echo.Context) error {
	ctx := context.Background()
	err := d.UseCase.OneCoin(ctx)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, true)
}
