package handler

import (
	_interface "main/features/room/model/interface"
	"main/features/room/model/request"
	mw "main/middleware"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ListRoomHandler struct {
	UseCase _interface.IListRoomUseCase
}

func NewListRoomHandler(c *echo.Echo, useCase _interface.IListRoomUseCase) _interface.IListRoomHandler {
	handler := &ListRoomHandler{
		UseCase: useCase,
	}
	c.GET("/v0.1/room", handler.List, mw.TokenChecker)
	return handler
}

// 방 목록 가져오기
// @Router /v0.1/room [get]
// @Summary 방 목록 가져오기
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description USER_NOT_FOUND : 유저를 찾을 수 없음
// @Description USER_ALREADY_EXISTED : 이미 존재하는 유저
// @Description ROOM_NOT_FOUND : 방을 찾을 수 없음
// @Description ROOM_FULL : 방이 꽉 참
// @Description ROOM_USER_NOT_FOUND : 방 유저를 찾을 수 없음
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Description PLAYER_STATE_CHANGE_FAILED : 플레이어 상태 변경 실패
// @Param tkn header string true "accessToken"
// @Param json body request.ReqList true "json body"
// @Produce json
// @Success 200 {object} boolean
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags room
func (d *ListRoomHandler) List(c echo.Context) error {
	ctx, _, _ := utils.CtxGenerate(c)
	req := &request.ReqList{}
	if err := utils.ValidateReq(c, req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err := d.UseCase.List(ctx, req.Page, req.PageSize)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, true)
}
