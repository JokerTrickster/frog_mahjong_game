package handler

import (
	"context"
	"encoding/json"
	_interface "main/features/game/model/interface"
	"main/features/game/model/response"
	_redis "main/utils/db/redis"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type ListGameHandler struct {
	UseCase _interface.IListGameUseCase
}

func NewListGameHandler(c *echo.Echo, useCase _interface.IListGameUseCase) _interface.IListGameHandler {
	handler := &ListGameHandler{
		UseCase: useCase,
	}
	c.GET("/v0.1/game/list", handler.ListGame)
	return handler
}

// 게임 정보를 가져온다.
// @Router /v0.1/game/list [get]
// @Summary 게임 정보를 가져온다.
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
// @Success 200 {object} response.ResListGame
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags app/game
func (d *ListGameHandler) ListGame(c echo.Context) error {
	ctx := context.Background()

	//business logic
	//redis
	cacheKey := "game:list"
	cardData, err := _redis.Client.Get(ctx, cacheKey).Result()
	if cardData == "" {
		// 2. 캐시에 데이터가 없을 경우 UseCase에서 조회
		res, err := d.UseCase.ListGame(ctx)
		if err != nil {
			return err
		}

		// 3. 조회된 데이터를 Redis에 캐시 (예: 1시간 TTL)
		data, err := json.Marshal(res)
		if err != nil {
			return err
		}
		err = _redis.Client.Set(ctx, cacheKey, data, 1*time.Hour).Err()
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
	var res response.ResListGame
	if err := json.Unmarshal([]byte(cardData), &res); err != nil {
		return err
	}

	// 캐시 히트 여부
	c.Response().Header().Set("X-Cache-Hit", "true")

	return c.JSON(http.StatusOK, res)
}
