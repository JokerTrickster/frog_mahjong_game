package _errors

import "net/http"

var (
	ErrWrongPassword  = "ERR_WRONG_PASSWORD"
	ErrInternalServer = "ERR_INTERNAL_SERVER"
	ErrGameInProgress = "ERR_GAME_IN_PROGRESS"
	ErrRoomFull       = "ERR_ROOM_FULL"
	ErrBadRequest     = "ERR_BAD_REQUEST"
	ErrAbnormalExit   = "ERR_ABNORMAL_EXIT"
	ErrDBServer       = "ERR_DB_SERVER"
	ErrInvalidToken   = "ERR_INVALID_TOKEN"
	ErrRoomOut        = "ROOM_OUT"
	ErrNotFoundCard   = "ERR_NOT_FOUND_CARD"
)

var (
	ErrCodeBadRequest = http.StatusBadRequest
	ErrCodeInternal   = http.StatusInternalServerError
)
