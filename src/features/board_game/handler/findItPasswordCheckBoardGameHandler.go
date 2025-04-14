package handler

import (
	_interface "main/features/board_game/model/interface"
	"main/features/board_game/model/request"
	mw "main/middleware"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type FindItPasswordCheckBoardGameHandler struct {
	UseCase _interface.IFindItPasswordCheckBoardGameUseCase
}

func NewFindItPasswordCheckBoardGameHandler(c *echo.Echo, useCase _interface.IFindItPasswordCheckBoardGameUseCase) _interface.IFindItPasswordCheckBoardGameHandler {
	handler := &FindItPasswordCheckBoardGameHandler{
		UseCase: useCase,
	}
	c.POST("/find-it/v0.1/game/join/password-check", handler.FindItPasswordCheck, mw.TokenChecker)
	return handler
}

// 틀린그림찾기 비밀번호 확인
// @Router /find-it/v0.1/game/join/password-check [post]
// @Summary 틀린그림찾기 비밀번호 확인
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description USER_NOT_FOUND : 유저를 찾을 수 없음
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Param tkn header string true "accessToken"
// @Param json body request.ReqFindItPasswordCheckBoardGame true "인증 코드 "
// @Produce json
// @Success 200 {object} bool
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags app/find-it/game
func (d *FindItPasswordCheckBoardGameHandler) FindItPasswordCheck(c echo.Context) error {
	ctx, _, _ := utils.CtxGenerate(c)
	req := &request.ReqFindItPasswordCheckBoardGame{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}
	res, err := d.UseCase.FindItPasswordCheck(ctx, req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}
