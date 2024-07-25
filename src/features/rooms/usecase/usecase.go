package usecase

import (
	"context"
	"main/features/rooms/model/request"
	"main/features/rooms/model/response"
	"main/utils"
	"main/utils/db/mysql"
)

func CreateResUserListRoom(userList []response.User, rooms mysql.Rooms) response.ResUserListRoom {
	res := response.ResUserListRoom{}
	for i := 0; i < len(userList); i++ {
		if userList[i].UserID == rooms.OwnerID {
			userList[i].Owner = true
		} else {
			userList[i].Owner = false
		}
	}
	res.Users = userList
	return res
}

func CreateRoomDTO(ctx context.Context, req *request.ReqCreate, uID uint) (mysql.Rooms, error) {

	result := mysql.Rooms{
		CurrentCount: 1,
		MaxCount:     req.MaxCount,
		MinCount:     req.MinCount,
		Name:         req.Name,
		State:        "wait",
		OwnerID:      int(uID),
	}
	if req.Password != "" {
		result.Password = req.Password
	}
	return result, nil
}

func V02CreateRoomDTO(ctx context.Context, req *request.ReqV02Create, uID uint) (mysql.Rooms, error) {

	result := mysql.Rooms{
		CurrentCount: 1,
		MaxCount:     2,
		MinCount:     2,
		Name:         "test",
		State:        "wait",
		OwnerID:      int(uID),
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

func V02CreateRoomUserDTO(uID uint, roomID int, playerState string) (mysql.RoomUsers, error) {
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
			OwnerID:      rooms[i].OwnerID,
			Password:     false,
			Created:      utils.TimeToEpochMillis(rooms[i].CreatedAt),
		}
		if rooms[i].Password != "" {
			room.Password = true
		}
		RoomList = append(RoomList, room)
	}
	res.Total = total
	res.Rooms = RoomList
	return res, nil
}
