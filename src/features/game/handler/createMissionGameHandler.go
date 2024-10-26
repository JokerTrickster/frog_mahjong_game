package handler

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CreateMissionGameHandler struct {
	UseCase _interface.ICreateMissionGameUseCase
}

func NewCreateMissionGameHandler(c *echo.Echo, useCase _interface.ICreateMissionGameUseCase) _interface.ICreateMissionGameHandler {
	handler := &CreateMissionGameHandler{
		UseCase: useCase,
	}
	c.POST("/v2.1/game/missions", handler.CreateMission)
	return handler
}

// 미션 생성하기
// @Router /v2.1/game/missions [post]
// @Summary 미션 생성하기
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description NOT_OWNER : 방장이 시작 요청을 하지 않음
// @Description NOT_FIRST_PLAYER : 첫 플레이어가 아님
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Produce json
// @Param json body request.ReqCreateMission true "json body"
// @Success 200 {object} bool
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags game
func (d *CreateMissionGameHandler) CreateMission(c echo.Context) error {

	ctx := context.Background()
	//business logic
	req := &request.ReqCreateMission{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}
	err := d.UseCase.CreateMission(ctx, req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, true)
}
