package handler

import (
	_interface "main/features/rooms/model/interface"
	"main/features/rooms/model/request"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ListRoomsHandler struct {
	UseCase _interface.IListRoomsUseCase
}

func NewListRoomsHandler(c *echo.Echo, useCase _interface.IListRoomsUseCase) _interface.IListRoomsHandler {
	handler := &ListRoomsHandler{
		UseCase: useCase,
	}
	c.GET("/v0.1/rooms", handler.List)
	return handler
}

// 방 목록 가져오기
// @Router /v0.1/rooms [get]
// @Summary 방 목록 가져오기
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
// @Param page query int false "조회할 페이지. 0부터 시작, 누락시 0으로 처리"
// @Param pageSize query int false "페이지당 알림 개수. 누락시 10으로 처리 "
// @Produce json
// @Success 200 {object} response.ResListRoom
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags rooms
func (d *ListRoomsHandler) List(c echo.Context) error {
	ctx, _, _ := utils.CtxGenerate(c)
	req := &request.ReqList{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	res, err := d.UseCase.List(ctx, req.Page, req.PageSize)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}
