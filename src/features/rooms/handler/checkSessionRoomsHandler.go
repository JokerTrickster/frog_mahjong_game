package handler

import (
	_interface "main/features/rooms/model/interface"
	"main/features/rooms/model/request"
	mw "main/middleware"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CheckSessionRoomsHandler struct {
	UseCase _interface.ICheckSessionRoomsUseCase
}

func NewCheckSessionRoomsHandler(c *echo.Echo, useCase _interface.ICheckSessionRoomsUseCase) _interface.ICheckSessionRoomsHandler {
	handler := &CheckSessionRoomsHandler{
		UseCase: useCase,
	}
	c.POST("/v0.1/rooms/session/check", handler.CheckSession, mw.TokenChecker)
	return handler
}

// 세션 ID 체크
// @Router /v0.1/rooms/session/check [post]
// @Summary 세션 ID 체크
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description USER_NOT_FOUND : 유저를 찾을 수 없음
// @Description USER_ALREADY_EXISTED : 이미 존재하는 유저
// @Description ROOM_NOT_FOUND : 방을 찾을 수 없음
// @Description ROOM_FULL : 방이 꽉 찼습니다.
// @Description ROOM_USER_NOT_FOUND : 방 유저를 찾을 수 없음
// @Description WRONG_PASSWORD : 비밀번호가 일치하지 않습니다.
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Description PLAYER_STATE_CHANGE_FAILED : 플레이어 상태 변경 실패
// @Param tkn header string true "accessToken"
// @Param json body request.ReqCheckSession true "json body"
// @Produce json
// @Success 200 {object} bool
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags rooms
func (d *CheckSessionRoomsHandler) CheckSession(c echo.Context) error {
	ctx, _, _ := utils.CtxGenerate(c)
	req := &request.ReqCheckSession{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}
	isExisted, err := d.UseCase.CheckSession(ctx, req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, isExisted)
}
