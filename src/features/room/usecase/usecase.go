package usecase

import (
	"main/features/room/model/request"
	"main/utils/db/mysql"
)

func CreateRoomDTO(req *request.ReqCreate, email string) mysql.Rooms {
	result := mysql.Rooms{
		CurrentCount: 0,
		MaxCount:     req.MaxCount,
		MinCount:     req.MinCount,
		Name:         req.Name,
		State:        "wait",
		Owner:        email,
	}
	if req.Password != "" {
		result.Password = req.Password
	}
	return result
}

func CreateRoomUserDTO(uID uint, roomID int) mysql.RoomUsers {
	return mysql.RoomUsers{
		UserID:      int(uID),
		RoomID:      roomID,
		Score:       0,
		CardCount:   0,
		PlayerState: "ready",
	}
}
