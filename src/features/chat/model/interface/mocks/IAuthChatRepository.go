// Code generated by mockery v2.42.2. DO NOT EDIT.

package mocks

import (
	context "context"
	mysql "main/utils/db/mysql"

	mock "github.com/stretchr/testify/mock"
)

// IAuthChatRepository is an autogenerated mock type for the IAuthChatRepository type
type IAuthChatRepository struct {
	mock.Mock
}

// FindOneUserInfo provides a mock function with given fields: ctx, userID
func (_m *IAuthChatRepository) FindOneUserInfo(ctx context.Context, userID uint) (*mysql.Users, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for FindOneUserInfo")
	}

	var r0 *mysql.Users
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint) (*mysql.Users, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint) *mysql.Users); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mysql.Users)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertOneChat provides a mock function with given fields: ctx, chatDTO
func (_m *IAuthChatRepository) InsertOneChat(ctx context.Context, chatDTO *mysql.Chats) error {
	ret := _m.Called(ctx, chatDTO)

	if len(ret) == 0 {
		panic("no return value specified for InsertOneChat")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *mysql.Chats) error); ok {
		r0 = rf(ctx, chatDTO)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewIAuthChatRepository creates a new instance of IAuthChatRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIAuthChatRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *IAuthChatRepository {
	mock := &IAuthChatRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}