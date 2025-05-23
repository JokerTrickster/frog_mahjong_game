package ws

import (
	"main/utils/db/mysql"
	"math/rand"
	"strconv"
	"time"
)

// 무작위로 4자리 숫자로만 이루어진 비밀번호 생성
func CreateRandomPassword() string {
	rand.Seed(time.Now().UnixNano())       // 현재 시간을 시드로 설정
	randomNumber := rand.Intn(9000) + 1000 // 100000 ~ 999999 사이의 숫자를 생성
	return strconv.Itoa(randomNumber)      // 숫자를 문자열로 변환

}

func CreatePlayTogetherRoomDTO(uID uint, count int, timer int, password string) mysql.GameRooms {
	result := mysql.GameRooms{
		CurrentCount: 0,
		MaxCount:     count,
		MinCount:     count,
		State:        "wait",
		OwnerID:      int(uID),
		Name:         "play together",
		Password:     password,
		GameID:       mysql.FROG,
		StartTime:    time.Now(),
	}
	return result
}

func CreatePlayTogetherRoomUserDTO(uID uint, roomID int, playerState string) mysql.FrogRoomUsers {
	result := mysql.FrogRoomUsers{
		UserID:         int(uID),
		RoomID:         roomID,
		Score:          0,
		OwnedCardCount: 0,
		PlayerState:    playerState,
		TurnNumber:     0,
	}
	return result
}
