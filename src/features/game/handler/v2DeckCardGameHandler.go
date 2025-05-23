package handler

import (
	"encoding/json"
	"fmt"
	_interface "main/features/game/model/interface"
	"main/features/game/model/response"
	mw "main/middleware"
	"main/utils"
	_redis "main/utils/db/redis"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type V2DeckCardGameHandler struct {
	UseCase _interface.IV2DeckCardGameUseCase
}

func NewV2DeckCardGameHandler(c *echo.Echo, useCase _interface.IV2DeckCardGameUseCase) _interface.IV2DeckCardGameHandler {
	handler := &V2DeckCardGameHandler{
		UseCase: useCase,
	}
	c.GET("/v2.1/game/:roomID/deck", handler.V2DeckCard, mw.TokenChecker)
	return handler
}

// 카드 정보 가져오기
// @Router /v2.1/game/{roomID}/deck [get]
// @Summary 카드 정보 가져오기
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
// @Param roomID path string true "roomID"
// @Produce json
// @Success 200 {object} response.ResV2DeckCardGame
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags game
func (d *V2DeckCardGameHandler) V2DeckCard(c echo.Context) error {
	ctx, userID, _ := utils.CtxGenerate(c)
	roomID := c.Param("roomID")
	rID, _ := strconv.Atoi(roomID)

	//business logic
	//redis
	cacheKey := fmt.Sprintf("game:bird:%s:deck", roomID)
	cardData, err := _redis.Client.Get(ctx, cacheKey).Result()
	if cardData == "" {
		// 2. 캐시에 데이터가 없을 경우 UseCase에서 조회
		res, err := d.UseCase.V2DeckCard(ctx, int(userID), rID)
		if err != nil {
			return err
		}

		// 3. 조회된 데이터를 Redis에 캐시 (예: 1시간 TTL)
		data, err := json.Marshal(res)
		if err != nil {
			return err
		}
		err = _redis.Client.Set(ctx, cacheKey, data, time.Minute).Err()
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
	var res response.ResV2DeckCardGame
	if err := json.Unmarshal([]byte(cardData), &res); err != nil {
		return err
	}

	// 캐시 히트 여부
	c.Response().Header().Set("X-Cache-Hit", "true")

	return c.JSON(http.StatusOK, res)
}
