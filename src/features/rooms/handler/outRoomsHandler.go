package handler

import (
	_interface "main/features/rooms/model/interface"
	"main/features/rooms/model/request"
	mw "main/middleware"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type OutRoomsHandler struct {
	UseCase _interface.IOutRoomsUseCase
}

func NewOutRoomsHandler(c *echo.Echo, useCase _interface.IOutRoomsUseCase) _interface.IOutRoomsHandler {
	handler := &OutRoomsHandler{
		UseCase: useCase,
	}
	c.PUT("/v0.1/rooms/out", handler.Out, mw.TokenChecker)
	return handler
}

// 방 나가기
// @Router /v0.1/rooms/out [put]
// @Summary 방 나가기
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
// @Param json body request.ReqOut true "json body"
// @Produce json
// @Success 200 {object} boolean
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags rooms
func (d *OutRoomsHandler) Out(c echo.Context) error {
	ctx, uID, _ := utils.CtxGenerate(c)
	req := &request.ReqOut{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}
	err := d.UseCase.Out(ctx, uID, req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, true)
}
