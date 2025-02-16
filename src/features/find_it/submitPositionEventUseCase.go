package find_it

import "main/utils/db/mysql"

func CreateUserCorrectPosition(roomID, userID uint, round, imageID, correctID int) *mysql.FindItUserCorrectPositions {
	return &mysql.FindItUserCorrectPositions{
		RoomID:            int(roomID),
		UserID:            int(userID),
		Round:             round,
		ImageID:           imageID,
		CorrectPositionID: correctID,
	}
}
