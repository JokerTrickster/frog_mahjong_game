package v2ws

import (
	"main/utils/db/mysql"
	"time"

	"golang.org/x/exp/rand"
)

func StartUpdateRoomUsers(roomUsers []mysql.RoomUsers) ([]mysql.RoomUsers, error) {
	visited := make(map[int]bool, len(roomUsers)+1)

	for i := range roomUsers {
		roomUsers[i].PlayerState = "play"
		for {
			// 플레이 순번을 인원수에 맞게 랜덤으로 생성하되 중복되지 않게 생성
			random := rand.Intn(len(roomUsers)) + 1
			if !visited[random] {
				roomUsers[i].TurnNumber = random
				visited[random] = true
				break
			}
		}
	}
	return roomUsers, nil
}

func CreateInitCards(roomID uint, birdCards []*mysql.BirdCards) []mysql.UserBirdCards {
	cards := make([]mysql.UserBirdCards, 0)
	// 총 카드 수 만큼 생성하면 된다.
	for i := 0; i < len(birdCards); i++ {
		cards = append(cards, mysql.UserBirdCards{
			RoomID: int(roomID),
			CardID: int(birdCards[i].ID),
			State:  "none",
			UserID: 0,
		})
	}

	return cards
}

func StartUpdateRoom(roomID uint) mysql.Rooms {
	//시작 시간 (epoch time milliseconds)
	return mysql.Rooms{
		State:     "play",
		StartTime: time.Now(),
	}
}
