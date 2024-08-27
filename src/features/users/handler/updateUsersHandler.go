package handler

import (
	_interface "main/features/users/model/interface"
	"main/features/users/model/request"
	mw "main/middleware"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UpdateUsersHandler struct {
	UseCase _interface.IUpdateUsersUseCase
}

func NewUpdateUsersHandler(c *echo.Echo, useCase _interface.IUpdateUsersUseCase) _interface.IUpdateUsersHandler {
	handler := &UpdateUsersHandler{
		UseCase: useCase,
	}
	c.PUT("/v0.1/users", handler.Update, mw.TokenChecker)
	return handler
}

// 유저 정보 수정하기
// @Router /v0.1/users [put]
// @Summary 유저 정보 수정하기
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description USER_NOT_FOUND : 유저를 찾을 수 없음
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Param tkn header string true "accessToken"
// @Param json body request.ReqUpdateUsers true "json body"
// @Produce json
// @Success 200 {object} bool
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags user
func (d *UpdateUsersHandler) Update(c echo.Context) error {
	ctx, userID, _ := utils.CtxGenerate(c)
	req := &request.ReqUpdateUsers{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}
	err := d.UseCase.Update(ctx, userID, req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, true)
}
