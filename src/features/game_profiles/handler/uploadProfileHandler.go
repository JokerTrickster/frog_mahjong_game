package handler

import (
	"context"
	_interface "main/features/game_profiles/model/interface"
	"net/http"
	"strconv"

	"main/features/game_profiles/model/entity"

	"github.com/labstack/echo/v4"
)

type UploadProfilesHandler struct {
	UseCase _interface.IUploadProfilesUseCase
}

func NewUploadProfilesHandler(c *echo.Echo, useCase _interface.IUploadProfilesUseCase) _interface.IUploadProfilesHandler {
	handler := &UploadProfilesHandler{
		UseCase: useCase,
	}
	c.POST("/board-game/v0.1/profiles/image", handler.Upload)
	return handler
}

// 프로필 이미지 업로드하기
// @Router /board-game/v0.1/profiles/image [post]
// @Summary 프로필 이미지 업로드하기
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description USER_NOT_FOUND : 유저를 찾을 수 없음
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Param name formData string false "profile image name"
// @Param totalCount formData int false "토탈 카운트"
// @Param image formData file false "프로필 이미지 파일"
// @Param description formData string false "프로필 설명"
// @Produce json
// @Success 200 {object} bool
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags app/board-game/profile
func (d *UploadProfilesHandler) Upload(c echo.Context) error {
	ctx := context.Background()
	file, err := c.FormFile("image")
	if err != nil {
		return err
	}
	totalCount, err := strconv.Atoi(c.FormValue("totalCount"))
	if err != nil {
		return err
	}
	name := c.FormValue("name")
	if err != nil {
		return err
	}
	description := c.FormValue("description")
	entity := entity.ImageUploadProfileEntity{
		Image:       file,
		Name:        name,
		TotalCount:  totalCount,
		Description: description,
	}

	err = d.UseCase.Upload(ctx, entity)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, true)
}
