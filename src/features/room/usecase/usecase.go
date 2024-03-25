package usecase

import (
	"context"
	"fmt"
	"main/features/room/model/request"
	"main/utils/db/mysql"
)

func CreateRoomDTO(ctx context.Context, req *request.ReqCreate, email string) (mysql.Rooms, error) {

	result := mysql.Rooms{
		CurrentCount: 1,
		MaxCount:     req.MaxCount,
		MinCount:     req.MinCount,
		Name:         req.Name,
		State:        "wait",
		Owner:        email,
	}
	if req.Password != "" {
		result.Password = req.Password
	}
	return result, nil
}

func CreateRoomUserDTO(uID uint, roomID int, playerState string) (mysql.RoomUsers, error) {
	result := mysql.RoomUsers{
		UserID:      int(uID),
		RoomID:      roomID,
		Score:       0,
		CardCount:   0,
		PlayerState: playerState,
	}
	fmt.Println(result)
	return result, nil
}
