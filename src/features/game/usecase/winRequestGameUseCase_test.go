package usecase

// TestWinRequestGameUseCase_WinRequest 함수는 WinRequestGameUseCase 의 WinRequest 메서드를 테스트합니다.
// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// given-when-then 패턴을 사용하여 테스트를 작성합니다.
// 테스트 매개변수 : req *request.ReqWinRequest, mysql.RoomUsers, error
// 테스트 케이스:
// - 유저가 소유하고 있는 카드가 6장이고 플레이 상태가 play인 경우
// - 유저가 소유하고 있는 카드가 6장이고 플레이 상태가 loan인 경우
// - 유저가 소유하고 있는 카드가 5장이고 플레이 상태가 play-wait인 경우
// 테스트 경로: src/features/game/usecase/winRequestGameUseCase_test.go

import (
	"context"
	"main/features/game/model/interface/mocks"
	"main/features/game/model/request"
	"main/utils/db/mysql"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"gopkg.in/go-playground/assert.v1"
)

// 테이블 주도 테스트 방법을 사용하여 여러 시나리오를 테스트합니다.
func TestWinRequestGameUseCase_WinRequest(t *testing.T) {
	//테스트 케이스를 정의합니다.
	tests := []struct {
		name          string
		req           *request.ReqWinRequest
		roomUser      mysql.RoomUsers
		expectedValue bool
		expectedError error
	}{
		{
			name: "유저가 소유하고 있는 카드가 6장이고 플레이 상태가 play인 경우",
			req: &request.ReqWinRequest{
				UserID: 1,
				RoomID: 1,
			},
			roomUser: mysql.RoomUsers{
				UserID:         1,
				RoomID:         1,
				PlayerState:    "play",
				OwnedCardCount: 6,
			},
			expectedValue: true,
			expectedError: nil,
		},
		{
			name: "유저가 소유하고 있는 카드가 6장이고 플레이 상태가 loan인 경우",
			req: &request.ReqWinRequest{
				UserID: 1,
				RoomID: 1,
			},
			roomUser: mysql.RoomUsers{
				UserID:         1,
				RoomID:         1,
				PlayerState:    "loan",
				OwnedCardCount: 6,
			},
			expectedValue: true,
			expectedError: nil,
		},
		{
			name: "유저가 소유하고 있는 카드가 5장이고 플레이 상태가 play-wait인 경우",
			req: &request.ReqWinRequest{
				UserID: 1,
				RoomID: 1,
			},
			roomUser: mysql.RoomUsers{
				UserID:         1,
				RoomID:         1,
				PlayerState:    "play-wait",
				OwnedCardCount: 5,
			},
			expectedValue: false,
			expectedError: nil,
		},
	}

	//테스트 케이스를 순회하며 테스트를 수행합니다.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			//모킹 객체를 생성합니다.
			mockRepo := new(mocks.IWinRequestGameRepository)
			timeout := time.Second * 5
			uc := NewWinRequestGameUseCase(mockRepo, timeout)

			//모킹 객체의 GetRoomUser 메서드를 설정합니다.
			mockRepo.On("GetRoomUser", mock.Anything, tt.req.UserID, tt.req.RoomID).Return(tt.roomUser, tt.expectedError)

			//WinRequest 메서드를 호출합니다.
			result, err := uc.WinRequest(context.Background(), tt.req)

			//결과를 검증합니다.
			assert.Equal(t, result, tt.expectedValue)
			assert.Equal(t, err, tt.expectedError)

			//모킹 객체의 GetRoomUser 메서드가 호출되었는지 검증합니다.
			mockRepo.AssertExpectations(t)
		})
	}
}
