package handler

import (
	_interface "main/features/chat/model/interface"
	"main/features/chat/model/request"
	mw "main/middleware"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type HistoryChatHandler struct {
	UseCase _interface.IHistoryChatUseCase
}

func NewHistoryChatHandler(c *echo.Echo, useCase _interface.IHistoryChatUseCase) _interface.IHistoryChatHandler {
	handler := &HistoryChatHandler{
		UseCase: useCase,
	}
	c.GET("/v0.1/chats", handler.History, mw.TokenChecker)
	return handler
}

// 채팅 리스트 가져오기
// @Router /v0.1/chats [get]
// @Summary 채팅 리스트 가져오기
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description NOT_OWNER : 방장이 시작 요청을 하지 않음
// @Description NOT_FIRST_PLAYER : 첫 플레이어가 아님
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Param tkn header string true "accessToken"
// @Param roomID query int true "roomID"
// @Param page query int false "조회할 페이지. 0부터 시작, 누락시 0으로 처리"
// @Param pageSize query int false "페이지당 알림 개수. 누락시 10으로 처리 "
// @Produce json
// @Success 200 {object} response.ResHistoryChat
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags chats
func (d *HistoryChatHandler) History(c echo.Context) error {
	ctx, _, _ := utils.CtxGenerate(c)
	req := &request.ReqHistory{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	res, err := d.UseCase.History(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
