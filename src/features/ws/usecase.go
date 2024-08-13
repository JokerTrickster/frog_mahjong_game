package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/features/ws/model/entity"
	"main/features/ws/repository"
	"main/utils/db/mysql"
)

func CalcPlayTurn(playTurn, playerCount int) int {
	return (playTurn % playerCount) + 1
}

func CreateRoomInfoMSG(ctx context.Context, roomID uint, playTurn int) *entity.RoomInfo {
	roomInfoMsg := entity.RoomInfo{}

	// 현재 참여하고 있는 유저에 대한 정보를 가져와서 메시지 전달한다.
	preloadUsers, err := repository.ImportSingleCardFindAllRoomUsers(ctx, roomID)
	if err != nil {
		log.Println(err)
	}
	//유저 정보 저장
	for _, roomUser := range preloadUsers {
		user := entity.User{
			ID:          uint(roomUser.UserID),
			PlayerState: roomUser.PlayerState,
			Coin:        roomUser.User.Coin,
			Name:        roomUser.User.Name,
			Email:       roomUser.User.Email,
			TurnNumber:  roomUser.TurnNumber,
		}
		for _, card := range roomUser.Cards {
			if card.State == "owned" {
				user.Cards = append(user.Cards, &entity.Card{
					CardID: uint(card.CardID),
					UserID: uint(card.UserID),
				})
			} else if card.State == "discard" {
				user.DiscardedCards = append(user.DiscardedCards, &entity.Card{
					CardID: uint(card.CardID),
					UserID: uint(card.UserID),
				})
			}
		}

		if roomUser.Room.OwnerID == roomUser.UserID {
			user.IsOwner = true
		}
		roomInfoMsg.Users = append(roomInfoMsg.Users, &user)
	}

	//게임 정보 저장
	gameInfo := entity.GameInfo{
		PlayTurn: playTurn,
		AllReady: true,
	}
	roomInfoMsg.GameInfo = &gameInfo
	return &roomInfoMsg

}

func CreateMessage(roomInfoMsg *entity.RoomInfo) (string, error) {
	// 구조체를 JSON 문자열로 변환 (마샬링)
	jsonData, err := json.Marshal(roomInfoMsg)
	if err != nil {
		return "", fmt.Errorf("JSON 마샬링 에러: %s", err)
	}

	return string(jsonData), nil
}

func CalcScore(cards []*mysql.Cards, score int) error {
	if score >= 5 {
		return nil
	}
	return fmt.Errorf("점수가 부족합니다.")
}
