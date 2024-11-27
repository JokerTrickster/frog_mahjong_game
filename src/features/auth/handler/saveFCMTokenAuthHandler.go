package handler

import (
	_interface "main/features/auth/model/interface"
	"main/features/auth/model/request"
	mw "main/middleware"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type FCMTokenAuthHandler struct {
	UseCase _interface.IFCMTokenAuthUseCase
}

func NewFCMTokenAuthHandler(c *echo.Echo, useCase _interface.IFCMTokenAuthUseCase) _interface.IFCMTokenAuthHandler {
	handler := &FCMTokenAuthHandler{
		UseCase: useCase,
	}
	c.POST("/v0.1/auth/fcm-token", handler.FCMToken, mw.TokenChecker)
	return handler
}

// fcm 토큰 저장
// @Router /v0.1/auth/fcm-token [post]
// @Summary Fcm 토큰 저장
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description USER_NOT_EXIST : 유저가 존재하지 않음
// @Description USER_ALREADY_EXISTED : 유저가 이미 존재
// @Description USER_GOOGLE_ALREADY_EXISTED : 구글 계정이 이미 존재
// @Description PASSWORD_NOT_MATCH : 비밀번호가 일치하지 않음
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Param tkn header string true "accessToken"
// @Param json body request.ReqFCMToken true "이메일, 비밀번호"
// @Produce json
// @Success 200 {object} bool
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags auth
func (d *FCMTokenAuthHandler) FCMToken(c echo.Context) error {
	ctx, uID, _ := utils.CtxGenerate(c)

	req := &request.ReqFCMToken{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}
	err := d.UseCase.FCMToken(ctx, uID, req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, true)
}
