// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	context "context"
	mysql "main/utils/db/mysql"

	mock "github.com/stretchr/testify/mock"
)

// IStartGameRepository is an autogenerated mock type for the IStartGameRepository type
type IStartGameRepository struct {
	mock.Mock
}

// CheckOwner provides a mock function with given fields: c, email, roomID
func (_m *IStartGameRepository) CheckOwner(c context.Context, email string, roomID uint) error {
	ret := _m.Called(c, email, roomID)

	if len(ret) == 0 {
		panic("no return value specified for CheckOwner")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, uint) error); ok {
		r0 = rf(c, email, roomID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CheckReady provides a mock function with given fields: c, roomID
func (_m *IStartGameRepository) CheckReady(c context.Context, roomID uint) ([]mysql.RoomUsers, error) {
	ret := _m.Called(c, roomID)

	if len(ret) == 0 {
		panic("no return value specified for CheckReady")
	}

	var r0 []mysql.RoomUsers
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint) ([]mysql.RoomUsers, error)); ok {
		return rf(c, roomID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint) []mysql.RoomUsers); ok {
		r0 = rf(c, roomID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]mysql.RoomUsers)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(c, roomID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateCards provides a mock function with given fields: c, roomID, cards
func (_m *IStartGameRepository) CreateCards(c context.Context, roomID uint, cards []mysql.Cards) error {
	ret := _m.Called(c, roomID, cards)

	if len(ret) == 0 {
		panic("no return value specified for CreateCards")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, []mysql.Cards) error); ok {
		r0 = rf(c, roomID, cards)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateRoom provides a mock function with given fields: c, roomID, state
func (_m *IStartGameRepository) UpdateRoom(c context.Context, roomID uint, state string) error {
	ret := _m.Called(c, roomID, state)

	if len(ret) == 0 {
		panic("no return value specified for UpdateRoom")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, string) error); ok {
		r0 = rf(c, roomID, state)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateRoomUser provides a mock function with given fields: c, roomID, state
func (_m *IStartGameRepository) UpdateRoomUser(c context.Context, roomID uint, state string) error {
	ret := _m.Called(c, roomID, state)

	if len(ret) == 0 {
		panic("no return value specified for UpdateRoomUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, string) error); ok {
		r0 = rf(c, roomID, state)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewIStartGameRepository creates a new instance of IStartGameRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIStartGameRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *IStartGameRepository {
	mock := &IStartGameRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}