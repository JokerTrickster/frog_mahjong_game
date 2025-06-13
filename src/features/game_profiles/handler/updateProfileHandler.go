package handler

import (
	"context"
	_interface "main/features/game_profiles/model/interface"
	"strconv"

	"net/http"

	"github.com/labstack/echo/v4"
)

type UpdateProfilesHandler struct {
	UseCase _interface.IUpdateProfilesUseCase
}

func NewUpdateProfilesHandler(c *echo.Echo, useCase _interface.IUpdateProfilesUseCase) _interface.IUpdateProfilesHandler {
	handler := &UpdateProfilesHandler{
		UseCase: useCase,
	}
	c.PATCH("/board-game/v0.1/users/:userID/profiles/:profileID", handler.Update)
	return handler
}

// 프로필 이미지 변경하기
// @Router /board-game/v0.1/users/:userID/profiles/:profileID [patch]
// @Summary 프로필 이미지 변경하기
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description USER_NOT_FOUND : 유저를 찾을 수 없음
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Param userID path string true "유저 ID"
// @Param profileID path string true "프로필 ID"
// @Produce json
// @Success 200 {object} bool
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags app/board-game/profile
func (d *UpdateProfilesHandler) Update(c echo.Context) error {
	ctx := context.Background()
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		return err
	}
	profileID, err := strconv.Atoi(c.Param("profileID"))
	if err != nil {
		return err
	}

	res, err := d.UseCase.Update(ctx, userID, profileID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
