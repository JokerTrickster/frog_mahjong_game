package handler

import (
	"context"
	_interface "main/features/game_users/model/interface"
	"net/http"

	"github.com/labstack/echo/v4"
)

type FullCoinUsersHandler struct {
	UseCase _interface.IFullCoinUsersUseCase
}

func NewFullCoinUsersHandler(c *echo.Echo, useCase _interface.IFullCoinUsersUseCase) _interface.IFullCoinUsersHandler {
	handler := &FullCoinUsersHandler{
		UseCase: useCase,
	}
	c.POST("/v0.1/game/users/batch/coins/full", handler.FullCoin)
	return handler
}

// [배치용] 유저 코인 30까지 모두 회복
// @Router /v0.1/game/users/batch/coins/full [post]
// @Summary [배치용] 유저 코인 30까지 모두 회복
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
// @Tags game/user
func (d *FullCoinUsersHandler) FullCoin(c echo.Context) error {
	ctx := context.Background()
	err := d.UseCase.FullCoin(ctx)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, true)
}
