package sequence

import (
	"main/features/sequence/model/request"
	"main/utils/db/mysql"
)

func CreateGameResultDTO(roomUser mysql.GameRoomUsers, roomID uint, req request.ReqGameOverEvent) mysql.GameResults {
	gameResult := 0
	if roomUser.UserID == int(req.WinnerID) {
		gameResult = 1
	} else {
		gameResult = 0
	}
	result := mysql.GameResults{
		RoomID:   int(roomID),
		UserID:   int(roomUser.UserID),
		Result:   gameResult,
		GameType: mysql.SEQUENCE,
	}

	return result
}
