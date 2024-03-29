package handler

import (
	_interface "main/features/room/model/interface"
	"main/features/room/model/request"
	mw "main/middleware"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type OutRoomHandler struct {
	UseCase _interface.IOutRoomUseCase
}

func NewOutRoomHandler(c *echo.Echo, useCase _interface.IOutRoomUseCase) _interface.IOutRoomHandler {
	handler := &OutRoomHandler{
		UseCase: useCase,
	}
	c.POST("/v0.1/room/out", handler.Out, mw.TokenChecker)
	return handler
}

// 방 나가기
// @Router /v0.1/room/out [post]
// @Summary 방 나가기
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Param tkn header string true "accessToken"
// @Param json body request.ReqOut true "json body"
// @Produce json
// @Success 200 {object} boolean
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags room
func (d *OutRoomHandler) Out(c echo.Context) error {
	ctx, uID, _ := utils.CtxGenerate(c)
	req := &request.ReqOut{}
	if err := utils.ValidateReq(c, req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err := d.UseCase.Out(ctx, uID, req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, true)
}
