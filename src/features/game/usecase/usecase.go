package usecase

import "main/utils/db/mysql"

/*
					card 구성
	 1 2 3 4 5 6 7 8 9  (모두 레드) , 중 그린 발 레드
	 1 2 3 4 5 6 7 8 9  (1,5,7,9는 노말, 나머지 그린) , 중 그린 발 레드
	 1 2 3 4 5 6 7 8 9  (1,5,7,9는 노말, 나머지 그린) , 중 그린 발 레드
	 1 2 3 4 5 6 7 8 9  (1,5,7,9는 노말, 나머지 그린) , 중 그린 발 레드
*/
const (
	red    = "red"
	green  = "green"
	normal = "normal"
)

var cardNames = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "중", "발"}

func CreateInitCards(roomID uint) []mysql.Cards {
	cards := make([]mysql.Cards, 0)

	for i := 0; i < 4; i++ {
		for j := 0; j < 11; j++ {
			card := mysql.Cards{
				RoomID: int(roomID),
				Name:   cardNames[j],
				State:  "none",
			}
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
			if card.Name == "중" {
				card.Color = green
			} else if card.Name == "발" {
				card.Color = red
			}

			cards = append(cards, card)
		}
	}

	return cards
}

func CheckRoomUsersReady(roomUsers []mysql.RoomUsers) bool {
	for _, ru := range roomUsers {
		if ru.PlayerState != "ready" {
			return false
		}
	}
	return true
}
