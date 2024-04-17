package usecase

import (
	"context"
	_errors "main/features/room/model/errors"
	"main/features/room/model/interface/mocks"
	"main/features/room/model/response"
	"main/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"gopkg.in/go-playground/assert.v1"
)

// TestUserListRoomeUseCase_Discard 함수는 UserListRoomUseCase 의 UserList 메서드를 테스트합니다.
// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// given-when-then 패턴을 사용하여 테스트를 작성합니다.
// 매개변수는 roomID uint
// 응답은 response.ResUserListRoom, error
// 테스트 케이스:
// 1. roomID 가 33 인 경우 유저가 2명 존재항는 경우
// 2. roomID 가 34 인 경우 유저가 0명 존재하는 경우
// 테스트 경로: src/features/room/usecase/userListRoomUseCase_test.go
// 함수명 : TestUserListRoomUseCase_UserList

func TestUserListRoomUseCase_UserList(t *testing.T) {
	tests := []struct {
		name     string
		roomID   uint
		userList []response.User
		err      error
	}{
		{
			name:   "Test Case 1 success",
			roomID: 33,
			userList: []response.User{
				{
					UserID:         1,
					RoomUserID:     1,
					PlayerState:    "test1",
					TurnNumber:     1,
					OwnedCardCount: 0,
					RoomID:         33,
					Score:          1,
					UserName:       "test1",
					UserEmail:      "test1",
				},
				{
					UserID:         2,
					RoomUserID:     2,
					PlayerState:    "test2",
					TurnNumber:     2,
					OwnedCardCount: 0,
					RoomID:         33,
					Score:          2,
					UserName:       "test2",
					UserEmail:      "test2",
				},
			},
			err: nil,
		},
		{
			name:     "Test Case 2 fail",
			roomID:   34,
			userList: []response.User{},
			err:      utils.ErrorMsg(context.TODO(), utils.ErrNotFound, utils.Trace(), _errors.ErrRoomUserNotFound.Error(), utils.ErrFromClient),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//given
			mockRoomRepository := new(mocks.IUserListRoomRepository)
			mockRoomRepository.On("FindRoomUser", mock.Anything, mock.Anything).Return(tt.userList, tt.err)
			us := NewUserListRoomUseCase(mockRoomRepository, time.Second*8)
			//when
			_, err := us.UserList(context.TODO(), tt.roomID)
			//then
			assert.Equal(t, tt.err, err)
		})
	}
}
