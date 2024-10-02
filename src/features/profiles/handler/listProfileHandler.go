package handler

import (
	"context"
	_interface "main/features/profiles/model/interface"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ListProfilesHandler struct {
	UseCase _interface.IListProfilesUseCase
}

func NewListProfilesHandler(c *echo.Echo, useCase _interface.IListProfilesUseCase) _interface.IListProfilesHandler {
	handler := &ListProfilesHandler{
		UseCase: useCase,
	}
	c.GET("/v0.1/profiles", handler.List)
	return handler
}

// 프로필 리스트 가져오기
// @Router /v0.1/profiles [get]
// @Summary 프로필 리스트 가져오기
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description USER_NOT_FOUND : 유저를 찾을 수 없음
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Produce json
// @Success 200 {object} response.ResListProfile
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags profile
func (d *ListProfilesHandler) List(c echo.Context) error {
	ctx := context.Background()
	res, err := d.UseCase.List(ctx)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}
