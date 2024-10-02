package handler

import (
	_interface "main/features/users/model/interface"
	mw "main/middleware"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ListProfilesUsersHandler struct {
	UseCase _interface.IListProfilesUsersUseCase
}

func NewListProfilesUsersHandler(c *echo.Echo, useCase _interface.IListProfilesUsersUseCase) _interface.IListProfilesUsersHandler {
	handler := &ListProfilesUsersHandler{
		UseCase: useCase,
	}
	c.GET("/v0.1/users/profiles", handler.ListProfiles, mw.TokenChecker)
	return handler
}

// 유저 프로필 리스트 가져오기
// @Router /v0.1/users/profiles [get]
// @Summary 유저 프로필 리스트 가져오기
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
// @Success 200 {object} response.ResListProfileUser
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags user
func (d *ListProfilesUsersHandler) ListProfiles(c echo.Context) error {
	ctx, userID, _ := utils.CtxGenerate(c)
	res, err := d.UseCase.ListProfiles(ctx, userID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}
