package handler

import (
	_interface "main/features/game_users/model/interface"
	mw "main/middleware"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type DeleteUsersHandler struct {
	UseCase _interface.IDeleteUsersUseCase
}

func NewDeleteUsersHandler(c *echo.Echo, useCase _interface.IDeleteUsersUseCase) _interface.IDeleteUsersHandler {
	handler := &DeleteUsersHandler{
		UseCase: useCase,
	}
	c.DELETE("/v0.1/game/users", handler.Delete, mw.TokenChecker)
	return handler
}

// 회원 탈퇴하기
// @Router /v0.1/game/users [delete]
// @Summary 회원 탈퇴하기
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description USER_NOT_FOUND : 유저를 찾을 수 없음
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Param tkn header string true "accessToken"
// @Produce json
// @Success 200 {object} bool
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags game/user
func (d *DeleteUsersHandler) Delete(c echo.Context) error {
	ctx, userID, _ := utils.CtxGenerate(c)

	err := d.UseCase.Delete(ctx, userID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, true)
}
