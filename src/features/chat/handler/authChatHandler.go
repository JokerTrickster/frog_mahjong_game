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

func (d *AuthChatHandler) Auth(c echo.Context) error {

	ctx, userID, _ := utils.CtxGenerate(c)

	// business logic
	_, err := d.UseCase.Auth(ctx, userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, true)
}
