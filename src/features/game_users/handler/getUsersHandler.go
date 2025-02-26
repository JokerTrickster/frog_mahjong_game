package handler

import (
	_interface "main/features/game_users/model/interface"
	mw "main/middleware"
	"main/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type GetUsersHandler struct {
	UseCase _interface.IGetUsersUseCase
}

func NewGetUsersHandler(c *echo.Echo, useCase _interface.IGetUsersUseCase) _interface.IGetUsersHandler {
	handler := &GetUsersHandler{
		UseCase: useCase,
	}
	c.GET("/v0.1/game/users/:userID", handler.Get, mw.TokenChecker)
	return handler
}

// 유저 정보 가져오기
// @Router /v0.1/game/users/{userID} [get]
// @Summary 유저 정보 가져오기
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description USER_NOT_FOUND : 유저를 찾을 수 없음
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Param tkn header string true "accessToken"
// @Param userID path string true "userID"
// @Produce json
// @Success 200 {object} response.ResGetGameUser
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags game/user
func (d *GetUsersHandler) Get(c echo.Context) error {
	ctx, uID, _ := utils.CtxGenerate(c)
	pathUserID := c.Param("userID")
	puID, _ := strconv.Atoi(pathUserID)
	if pathUserID == "" || uID != uint(puID) {
		return utils.ErrorMsg(ctx, utils.ErrBadParameter, utils.Trace(), utils.HandleError("invalid user id", pathUserID), utils.ErrFromClient)
	}

	res, err := d.UseCase.Get(ctx, int(uID))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}
