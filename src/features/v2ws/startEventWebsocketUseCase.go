package v2ws

import (
	"main/utils/db/mysql"

	"golang.org/x/exp/rand"
)

func CheckRoomUsersReady(roomUsers []mysql.RoomUsers, ownerID uint) bool {
	for _, ru := range roomUsers {
		if ru.UserID == int(ownerID) {
			continue
		}
		if ru.PlayerState != "ready" {
			return false
		}
	}
	return true
}

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

func CreateInitCards(roomID uint, count int) []mysql.UserBirdCards {
	cards := make([]mysql.UserBirdCards, 0)
	// 총 카드 수 만큼 생성하면 된다.
	for i := 0; i < count; i++ {
		cards = append(cards, mysql.UserBirdCards{
			RoomID: int(roomID),
			CardID: i + 1,
			State:  "none",
			UserID: 0,
		})
	}

	return cards
}
