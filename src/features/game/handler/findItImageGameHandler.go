package handler

import (
	"context"
	"main/features/game/model/entity"
	_interface "main/features/game/model/interface"
	"net/http"

	"github.com/labstack/echo/v4"
)

type FindItImageGameHandler struct {
	UseCase _interface.IFindItImageGameUseCase
}

func NewFindItImageGameHandler(c *echo.Echo, useCase _interface.IFindItImageGameUseCase) _interface.IFindItImageGameHandler {
	handler := &FindItImageGameHandler{
		UseCase: useCase,
	}
	c.POST("/find-it/v0.1/game/image", handler.FindItImage)
	return handler
}

// 틀린그림찾기 이미지를 저장한다.
// @Router /find-it/v0.1/game/image [post]
// @Summary 틀린그림찾기 이미지를 저장한다.
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Produce json
// @Param image formData file false "틀린그림찾기 이미지 파일"
// @Success 200 {object} bool
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags game
func (d *FindItImageGameHandler) FindItImage(c echo.Context) error {
	ctx := context.Background()
	file, err := c.FormFile("image")
	if err != nil {
		return err
	}
	entity := &entity.FindItImageGameEntity{
		Image: file,
	}

	//business logic
	err = d.UseCase.FindItImage(ctx, entity)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, true)
}
