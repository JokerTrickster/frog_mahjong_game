package handler

import (
	"context"
	"fmt"

	_interface "main/features/game_users/model/interface"
	"main/features/game_users/model/request"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PushUsersHandler struct {
	UseCase _interface.IPushUsersUseCase
}

func NewPushUsersHandler(c *echo.Echo, useCase _interface.IPushUsersUseCase) _interface.IPushUsersHandler {
	handler := &PushUsersHandler{
		UseCase: useCase,
	}
	c.POST("/v0.1/game/users/push", handler.Push)
	return handler
}

// 전체 유저 메시지 전송하기
// @Router /v0.1/users/push [post]
// @Summary 전체 유저 메시지 전송하기
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description USER_NOT_FOUND : 유저가 존재하지 않음
// @Description ■ errCode with 401
// @Description INVALID_AUTH_CODE : 인증 코드 검증 실패
// @Description TOKEN_BAD : 잘못된 토큰
// @Description INVALID_ACCESS_TOKEN : 잘못된 액세스 토큰
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Description PLAYER_STATE_CHANGE_FAILED : 플레이어 상태 변경 실패
// @param json body request.ReqPushGameUsers true "메시지 내용"
// @Produce json
// @Success 200 {object} bool
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags game/user
func (d *PushUsersHandler) Push(c echo.Context) error {
	ctx := context.Background()
	req := &request.ReqPushGameUsers{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}
	err := d.UseCase.Push(ctx, req)
	if err != nil {
		fmt.Println(err)
		fmt.Println("여기 들어ㅏ오나?")
		return err
	}

	return c.JSON(http.StatusOK, true)
}
