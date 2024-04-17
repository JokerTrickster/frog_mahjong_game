package handler

import (
	_interface "main/features/room/model/interface"
	"main/features/room/model/request"
	mw "main/middleware"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserListRoomHandler struct {
	UseCase _interface.IUserListRoomUseCase
}

func NewUserListRoomHandler(c *echo.Echo, useCase _interface.IUserListRoomUseCase) _interface.IUserListRoomHandler {
	handler := &UserListRoomHandler{
		UseCase: useCase,
	}
	c.GET("/v0.1/room/user", handler.UserList, mw.TokenChecker)
	return handler
}

// 게임 준비 상태 변경
// @Router /v0.1/room/user [get]
// @Summary 게임 준비 상태 변경
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
// @Param json body request.ReqUserList true "json body"
// @Produce json
// @Success 200 {object} response.ResUserListRoom
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags room
func (d *UserListRoomHandler) UserList(c echo.Context) error {
	ctx, _, _ := utils.CtxGenerate(c)
	req := &request.ReqUserList{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}
	res, err := d.UseCase.UserList(ctx, req.RoomID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}
