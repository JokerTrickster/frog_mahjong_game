package handler

import (
	_interface "main/features/board_game/model/interface"
	"main/features/board_game/model/request"
	mw "main/middleware"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type FindItSoloPlayBoardGameHandler struct {
	UseCase _interface.IFindItSoloPlayBoardGameUseCase
}

func NewFindItSoloPlayBoardGameHandler(c *echo.Echo, useCase _interface.IFindItSoloPlayBoardGameUseCase) _interface.IFindItSoloPlayBoardGameHandler {
	handler := &FindItSoloPlayBoardGameHandler{
		UseCase: useCase,
	}
	c.POST("/find-it/v0.1/game/solo-play", handler.FindItSoloPlay, mw.TokenChecker)
	return handler
}

// 틀린그림찾기 솔로플레이 이미지 가져오기
// @Router /find-it/v0.1/game/solo-play [post]
// @Summary 틀린그림찾기 솔로플레이 이미지 가져오기
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description USER_NOT_FOUND : 유저를 찾을 수 없음
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Param tkn header string true "accessToken"
// @Param json body request.ReqFindItSoloPlayBoardGame true "플레이 라운드 수"
// @Produce json
// @Success 200 {object} response.ResFindItSoloPlayBoardGame
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags app/find-it/game
func (d *FindItSoloPlayBoardGameHandler) FindItSoloPlay(c echo.Context) error {
	ctx, userID, _ := utils.CtxGenerate(c)
	req := &request.ReqFindItSoloPlayBoardGame{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}
	res, err := d.UseCase.FindItSoloPlay(ctx, int(userID), req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}
