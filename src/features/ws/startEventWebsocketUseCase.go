package ws

import (
	"main/features/ws/model/entity"
	"main/utils/db/mysql"

	"golang.org/x/exp/rand"
)

func CheckRoomUsersReady(roomUsers []mysql.FrogRoomUsers, ownerID uint) bool {
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

func StartUpdateRoomUsers(roomUsers []mysql.FrogRoomUsers) ([]mysql.FrogRoomUsers, *entity.ErrorInfo) {
	visited := make(map[int]bool, len(roomUsers)+1)

	for i := range roomUsers {
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

/*
					card 구성
	 1 2 3 4 5 6 7 8 9  (모두 레드) , 중 그린 발 레드
	 1 2 3 4 5 6 7 8 9  (1,5,7,9는 노말, 나머지 그린) , 중 그린 발 레드
	 1 2 3 4 5 6 7 8 9  (1,5,7,9는 노말, 나머지 그린) , 중 그린 발 레드
	 1 2 3 4 5 6 7 8 9  (1,5,7,9는 노말, 나머지 그린) , 중 그린 발 레드
*/
const (
	red        = "red"
	green      = "green"
	normal     = "normal"
	allGreen   = "allGreen"
	superRed   = "superRed"
	tangYao    = "tangYao"
	chanTa     = "chanTa"
	chinYao    = "chinYao"
	dora       = "dora"
	same       = "same"
	continuous = "continuous"
)

var cardNames = []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "chung", "bal"}

func CreateInitCards(roomID uint, cards []mysql.FrogCards) []mysql.FrogUserCards {
	var userCards []mysql.FrogUserCards
	for _, card := range cards {
		userCard := mysql.FrogUserCards{
			UserID: 0,
			RoomID: int(roomID),
			CardID: int(card.Model.ID),
			State:  "none",
		}
		userCards = append(userCards, userCard)
	}

	return userCards
}
