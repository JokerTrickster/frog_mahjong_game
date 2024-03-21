package handler

import (
	_interface "main/features/room/model/interface"
	"main/features/room/model/request"
	mw "main/middleware"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CreateRoomHandler struct {
	UseCase _interface.ICreateRoomUseCase
}

func NewCreateRoomHandler(c *echo.Echo, useCase _interface.ICreateRoomUseCase) _interface.ICreateRoomHandler {
	handler := &CreateRoomHandler{
		UseCase: useCase,
	}
	c.POST("/v0.1/room/create", handler.Create, mw.TokenChecker)
	return handler
}

// 방 생성
// @Router /v0.1/room/create [post]
// @Summary 방 생성
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Param tkn header string true "accessToken"
// @Param json body request.ReqCreate true "json body"
// @Produce json
// @Success 200 {object} boolean
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags room
func (d *CreateRoomHandler) Create(c echo.Context) error {
	ctx, uID, email := utils.CtxGenerate(c)
	req := &request.ReqCreate{}
	if err := utils.ValidateReq(c, req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err := d.UseCase.Create(ctx, uID, email, req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, true)
}
