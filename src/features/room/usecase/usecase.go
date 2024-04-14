package usecase

import (
	"context"
	"main/features/room/model/request"
	"main/features/room/model/response"
	"main/utils"
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
		UserID:         int(uID),
		RoomID:         roomID,
		Score:          0,
		OwnedCardCount: 0,
		PlayerState:    playerState,
	}
	return result, nil
}

func CreateResListRoom(rooms []mysql.Rooms, total int) (response.ResListRoom, error) {
	res := response.ResListRoom{}
	RoomList := make([]response.ListRoom, 0, len(rooms))
	for i := 0; i < len(rooms); i++ {
		room := response.ListRoom{
			ID:           int(rooms[i].ID),
			CurrentCount: rooms[i].CurrentCount,
			MaxCount:     rooms[i].MaxCount,
			MinCount:     rooms[i].MinCount,
			Name:         rooms[i].Name,
			State:        rooms[i].State,
			Owner:        rooms[i].Owner,
			Created:      utils.TimeToEpochMillis(rooms[i].CreatedAt),
		}
		RoomList = append(RoomList, room)
	}
	res.Total = total
	res.Rooms = RoomList
	return res, nil
}
