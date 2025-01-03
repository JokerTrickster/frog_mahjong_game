package v2ws

import (
	"main/utils/db/mysql"
	"math/rand"
	"strconv"
	"time"
)

// 무작위로 6자리 숫자로만 이루어진 비밀번호 생성
func CreateRandomPassword() string {
	rand.Seed(time.Now().UnixNano())       // 현재 시간을 시드로 설정
	randomNumber := rand.Intn(9000) + 1000 // 1000 ~ 9999 사이의 숫자를 생성
	return strconv.Itoa(randomNumber)      // 숫자를 문자열로 변환

}

func CreatePlayTogetherRoomDTO(uID uint, count int, timer int, password string) mysql.Rooms {
	result := mysql.Rooms{
		CurrentCount: 0,
		MaxCount:     count,
		MinCount:     count,
		State:        "wait",
		OwnerID:      int(uID),
		Timer:        timer,
		PlayTurn:     0,
		Name:         "play together",
		Password:     password,
		StartTime:    time.Now(),
		GameID:       2,
	}
	return result
}

func CreatePlayTogetherRoomUserDTO(uID uint, roomID int) mysql.RoomUsers {
	result := mysql.RoomUsers{
		UserID:         int(uID),
		RoomID:         roomID,
		Score:          0,
		OwnedCardCount: 0,
		PlayerState:    "play",
		TurnNumber:     0,
	}
	return result
}

func CreatePlayTogetherUserItemDTO(uID uint, roomID uint, item mysql.Items) mysql.UserItems {
	result := mysql.UserItems{
		UserID:        int(uID),
		RoomID:        int(roomID),
		ItemID:        int(item.ID),
		RemainingUses: item.MaxUses,
	}
	return result
}
