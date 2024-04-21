package handler

import (
	_interface "main/features/rooms/model/interface"
	"main/features/rooms/model/request"
	mw "main/middleware"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CreateRoomsHandler struct {
	UseCase _interface.ICreateRoomsUseCase
}

func NewCreateRoomsHandler(c *echo.Echo, useCase _interface.ICreateRoomsUseCase) _interface.ICreateRoomsHandler {
	handler := &CreateRoomsHandler{
		UseCase: useCase,
	}
	c.POST("/v0.1/rooms/create", handler.Create, mw.TokenChecker)
	return handler
}

// 방 생성
// @Router /v0.1/rooms/create [post]
// @Summary 방 생성
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
// @Param json body request.ReqCreate true "json body"
// @Produce json
// @Success 200 {object} response.ResCreateRoom
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags Rooms
func (d *CreateRoomsHandler) Create(c echo.Context) error {
	ctx, uID, email := utils.CtxGenerate(c)
	req := &request.ReqCreate{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}
	res, err := d.UseCase.Create(ctx, uID, email, req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, res)
}
