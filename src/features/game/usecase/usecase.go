package usecase

import "main/utils/db/mysql"

func CreateInitCards(roomID uint) []mysql.Cards {
	cards := []mysql.Cards{}
	for i := 0; i < 4; i++ {
		for j := 0; j < 10; j++ {
			var card mysql.Cards
			card.RoomID = int(roomID)
			card.Name = j
			card.Color = i
			card.State = "none"
			cards = append(cards, card)
		}
	}
	return cards
}
