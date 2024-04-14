package usecase

import (
	"context"
	"main/features/room/model/interface/mocks"
	"main/features/room/model/response"
	"main/utils"
	"main/utils/db/mysql"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/gorm"
)

// TestListRoomUseCase_ListRoom 함수는 ListRoomUseCase 의 ListRoom 메서드를 테스트합니다.
// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// given-when-then 패턴을 사용하여 테스트를 작성합니다.
// 테스트 매개변수 : page, pageSize, response.ResListRoom, error
// 테스트 케이스:
// 1. page: 0, pageSize: 10, 방 목록이 1개인 경우
// 2. page: 0, pageSize: 10, 방 목록이 0개인 경우
// 테스트 경로: src/features/room/usecase/ListRoomUseCase_test.go

func TestListRoomUseCase_List(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		page     int
		pageSize int
		rooms    []mysql.Rooms
		res      *response.ResListRoom
		err      error
	}{
		{
			name:     "유저가 1명인 방 목록 조회 성공",
			page:     0,
			pageSize: 10,
			rooms: []mysql.Rooms{
				{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: now,
					},
					Name:         "test",
					Owner:        "ryan",
					CurrentCount: 1,
					MaxCount:     2,
					MinCount:     2,
					State:        "wait",
				},
			},
			res: &response.ResListRoom{
				Total: 1,
				Rooms: []response.ListRoom{
					{
						ID:           1,
						Name:         "test",
						Owner:        "ryan",
						CurrentCount: 1,
						MaxCount:     2,
						MinCount:     2,
						State:        "wait",
						Created:      utils.TimeToEpochMillis(now),
					},
				},
			},
			err: nil,
		},
		{
			name:     "방 목록이 없는 경우 조회 성공",
			page:     0,
			pageSize: 10,
			rooms:    []mysql.Rooms{},
			res: &response.ResListRoom{
				Total: 0,
				Rooms: []response.ListRoom{},
			},
			err: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// given
			mockRoomRepo := new(mocks.IListRoomRepository)
			mockRoomRepo.On("FindRoomList", mock.Anything, mock.Anything, mock.Anything).Return(tc.rooms, tc.err)
			mockRoomRepo.On("CountRoomList", mock.Anything).Return(tc.res.Total, tc.err)
			uc := NewListRoomUseCase(mockRoomRepo, 8*time.Second)

			// when
			res, err := uc.List(context.Background(), tc.page, tc.pageSize)

			// then
			assert.Equal(t, tc.res, res)
			assert.Equal(t, tc.err, err)
		})
	}
}

// TestCreateResListRoom 함수는 ListRoomUseCase 의 CreateResListRoom 메서드를 테스트합니다.
// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// given-when-then 패턴을 사용하여 테스트를 작성합니다.
// 테스트 케이스:
// 1. 방 목록이 1개인 경우
// 2. 방 목록이 0개인 경우
// 테스트 경로: src/features/room/usecase/listRoomUseCase_test.go

func TestCreateResListRoom(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name  string
		rooms []mysql.Rooms
		total int
		res   *response.ResListRoom
		err   error
	}{
		{
			name: "방 목록이 1개인 경우",
			rooms: []mysql.Rooms{
				{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: now,
					},
					Name:         "test",
					Owner:        "ryan",
					CurrentCount: 1,
					MaxCount:     2,
					MinCount:     2,
					State:        "wait",
				},
			},
			total: 1,
			res: &response.ResListRoom{
				Total: 1,
				Rooms: []response.ListRoom{
					{
						ID:           1,
						Name:         "test",
						Owner:        "ryan",
						CurrentCount: 1,
						MaxCount:     2,
						MinCount:     2,
						State:        "wait",
						Created:      utils.TimeToEpochMillis(now),
					},
				},
			},
			err: nil,
		},
		{
			name:  "방 목록이 0개인 경우",
			rooms: []mysql.Rooms{},
			total: 0,
			res: &response.ResListRoom{
				Total: 0,
				Rooms: []response.ListRoom{},
			},
			err: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// when
			res, err := CreateResListRoom(tc.rooms, tc.total)

			// then
			assert.Equal(t, tc.res, res)
			assert.Equal(t, tc.err, err)
		})
	}
}
