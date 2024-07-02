package handler

import (
	_interface "main/features/chat/model/interface"
	mw "main/middleware"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthChatHandler struct {
	UseCase _interface.IAuthChatUseCase
}

func NewAuthChatHandler(c *echo.Echo, useCase _interface.IAuthChatUseCase) _interface.IAuthChatHandler {
	handler := &AuthChatHandler{
		UseCase: useCase,
	}
	c.GET("/v0.1/chat/auth", handler.Auth, mw.TokenChecker)
	return handler
}

// 챗 유저 검증
// @Router /v0.1/chat/auth [get]
// @Summary 챗 유저 검증
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description NOT_OWNER : 방장이 시작 요청을 하지 않음
// @Description NOT_FIRST_PLAYER : 첫 플레이어가 아님
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Param tkn header string true "accessToken"
// @Produce json
// @Success 200 {object} bool
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags chat
func (d *AuthChatHandler) Auth(c echo.Context) error {

	ctx, userID, _ := utils.CtxGenerate(c)

	// business logic
	_, err := d.UseCase.Auth(ctx, userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, true)
}
