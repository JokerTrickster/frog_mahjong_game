// Code generated by mockery v2.42.2. DO NOT EDIT.

package mocks

import (
	context "context"
	mysql "main/utils/db/mysql"

	mock "github.com/stretchr/testify/mock"
)

// IWinRequestGameRepository is an autogenerated mock type for the IWinRequestGameRepository type
type IWinRequestGameRepository struct {
	mock.Mock
}

// GetRoomUser provides a mock function with given fields: c, userID, roomID
func (_m *IWinRequestGameRepository) GetRoomUser(c context.Context, userID uint, roomID uint) (mysql.RoomUsers, error) {
	ret := _m.Called(c, userID, roomID)

	if len(ret) == 0 {
		panic("no return value specified for GetRoomUser")
	}

	var r0 mysql.RoomUsers
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, uint) (mysql.RoomUsers, error)); ok {
		return rf(c, userID, roomID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint, uint) mysql.RoomUsers); ok {
		r0 = rf(c, userID, roomID)
	} else {
		r0 = ret.Get(0).(mysql.RoomUsers)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint, uint) error); ok {
		r1 = rf(c, userID, roomID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIWinRequestGameRepository creates a new instance of IWinRequestGameRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIWinRequestGameRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *IWinRequestGameRepository {
	mock := &IWinRequestGameRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}