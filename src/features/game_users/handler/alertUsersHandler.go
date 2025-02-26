package handler

import (
	_interface "main/features/game_users/model/interface"
	"main/features/game_users/model/request"
	mw "main/middleware"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AlertUsersHandler struct {
	UseCase _interface.IAlertUsersUseCase
}

func NewAlertUsersHandler(c *echo.Echo, useCase _interface.IAlertUsersUseCase) _interface.IAlertUsersHandler {
	handler := &AlertUsersHandler{
		UseCase: useCase,
	}
	c.PUT("/v0.1/game/users/alert", handler.Alert, mw.TokenChecker)
	return handler
}

// 유저 알람 활성화/비활성화 수정하기
// @Router /v0.1/game/users/alert [put]
// @Summary 유저 알람 활성화/비활성화 수정하기
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description USER_NOT_FOUND : 유저를 찾을 수 없음
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Param tkn header string true "accessToken"
// @Param json body request.ReqAlertGameUsers true "활성화 (true) / 비활성화 (false)"
// @Produce json
// @Success 200 {object} bool
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags game/user
func (d *AlertUsersHandler) Alert(c echo.Context) error {
	ctx, userID, _ := utils.CtxGenerate(c)
	req := &request.ReqAlertGameUsers{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}
	err := d.UseCase.Alert(ctx, userID, req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, true)
}
