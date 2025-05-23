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
				Count:  j + 1,
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
			if card.Count == 10 {
				card.Color = red
			} else if card.Count == 11 {
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

func ScoreCalculate(e *entity.ScoreCalculateEntity, doraCard *entity.ResultCardEntity) (int, []string, error) {
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
	if card1.Count == card2.Count && card2.Count == card3.Count {
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
	card1Int := card1.Count
	card2Int := card2.Count
	card3Int := card3.Count

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
		if card.Count != 1 && card.Count != 9 && card.Count != 10 && card.Count != 11 {
			return false
		}
	}
	return true
}

// 카드가 1~9까지 숫자로만 이루어져 있는지 체크
func IsCheckedNumberCard(cards []entity.ScoreCalculateCard) bool {
	for _, card := range cards {
		if card.Count == 10 || card.Count == 11 {
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
		if card.Count == 1 || card.Count == 9 || card.Count == 10 || card.Count == 11 {
			return false
		}
	}
	return true
}

// 3. 찬타 : index 0,1,2 / 3,4,5 모두 1/9/발 중을 포함하고 있을 때 true 아니라면 false
func IsCheckedChanTa(cards []entity.ScoreCalculateCard) bool {
	result1 := false
	result2 := false
	if (cards[0].Count == 1 || cards[0].Count == 9 || cards[0].Count == 10 || cards[0].Count == 11) ||
		(cards[1].Count == 1 || cards[1].Count == 9 || cards[1].Count == 10 || cards[1].Count == 11) ||
		(cards[2].Count == 1 || cards[2].Count == 9 || cards[2].Count == 10 || cards[2].Count == 11) {
		result1 = true
	}

	if (cards[3].Count == 1 || cards[3].Count == 9 || cards[3].Count == 10 || cards[3].Count == 11) ||
		(cards[4].Count == 1 || cards[4].Count == 9 || cards[4].Count == 10 || cards[4].Count == 11) ||
		(cards[5].Count == 1 || cards[5].Count == 9 || cards[5].Count == 10 || cards[5].Count == 11) {
		result2 = true
	}
	if result1 && result2 {
		return true
	}
	return false
}

// 4. 도라 : dora 하나당 1점 추가
func IsCheckedDora(card entity.ScoreCalculateCard, doraCard *entity.ResultCardEntity) bool {
	if card.Count == doraCard.Count {
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
					Count:  cardDTO.Count,
					Color:  cardDTO.Color,
				})
				break
			}
		}
	}

	return result
}

func CreateResultEntity(cardsDTO []*entity.ResultCardEntity, cards []request.ResultCard) *entity.ScoreCalculateEntity {
	result := &entity.ScoreCalculateEntity{
		Cards: make([]entity.ScoreCalculateCard, 0),
	}
	for _, card := range cards {
		for _, cardDTO := range cardsDTO {
			if card.CardID == uint(cardDTO.CardID) {
				result.Cards = append(result.Cards, entity.ScoreCalculateCard{
					CardID: uint(cardDTO.CardID),
					Count:  cardDTO.Count,
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
			Count: card.Count,
			Color: card.Color,
		}
		res.Cards = append(res.Cards, c)
	}
	return res

}

func CreateResultCardEntity(cardDTO *mysql.FrogUserCards, frogCard *mysql.FrogCards) *entity.ResultCardEntity {
	return &entity.ResultCardEntity{
		CardID: uint(cardDTO.CardID),
		Count:  frogCard.Count,
		Color:  frogCard.Color,
		State:  cardDTO.State,
	}
}

func CreateResV2DrawResult(userMissions []mysql.UserMissions) response.ResV2DrawResult {

	// Temporary map to group missions by UserID
	userMissionMap := make(map[int][]int)

	// Group missions by UserID
	for _, userMission := range userMissions {
		userMissionMap[userMission.UserID] = append(userMissionMap[userMission.UserID], userMission.MissionID)
	}

	// Create the response
	var res response.ResV2DrawResult
	for userID, missions := range userMissionMap {
		res.Users = append(res.Users, response.DrawResult{
			UserID:          userID,
			SuccessMissions: missions,
		})
	}

	return res
}

// find-it
func CreateResResult(roomSettingDTO *mysql.FindItRoomSettings, userCorrectPositionsDTO []*mysql.FindItUserCorrectPositions, userDTOs []*mysql.GameUsers) response.ResFindItResult {
	userCorrectCount := make(map[int]int)

	// 유저별 맞춘 개수 카운트
	for _, userCorrectPosition := range userCorrectPositionsDTO {
		userCorrectCount[userCorrectPosition.UserID]++
	}
	res := response.ResFindItResult{
		Round: roomSettingDTO.Round,
	}
	// 유저 ID를 이름으로 변환하기 위한 map 생성
	userNames := make(map[int]string)
	for _, userDTO := range userDTOs {
		userNames[int(userDTO.ID)] = userDTO.Name
	}

	// 응답 구조에 맞추어 데이터 저장 (모든 유저 포함)
	var users []response.UserResult
	for userID, name := range userNames {
		count, exists := userCorrectCount[userID]
		if !exists {
			count = 0 // 맞춘 데이터가 없는 유저는 0으로 처리
		}

		users = append(users, response.UserResult{
			Name:              name,
			TotalCorrectCount: count,
		})
	}
	res.Users = users

	return res
}

func CreateResListGame(gameList []*mysql.Games) response.ResListGame {
	res := response.ResListGame{
		TotalCount: len(gameList),
	}
	for _, game := range gameList {
		g := response.GameInfo{
			Title:       game.Title,
			Description: game.Description,
			IsEnabled:   game.IsEnabled,
			YoutubeUrl:  game.YoutubeUrl,
			HashTag:     game.HashTag,
			Category:    game.Category,
		}
		// s3 에서 서명된 url로 응답
		imageUrl, err := _aws.ImageGetSignedURL(context.TODO(), game.Image, _aws.ImgTypeBoardGame)
		if err != nil {
			fmt.Println(err)
			continue
		}
		g.Image = imageUrl
		res.Games = append(res.Games, g)

	}
	return res
}

func CreateImageDTO(imageInfo request.ImageInfo) *mysql.FindItImages {
	return &mysql.FindItImages{
		Level:            imageInfo.Level,
		NormalImageUrl:   imageInfo.NormalImage,
		AbnormalImageUrl: imageInfo.AbnormalImage,
	}
}

func CreateImageCorrectDTO(imageID int, imageInfo request.ImageInfo) []*mysql.FindItImageCorrectPositions {
	imageCorrectDTOs := make([]*mysql.FindItImageCorrectPositions, 0)
	for _, position := range imageInfo.Positions {
		imageCorrectDTOs = append(imageCorrectDTOs, &mysql.FindItImageCorrectPositions{
			ImageID:   imageID,
			XPosition: position.X,
			YPosition: position.Y,
		})
	}
	return imageCorrectDTOs
}
