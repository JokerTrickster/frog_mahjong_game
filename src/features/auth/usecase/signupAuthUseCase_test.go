package usecase

// Write code to test the signupAuthUseCase.go file

import (
	"context"
	"main/features/auth/model/interface/mocks"
	"main/features/auth/model/request"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"gopkg.in/go-playground/assert.v1"
)

func TestSignupAuthUseCase_Signup(t *testing.T) {
	mockSignupAuthRepository := new(mocks.ISignupAuthRepository)
	tests := []struct {
		name string
		req  request.ReqSignup
		want error
	}{
		{"success1", request.ReqSignup{Email: "test01@test.com", Password: "1234"}, nil},
		{"success2", request.ReqSignup{Email: "test02@test.com", Password: "1234"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//given
			mockSignupAuthRepository.On("FindOneUserAuth", mock.Anything, mock.Anything).Return(nil)       //mock
			mockSignupAuthRepository.On("InsertOneUserDTO", mock.Anything, mock.Anything).Return("1", nil) //mock
			mockSignupAuthRepository.On("InsertOneUserAuthDTO", mock.Anything, mock.Anything).Return(nil)  //mock
			us := NewSignupAuthUseCase(mockSignupAuthRepository, 8*time.Second)

			//when
			got := us.Signup(context.TODO(), &tt.req)
			//then
			assert.Equal(t, tt.want, got)
		})
	}
}
