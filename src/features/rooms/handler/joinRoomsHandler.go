package handler

import (
	_interface "main/features/rooms/model/interface"
	"main/features/rooms/model/request"
	mw "main/middleware"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type JoinRoomsHandler struct {
	UseCase _interface.IJoinRoomsUseCase
}

func NewJoinRoomsHandler(c *echo.Echo, useCase _interface.IJoinRoomsUseCase) _interface.IJoinRoomsHandler {
	handler := &JoinRoomsHandler{
		UseCase: useCase,
	}
	c.POST("/v0.1/rooms/join", handler.Join, mw.TokenChecker)
	return handler
}

// 방 참여
// @Router /v0.1/rooms/join [post]
// @Summary 방 참여
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
// @Param tkn header string true "accessToken"
// @Param json body request.ReqJoin true "json body"
// @Produce json
// @Success 200 {object} response.ResJoinRoom
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags Rooms
func (d *JoinRoomsHandler) Join(c echo.Context) error {
	ctx, uID, email := utils.CtxGenerate(c)
	req := &request.ReqJoin{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}
	res, err := d.UseCase.Join(ctx, uID, email, req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}