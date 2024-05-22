package usecase

import (
	"context"
	_errors "main/features/users/model/errors"
	"main/features/users/model/interface/mocks"

	"main/utils"
	"main/utils/db/mysql"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/gorm"
)

// TestListUsersUseCase_List 함수는 GetUsersUseCase 의 List 메서드를 테스트합니다.
// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// given-when-then 패턴을 사용하여 테스트를 작성합니다.
// 구조체 필드 : name string, userDTO []mysql.Users, total int, wantErr error
// 반환값 : response.ResListUser, error
// 테스트 케이스:
// 1. userID 가 존재해서 정상적으로 응답한 경우
// 2. userID 가 존재하지 않아서 에러가 발생한 경우
// 테스트 경로: src/features/users/usecase/listUsersUseCase_test.go

func TestListUsersUseCase_List(t *testing.T) {
	tests := []struct {
		name     string
		userList []mysql.Users
		total    int
		wantErr  error
	}{
		{
			name: "userID 가 존재해서 정상적으로 응답한 경우",
			userList: []mysql.Users{
				{
					Model: gorm.Model{
						ID: 1,
					},
					Name:  "test",
					Email: "test@gmail.com",
				},
			},
			total:   1,
			wantErr: nil,
		},
		{
			name:     "userID 가 존재하지 않아서 에러가 발생한 경우",
			userList: []mysql.Users{},
			total:    0,
			wantErr:  utils.ErrorMsg(context.TODO(), utils.ErrBadParameter, utils.Trace(), _errors.ErrUserNotFound.Error(), utils.ErrFromClient),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			mockListUsersRepository := new(mocks.IListUsersRepository)
			mockListUsersRepository.On("FindUsers", mock.Anything).Return(tt.userList, tt.wantErr) //mock
			mockListUsersRepository.On("CountUsers", mock.Anything).Return(tt.total, tt.wantErr)   //mock
			us := NewListUsersUseCase(mockListUsersRepository, 8*time.Second)

			// when
			_, err := us.List(context.TODO())

			// then
			assert.Equal(t, tt.wantErr, err)

		})
	}
}
