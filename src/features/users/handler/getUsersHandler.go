package handler

import (
	_interface "main/features/users/model/interface"
	"main/features/users/model/request"
	mw "main/middleware"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type GetUsersHandler struct {
	UseCase _interface.IGetUsersUseCase
}

func NewGetUsersHandler(c *echo.Echo, useCase _interface.IGetUsersUseCase) _interface.IGetUsersHandler {
	handler := &GetUsersHandler{
		UseCase: useCase,
	}
	c.GET("/v0.1/users/:userID", handler.Get, mw.TokenChecker)
	return handler
}

// 유저 정보 가져오기
// @Router /v0.1/users/{userID} [get]
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
// @Param userID query string true "유저 아이디"
// @Produce json
// @Success 200 {object} response.ResGetUser
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags user
func (d *GetUsersHandler) Get(c echo.Context) error {
	ctx, _, _ := utils.CtxGenerate(c)
	req := &request.ReqGetUser{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}
	res, err := d.UseCase.Get(ctx, req.UserID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, res)
}
