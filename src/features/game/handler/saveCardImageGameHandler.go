package handler

import (
	"context"
	"main/features/game/model/entity"
	_interface "main/features/game/model/interface"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SaveCardImageGameHandler struct {
	UseCase _interface.ISaveCardImageGameUseCase
}

func NewSaveCardImageGameHandler(c *echo.Echo, useCase _interface.ISaveCardImageGameUseCase) _interface.ISaveCardImageGameHandler {
	handler := &SaveCardImageGameHandler{
		UseCase: useCase,
	}
	c.POST("/v0.1/game/cards/image", handler.SaveCardImage)
	return handler
}

// 카드 이미지를 저장한다.
// @Router /v0.1/game/cards/image [post]
// @Summary 카드 이미지를 저장한다.
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description NOT_OWNER : 방장이 시작 요청을 하지 않음
// @Description NOT_FIRST_PLAYER : 첫 플레이어가 아님
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Produce json
// @Param image formData file false "새카드 이미지 파일"
// @Success 200 {object} bool
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags game
func (d *SaveCardImageGameHandler) SaveCardImage(c echo.Context) error {
	ctx := context.Background()
	file, err := c.FormFile("image")
	if err != nil {
		return err
	}
	entity := entity.SaveCardImageGameEntity{
		Image: file,
	}

	//business logic
	err = d.UseCase.SaveCardImage(ctx, entity)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, true)
}
