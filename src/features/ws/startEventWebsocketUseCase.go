package ws

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

func CreateInitCards(roomID uint) []mysql.Cards {
	cards := make([]mysql.Cards, 0)
	cardID := 1
	for i := 0; i < 4; i++ {
		for j := 0; j < 11; j++ {
			card := mysql.Cards{
				RoomID: int(roomID),
				Name:   cardNames[j],
				State:  "none",
				CardID: cardID,
			}
			cardID++
			if i == 0 {
				card.Color = red
			} else {
				// 1,5,7,9는 노말, 나머지 그린
				if j == 0 || j == 4 || j == 6 || j == 8 {
					card.Color = normal
				} else {
					card.Color = green
				}
			}
			if card.Name == "chung" {
				card.Color = red
			} else if card.Name == "bal" {
				card.Color = green
			}

			cards = append(cards, card)
		}
	}

	return cards
}
