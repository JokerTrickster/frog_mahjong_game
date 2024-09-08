package handler

import (
	"context"
	_interface "main/features/rooms/model/interface"
	"net/http"

	"github.com/labstack/echo/v4"
)

type MetaRoomsHandler struct {
	UseCase _interface.IMetaRoomsUseCase
}

func NewMetaRoomsHandler(c *echo.Echo, useCase _interface.IMetaRoomsUseCase) _interface.IMetaRoomsHandler {
	handler := &MetaRoomsHandler{
		UseCase: useCase,
	}
	c.GET("/v0.1/rooms/meta", handler.Meta)
	return handler
}

// 방 메타 데이터 가져오기
// @Router /v0.1/rooms/meta [get]
// @Summary 방 데이터 가져오기
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description USER_NOT_FOUND : 유저를 찾을 수 없음
// @Description USER_ALREADY_EXISTED : 이미 존재하는 유저
// @Description Room_NOT_FOUND : 방을 찾을 수 없음
// @Description Room_FULL : 방이 꽉 참
// @Description Room_USER_NOT_FOUND : 방 유저를 찾을 수 없음
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Description PLAYER_STATE_CHANGE_FAILED : 플레이어 상태 변경 실패
// @Produce json
// @Success 200 {object} response.ResMetaRoom
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags rooms
func (d *MetaRoomsHandler) Meta(c echo.Context) error {
	ctx := context.Background()
	res, err := d.UseCase.Meta(ctx)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}
