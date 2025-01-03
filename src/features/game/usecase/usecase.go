package usecase

import (
	"context"
	"fmt"
	"main/features/game/model/entity"
	_errors "main/features/game/model/errors"
	"main/features/game/model/request"
	"main/features/game/model/response"
	"main/utils"
	_aws "main/utils/aws"
	"main/utils/db/mysql"
	"math/rand"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

/*
					card 구성
	 1 2 3 4 5 6 7 8 9  (모두 레드) , 중 그린 발 레드
	 1 2 3 4 5 6 7 8 9  (1,5,7,9는 노말, 나머지 그린) , 중 레드 발 그린
	 1 2 3 4 5 6 7 8 9  (1,5,7,9는 노말, 나머지 그린) , 중 레드 발 그린
	 1 2 3 4 5 6 7 8 9  (1,5,7,9는 노말, 나머지 그린) , 중 레드 발 그린
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

// 숫자 문자열과 대응하는 숫자를 맵으로 정의
var numberMap = map[string]int{
	"one":   1,
	"two":   2,
	"three": 3,
	"four":  4,
	"five":  5,
	"six":   6,
	"seven": 7,
	"eight": 8,
	"nine":  9,
}

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
				card.Color = red
			} else if card.Name == "bal" {
				card.Color = green
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

func ScoreCalculate(e *entity.ScoreCalculateEntity, doraCard mysql.Cards) (int, []string, error) {
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
		if IsCheckedSameCard(e.Cards[i], e.Cards[i+1], e.Cards[i+2]) {
			score += 2
			bonuses = append(bonuses, same)
		} else if IsCheckedContinuousCard(e.Cards[i], e.Cards[i+1], e.Cards[i+2]) {
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
	if IsCheckedAllGreen(e.Cards) {
		return 10, []string{allGreen}, nil
	} else if IsCheckedAllRed(e.Cards) {
		return 20, []string{superRed}, nil
	} else if IsCheckedChinYao(e.Cards) {
		return 15, []string{chinYao}, nil
	}
	// 보너스 점수 계산
	// 1. 적패 : color가 red이면 1점 추가
	for _, card := range e.Cards {
		if IsCheckedRedCard(card) {
			bonuses = append(bonuses, red)
			score += 1
		}
	}

	// 2. 탕야오 : 모든 패가 2~8 사이의 패로만 이루어진 카드이면 1점 추가
	if IsCheckedTangYaoCard(e.Cards) {
		bonuses = append(bonuses, tangYao)
		score += 1
	}
	// 3. 찬타 : 두 개의 몸통 모두 1/9/발/중을 포함하고 있을 때 2점 추가
	if IsCheckedChanTa(e.Cards) {
		bonuses = append(bonuses, chanTa)
		score += 2
	}

	// 4. 도라 : dora 하나당 1점 추가
	for _, card := range e.Cards {
		if IsCheckedDora(card, doraCard) {
			bonuses = append(bonuses, dora)
			score += 1
		}
	}

	return score, bonuses, nil
}

// 모두 이름이 같은 카드라면 true 아니라면 false
func IsCheckedSameCard(card1, card2, card3 entity.ScoreCalculateCard) bool {
	if card1.Name == card2.Name && card2.Name == card3.Name {
		return true
	}
	return false
}

// 문자열을 숫자로 변환하는 함수
func convertToNumber(word string) int {
	word = strings.ToLower(word) // 대소문자 구분을 없애기 위해 소문자로 변환
	number, _ := numberMap[word]
	return number
}

// 3개 카드가 연속된 숫자라면 true 아니라면 false
func IsCheckedContinuousCard(card1, card2, card3 entity.ScoreCalculateCard) bool {
	// 카드가 숫자가 아니라면 false 반환
	if !IsCheckedNumberCard([]entity.ScoreCalculateCard{card1, card2, card3}) {
		return false
	}

	// 카드 이름을 int형으로 변경 후 연속된 숫자인지 체크
	card1Int := convertToNumber(card1.Name)
	card2Int := convertToNumber(card2.Name)
	card3Int := convertToNumber(card3.Name)

	if card1Int+1 == card2Int && card2Int+1 == card3Int {
		return true
	}
	return false

}

// 모든 카드가 green이면 true 아니라면 false
func IsCheckedAllGreen(cards []entity.ScoreCalculateCard) bool {
	for _, card := range cards {
		if card.Color != green {
			return false
		}
	}
	return true
}

// 모든 카드가 red이면 true 아니라면 false
func IsCheckedAllRed(cards []entity.ScoreCalculateCard) bool {
	for _, card := range cards {
		if card.Color != red {
			return false
		}
	}
	return true
}

// 모든 패가 1,9,중,발로만 이루어진 카드이면 true 아니라면 false
func IsCheckedChinYao(cards []entity.ScoreCalculateCard) bool {
	for _, card := range cards {
		if card.Name != "one" && card.Name != "nine" && card.Name != "chung" && card.Name != "bal" {
			return false
		}
	}
	return true
}

// 카드가 1~9까지 숫자로만 이루어져 있는지 체크
func IsCheckedNumberCard(cards []entity.ScoreCalculateCard) bool {
	for _, card := range cards {
		if card.Name == "chung" || card.Name == "bal" {
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
func IsCheckedRedCard(cards entity.ScoreCalculateCard) bool {
	if cards.Color == red {
		return true
	}
	return false
}

// 2. 탕야오 : 모든 패가 2~8 사이의 패로만 이루어진 카드라면 true 아니라면 false
func IsCheckedTangYaoCard(cards []entity.ScoreCalculateCard) bool {
	for _, card := range cards {
		if card.Name == "one" || card.Name == "nine" || card.Name == "chung" || card.Name == "bal" {
			return false
		}
	}
	return true
}

// 3. 찬타 : index 0,1,2 / 3,4,5 모두 1/9/발 중을 포함하고 있을 때 true 아니라면 false
func IsCheckedChanTa(cards []entity.ScoreCalculateCard) bool {
	result1 := false
	result2 := false
	if (cards[0].Name == "one" || cards[0].Name == "nine" || cards[0].Name == "chung" || cards[0].Name == "bal") ||
		(cards[1].Name == "one" || cards[1].Name == "nine" || cards[1].Name == "chung" || cards[1].Name == "bal") ||
		(cards[2].Name == "one" || cards[2].Name == "nine" || cards[2].Name == "chung" || cards[2].Name == "bal") {
		result1 = true
	}

	if (cards[3].Name == "one" || cards[3].Name == "nine" || cards[3].Name == "chung" || cards[3].Name == "bal") ||
		(cards[4].Name == "one" || cards[4].Name == "nine" || cards[4].Name == "chung" || cards[4].Name == "bal") ||
		(cards[5].Name == "one" || cards[5].Name == "nine" || cards[5].Name == "chung" || cards[5].Name == "bal") {
		result2 = true
	}
	if result1 && result2 {
		return true
	}
	return false
}

// 4. 도라 : dora 하나당 1점 추가
func IsCheckedDora(card entity.ScoreCalculateCard, doraCard mysql.Cards) bool {
	if card.Name == doraCard.Name {
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

func CreateScoreCalculateEntitySQL(userID uint, req *request.ReqScoreCalculate) *entity.ScoreCalculateEntitySQL {
	entitySQL := &entity.ScoreCalculateEntitySQL{
		UserID: userID,
		RoomID: req.RoomID,
	}
	cards := make([]uint, 0)
	for _, card := range req.Cards {
		cards = append(cards, card.CardID)
	}
	entitySQL.Cards = cards
	return entitySQL

}
func CreateResultEntitySQL(userID uint, req *request.ReqResult) *entity.ResultEntitySQL {
	entitySQL := &entity.ResultEntitySQL{
		RoomID: req.RoomID,
	}
	cards := make([]uint, 0)
	for _, card := range req.Cards {
		cards = append(cards, card.CardID)
	}
	entitySQL.Cards = cards
	return entitySQL
}

func CardValidation(cardsDTO []mysql.Cards, cards []request.ScoreCard) error {
	// 카드가 6장인지 체크
	if len(cardsDTO) != 6 {
		return utils.ErrorMsg(context.TODO(), utils.ErrBadParameter, utils.Trace(), _errors.ErrNotEnoughCard.Error(), utils.ErrFromClient)
	}
	// 보유하고 있는 카드와 전달받은 카드가 일치하는지 체크
	sameCard := 0
	for _, card := range cards {
		for _, cardDTO := range cardsDTO {
			if card.CardID == uint(cardDTO.CardID) {
				sameCard++
				break
			}
		}
	}
	if sameCard != 6 {
		return utils.ErrorMsg(context.TODO(), utils.ErrBadParameter, utils.Trace(), _errors.ErrBadRequestCard.Error(), utils.ErrFromClient)
	}

	return nil
}

func CreateScoreCalculateEntity(cardsDTO []mysql.Cards, cards []request.ScoreCard) *entity.ScoreCalculateEntity {
	result := &entity.ScoreCalculateEntity{
		Cards: make([]entity.ScoreCalculateCard, 0),
	}
	for _, card := range cards {
		for _, cardDTO := range cardsDTO {
			if card.CardID == uint(cardDTO.CardID) {
				result.Cards = append(result.Cards, entity.ScoreCalculateCard{
					CardID: uint(cardDTO.CardID),
					Name:   cardDTO.Name,
					Color:  cardDTO.Color,
				})
				break
			}
		}
	}

	return result
}

func CreateResultEntity(cardsDTO []mysql.Cards, cards []request.ResultCard) *entity.ScoreCalculateEntity {
	result := &entity.ScoreCalculateEntity{
		Cards: make([]entity.ScoreCalculateCard, 0),
	}
	for _, card := range cards {
		for _, cardDTO := range cardsDTO {
			if card.CardID == uint(cardDTO.CardID) {
				result.Cards = append(result.Cards, entity.ScoreCalculateCard{
					CardID: uint(cardDTO.CardID),
					Name:   cardDTO.Name,
					Color:  cardDTO.Color,
				})
				break
			}
		}
	}
	return result
}

func CreateV2ReportDTO(userID uint, req *request.ReqV2Report) *mysql.Reports {
	return &mysql.Reports{
		TargetUserID:   int(req.TargetUserID),
		ReporterUserID: int(userID),
		CategoryID:     int(req.CategoryID),
		Reason:         req.Reason,
	}
}
func CreateReportDTO(userID uint, req *request.ReqReport) *mysql.Reports {
	return &mysql.Reports{
		TargetUserID:   int(req.TargetUserID),
		ReporterUserID: int(userID),
		CategoryID:     int(req.CategoryID),
		Reason:         req.Reason,
	}
}
func CreateResMetaGame(categoryList []mysql.Categories) response.ResMetaGame {
	res := response.ResMetaGame{}
	for _, category := range categoryList {
		category := response.Category{
			ID:     uint(category.ID),
			Reason: category.Reason,
		}
		res.Categories = append(res.Categories, category)
	}
	return res
}

func CreateRandomCardIDList() response.ResDeckCardGame {
	res := response.ResDeckCardGame{}
	cardIDList := make([]int, 0)
	for i := 1; i <= 44; i++ {
		cardIDList = append(cardIDList, i)
	}
	rand.Shuffle(len(cardIDList), func(i, j int) {
		cardIDList[i], cardIDList[j] = cardIDList[j], cardIDList[i]
	})
	res.CardIDList = cardIDList
	return res

}

func CreateResListMission(missionList []*mysql.Missions) response.ResListMissionGame {
	res := response.ResListMissionGame{}
	for _, mission := range missionList {
		m := response.Mission{
			ID:          int(mission.ID),
			Title:       mission.Name,
			Description: mission.Description,
		}
		if mission.Image != "" {
			imageUrl, err := _aws.ImageGetSignedURL(context.TODO(), mission.Image, _aws.ImgTypeMission)
			if err != nil {
				return response.ResListMissionGame{}
			}
			m.Image = imageUrl
		}
		res.Missions = append(res.Missions, m)
	}
	return res
}
func CreateMissionDTO(req *request.ReqCreateMission, fileName string) *mysql.Missions {
	return &mysql.Missions{
		Name:        req.Name,
		Description: req.Description,
		Image:       fileName,
	}
}

func CreateResV2ListCard(cards []*mysql.BirdCards, count int) response.ResV2ListCardGame {
	res := response.ResV2ListCardGame{
		TotalCount: count,
	}
	for _, card := range cards {
		c := response.BirdCard{
			ID:            int(card.ID),
			Name:          card.Name,
			Size:          card.Size,
			Habitat:       card.Habitat,
			BeakDirection: card.BeakDirection,
			Nest:          card.Nest,
		}
		// s3 에서 서명된 url로 응답
		imageUrl, err := _aws.ImageGetSignedURL(context.TODO(), card.Image, _aws.ImgTypeBirdCard)
		if err != nil {
			fmt.Println(err)
			return response.ResV2ListCardGame{}
		}
		c.Image = imageUrl
		res.Cards = append(res.Cards, c)
	}
	return res
}

func CreateV2RandomCardIDList() response.ResV2DeckCardGame {
	res := response.ResV2DeckCardGame{}
	cardIDList := make([]int, 0)
	for i := 1; i <= 9; i++ {
		cardIDList = append(cardIDList, i)
	}
	rand.Shuffle(len(cardIDList), func(i, j int) {
		cardIDList[i], cardIDList[j] = cardIDList[j], cardIDList[i]
	})
	res.CardIDList = cardIDList
	return res

}

func CreateBirdCardsDTO(req *request.ReqSaveCardInfo) []mysql.BirdCards {
	birdCardsDTO := make([]mysql.BirdCards, 0)
	for _, card := range req.Cards {
		birdCardsDTO = append(birdCardsDTO, mysql.BirdCards{
			Name:          card.Name,
			Size:          card.Size,
			Habitat:       card.Habitat,
			BeakDirection: card.BeakDirection,
			Nest:          card.Nest,
			Image:         card.Image,
		})
	}
	return birdCardsDTO
}

func UpdateBirdCardsDTO(req *request.ReqUpdateCard) mysql.BirdCards {
	result := mysql.BirdCards{
		Model: gorm.Model{
			ID: uint(req.Card.CardID),
		},
	}
	if req.Card.Name != "" {
		result.Name = req.Card.Name
	}

	if req.Card.Size != 0 {
		result.Size = req.Card.Size
	}

	if req.Card.Habitat != "" {
		result.Habitat = req.Card.Habitat
	}

	if req.Card.BeakDirection != "" {
		result.BeakDirection = req.Card.BeakDirection
	}

	if req.Card.Nest != "" {
		result.Nest = req.Card.Nest
	}
	return result
}

func CreateResListCard(cards []*mysql.FrogCards, count int) response.ResListCardGame {
	res := response.ResListCardGame{
		TotalCount: count,
	}
	for _, card := range cards {
		c := response.FrogCard{
			ID:    int(card.ID),
			Name:  card.Name,
			Color: card.Color,
		}
		res.Cards = append(res.Cards, c)
	}
	return res

}
