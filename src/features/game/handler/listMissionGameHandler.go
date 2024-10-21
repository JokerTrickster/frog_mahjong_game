package handler

import (
	"context"
	"encoding/json"
	"fmt"
	_interface "main/features/game/model/interface"
	"main/features/game/model/response"
	_redis "main/utils/db/redis"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type ListMissionGameHandler struct {
	UseCase _interface.IListMissionGameUseCase
}

func NewListMissionGameHandler(c *echo.Echo, useCase _interface.IListMissionGameUseCase) _interface.IListMissionGameHandler {
	handler := &ListMissionGameHandler{
		UseCase: useCase,
	}
	c.GET("/v0.1/game/missions", handler.ListMission)
	return handler
}

// 미션 리스트 가져오기
// @Router /v0.1/game/missions [get]
// @Summary 미션 리스트 가져오기
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
// @Success 200 {object} response.ResListMissionGame
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags game
func (d *ListMissionGameHandler) ListMission(c echo.Context) error {

	ctx := context.Background()
	//business logic
	//redis
	cacheKey := fmt.Sprintf("game:missions")
	missionData, err := _redis.Client.Get(ctx, cacheKey).Result()
	if missionData == "" {
		// 2. 캐시에 데이터가 없을 경우 UseCase에서 조회
		res, err := d.UseCase.ListMission(ctx)
		if err != nil {
			return err
		}

		// 3. 조회된 데이터를 Redis에 캐시 (예: 1시간 TTL)
		data, err := json.Marshal(res)
		if err != nil {
			return err
		}
		err = _redis.Client.Set(ctx, cacheKey, data, time.Hour).Err()
		if err != nil {
			return err
		}

		// 캐시 히트 여부
		c.Response().Header().Set("X-Cache-Hit", "false")

		// 4. DB에서 조회한 데이터 반환
		return c.JSON(http.StatusOK, res)
	} else if err != nil {
		// Redis 오류 처리
		return err
	}

	// 5. 캐시된 데이터가 있을 경우 반환
	var res response.ResListMissionGame
	if err := json.Unmarshal([]byte(missionData), &res); err != nil {
		return err
	}

	// 캐시 히트 여부
	c.Response().Header().Set("X-Cache-Hit", "true")

	return c.JSON(http.StatusOK, res)
}
