package _errors

import "net/http"

var (
	// 에러 타입 정의
	ErrWrongPassword               = "ERR_WRONG_PASSWORD"
	ErrInternalServer              = "ERR_INTERNAL_SERVER"
	ErrGameInProgress              = "ERR_GAME_IN_PROGRESS"
	ErrRoomFull                    = "ERR_ROOM_FULL"
	ErrBadRequest                  = "ERR_BAD_REQUEST"
	ErrAbnormalExit                = "ERR_ABNORMAL_EXIT"
	ErrCryptoFailed                = "ERR_CRYPTO_FAILED"
	ErrUnmarshalFailed             = "ERR_UNMARSHAL_FAILED"
	ErrHintItemFailed              = "ERR_HINT_ITEM_FAILED"
	ErrTimerItemFailed             = "ERR_TIMER_ITEM_FAILED"
	ErrMarshalFailed               = "ERR_MARSHAL_FAILED"
	ErrDisconnect                  = "ERR_DISCONNECT"
	ErrDBServer                    = "ERR_DB_SERVER"
	ErrInvalidToken                = "ERR_INVALID_TOKEN"
	ErrRoomOut                     = "ERR_ROOM_OUT"
	ErrNotFoundCard                = "ERR_NOT_FOUND_CARD"
	ErrItemNotAvailable            = "ERR_ITEM_NOT_AVAILABLE"
	ErrRoomUsersNotFound           = "ERR_ROOM_USERS_NOT_FOUND"
	ErrDeleteCardFailed            = "ERR_DELETE_CARD_FAILED"
	ErrDeleteRoomFailed            = "ERR_DELETE_ROOM_FAILED"
	ErrDeleteMapFailed             = "ERR_DELETE_MAP_FAILED"
	ErrDeleteGameRoomSettingFailed = "ERR_DELETE_GAME_ROOM_SETTING_FAILED"
	ErrDeleteRoomUserFailed        = "ERR_DELETE_ROOM_USER_FAILED"
	ErrUpdateUserStateFailed       = "ERR_UPDATE_USER_STATE_FAILED"
	ErrUserNotFound                = "ERR_USER_NOT_FOUND"
	ErrRoomNotFound                = "ERR_ROOM_NOT_FOUND"
	ErrUpdateFailed                = "ERR_UPDATE_FAILED"
	ErrRoundIncreaseFailed         = "ERR_ROUND_INCREASE_FAILED"
	ErrDeleteFailed                = "ERR_DELETE_FAILED"
	ErrInvalidRequest              = "ERR_INVALID_REQUEST"
	ErrCreateFailed                = "ERR_CREATE_FAILED"
	ErrCountFailed                 = "ERR_COUNT_FAILED"
	ErrUnauthorizedAction          = "ERR_UNAUTHORIZED_ACTION"
	ErrFetchFailed                 = "ERR_FETCH_FAILED"
	ErrAlreadyGame                 = "ERR_ALREADY_GAME"
	ErrGameTerminated              = "ERR_GAME_TERMINATED" // 게임 강제 중단
	ErrSlimeWarUsersNotFound       = "ERR_SLIME_WAR_USERS_NOT_FOUND"
	ErrSlimeWarUserUpdateFailed    = "ERR_SLIME_WAR_USER_UPDATE_FAILED"
)

var (
	ErrCodeInternal         = http.StatusInternalServerError // 500
	ErrCodeNotFound         = http.StatusNotFound            // 404
	ErrCodeBadRequest       = http.StatusBadRequest          // 400
	ErrCodeForbidden        = http.StatusForbidden           // 403
	ErrCodeInvalidUserCount = http.StatusBadRequest          // 400
)
