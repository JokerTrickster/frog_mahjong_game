package usecase

import (
	"context"
	_errors "main/features/game/model/errors"
	"main/features/game/model/request"
	"main/utils"
	"main/utils/db/mysql"
	"math/rand"
	"strconv"
)

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
			if card.Name == "chung" {
				card.Color = green
			} else if card.Name == "bal" {
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

func CreateUpdateRoomUser(roomUser mysql.RoomUsers, req *request.ReqDiscard) mysql.RoomUsers {
	roomUser.OwnedCardCount -= 1
	roomUser.PlayerState = "play_wait"
	return roomUser
}

func ScoreCalculate(req *request.ReqScoreCalculate, doraCard mysql.Cards) (int, []string, error) {
	// 카드의 점수를 계산한다.
	score := 0
	// 카드는 1,2,3,4,5,6,7,8,9,중,발 이 있다.
	// 인덱스 0,1,2,3,4,5 에서 0,1,2 와 3,4,5 이 아래 조건 1번 2번 중 하나를 만족하면 점수 계산한다.
	// 1. 3개 카드가 연속된 숫자이면 1점
	// 2. 3개 카드가 같은 숫자이면 2점
	// 3. 중 이나 발 카드는 같은 글자로만 이루어져야 한다. (중 중 중, 발 발 발) 2점
	bonuses := make([]string, 0)

	// 3개 카드가 모두 같은 카드인지 체크
	for i := 0; i < 6; i += 3 {
		if IsCheckedSameCard(req.Cards[i], req.Cards[i+1], req.Cards[i+2]) {
			score += 2
			bonuses = append(bonuses, same)
		} else if IsCheckedContinuousCard(req.Cards[i], req.Cards[i+1], req.Cards[i+2]) {
			score += 1
			bonuses = append(bonuses, continuous)
		} else {
			return 0, []string{}, utils.ErrorMsg(context.TODO(), utils.ErrNotEnoughCond, utils.Trace(), _errors.ErrNotEnoughCondition.Error(), utils.ErrFromClient)
		}
	}

	// 역만 보너스 점수 계산하기
	// 1. 올 그린 : 카드가(2,3,4,6,8,발)로만 이루어진 카드이고 color가 green이면 10점
	// 2. 슈퍼 레드 : 모든 카드가(1,2,3,4,5,6,7,8,9,중)로만 이루어진 카드이고 color가 red이면 20점
	// 3. 칭야오 : 모든 패가 1,9,중,발로만 이루어진 카드이면 15점

	// 역만 보너스 점수 계산
	if IsCheckedAllGreen(req.Cards) {
		return 10, []string{allGreen}, nil
	} else if IsCheckedAllRed(req.Cards) {
		return 20, []string{superRed}, nil
	} else if IsCheckedChinYao(req.Cards) {
		return 15, []string{chinYao}, nil
	}
	// 보너스 점수 계산
	// 1. 적패 : color가 red이면 1점 추가
	for _, card := range req.Cards {
		if IsCheckedRedCard(card) {
			bonuses = append(bonuses, red)
			score += 1
		}
	}

	// 2. 탕야오 : 모든 패가 2~8 사이의 패로만 이루어진 카드이면 1점 추가
	if IsCheckedTangYaoCard(req.Cards) {
		bonuses = append(bonuses, tangYao)
		score += 1
	}
	// 3. 찬타 : 두 개의 몸통 모두 1/9/발/중을 포함하고 있을 때 2점 추가
	if IsCheckedChanTa(req.Cards) {
		bonuses = append(bonuses, chanTa)
		score += 2
	}

	// 4. 도라 : dora 하나당 1점 추가
	for _, card := range req.Cards {
		if IsCheckedDora(card, doraCard) {
			bonuses = append(bonuses, dora)
			score += 1
		}
	}

	return score, bonuses, nil
}

// 모두 이름이 같은 카드라면 true 아니라면 false
func IsCheckedSameCard(card1, card2, card3 request.ScoreCard) bool {
	if card1.Name == card2.Name && card2.Name == card3.Name {
		return true
	}
	return false
}

// 3개 카드가 연속된 숫자라면 true 아니라면 false
func IsCheckedContinuousCard(card1, card2, card3 request.ScoreCard) bool {
	// 카드가 숫자가 아니라면 false 반환
	if !IsCheckedNumberCard([]request.ScoreCard{card1, card2, card3}) {
		return false
	}
	// 카드 이름을 int형으로 변경 후 연속된 숫자인지 체크
	card1Int := ConvertStringToInt(card1.Name)
	card2Int := ConvertStringToInt(card2.Name)
	card3Int := ConvertStringToInt(card3.Name)

	if card1Int+1 == card2Int && card2Int+1 == card3Int {
		return true
	}
	return false

}

// 모든 카드가 green이면 true 아니라면 false
func IsCheckedAllGreen(cards []request.ScoreCard) bool {
	for _, card := range cards {
		if card.Color != green {
			return false
		}
	}
	return true
}

// 모든 카드가 red이면 true 아니라면 false
func IsCheckedAllRed(cards []request.ScoreCard) bool {
	for _, card := range cards {
		if card.Color != red {
			return false
		}
	}
	return true
}

// 모든 패가 1,9,중,발로만 이루어진 카드이면 true 아니라면 false
func IsCheckedChinYao(cards []request.ScoreCard) bool {
	for _, card := range cards {
		if card.Name != "1" && card.Name != "9" && card.Name != "중" && card.Name != "발" {
			return false
		}
	}
	return true
}

// 카드가 1~9까지 숫자로만 이루어져 있는지 체크
func IsCheckedNumberCard(cards []request.ScoreCard) bool {
	for _, card := range cards {
		if card.Name == "중" || card.Name == "발" {
			return false
		}
	}
	return true
}

// string을 int로 변환
func ConvertStringToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

// 1. 적패 : color가 red이면 1점 추가
func IsCheckedRedCard(cards request.ScoreCard) bool {
	if cards.Color == red {
		return true
	}
	return false
}

// 2. 탕야오 : 모든 패가 2~8 사이의 패로만 이루어진 카드라면 true 아니라면 false
func IsCheckedTangYaoCard(cards []request.ScoreCard) bool {
	for _, card := range cards {
		if card.Name == "1" || card.Name == "9" || card.Name == "중" || card.Name == "발" {
			return false
		}
	}
	return true
}

// 3. 찬타 : index 0,1,2 / 3,4,5 모두 1/9/발 중을 포함하고 있을 때 true 아니라면 false
func IsCheckedChanTa(cards []request.ScoreCard) bool {
	result1 := false
	result2 := false
	if (cards[0].Name == "1" || cards[0].Name == "9" || cards[0].Name == "중" || cards[0].Name == "발") ||
		(cards[1].Name == "1" || cards[1].Name == "9" || cards[1].Name == "중" || cards[1].Name == "발") ||
		(cards[2].Name == "1" || cards[2].Name == "9" || cards[2].Name == "중" || cards[2].Name == "발") {
		result1 = true
	}

	if (cards[3].Name == "1" || cards[3].Name == "9" || cards[3].Name == "중" || cards[3].Name == "발") ||
		(cards[4].Name == "1" || cards[4].Name == "9" || cards[4].Name == "중" || cards[4].Name == "발") ||
		(cards[5].Name == "1" || cards[5].Name == "9" || cards[5].Name == "중" || cards[5].Name == "발") {
		result2 = true
	}
	if result1 && result2 {
		return true
	}
	return false
}

// 4. 도라 : dora 하나당 1점 추가
func IsCheckedDora(card request.ScoreCard, doraCard mysql.Cards) bool {
	if card.Name == doraCard.Name && card.Color == doraCard.Color {
		return true
	}
	return false
}

func IsCheckedWinRequest(roomUser mysql.RoomUsers, score int) bool {
	if (roomUser.PlayerState == "play" || roomUser.PlayerState == "loan") && roomUser.OwnedCardCount == 6 {
		if score >= 5 {
			return true
		}
	}
	return false
}
