package handler

import (
	mw "main/middleware"
	"main/utils"
	"net/http"

	_interface "main/features/game/model/interface"

	"main/features/game/model/request"

	"github.com/labstack/echo/v4"
)

type ReportGameHandler struct {
	UseCase _interface.IReportGameUseCase
}

func NewReportGameHandler(c *echo.Echo, useCase _interface.IReportGameUseCase) _interface.IReportGameHandler {
	handler := &ReportGameHandler{
		UseCase: useCase,
	}
	c.POST("/v0.1/game/report", handler.Report, mw.TokenChecker)
	return handler
}

// 유저 신고하기
// @Router /v0.1/game/report [post]
// @Summary 유저 신고하기
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description NOT_ALL_USERS_READY : 모든 유저가 준비되지 않음
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Description
// @Param tkn header string true "accessToken"
// @Param json body request.ReqReport true "json body"
// @Produce json
// @Success 200 {object} bool
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags game
func (d *ReportGameHandler) Report(c echo.Context) error {
	ctx, userID, _ := utils.CtxGenerate(c)
	req := &request.ReqReport{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}
	err := d.UseCase.Report(ctx, userID, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, true)
}
