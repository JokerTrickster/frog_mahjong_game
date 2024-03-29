// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	context "context"
	mysql "main/utils/db/mysql"

	mock "github.com/stretchr/testify/mock"
)

// ICreateRoomRepository is an autogenerated mock type for the ICreateRoomRepository type
type ICreateRoomRepository struct {
	mock.Mock
}

// FindOneAndUpdateUser provides a mock function with given fields: ctx, uID, roomID
func (_m *ICreateRoomRepository) FindOneAndUpdateUser(ctx context.Context, uID uint, roomID uint) error {
	ret := _m.Called(ctx, uID, roomID)

	if len(ret) == 0 {
		panic("no return value specified for FindOneAndUpdateUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, uint) error); ok {
		r0 = rf(ctx, uID, roomID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InsertOneRoom provides a mock function with given fields: ctx, roomDTO
func (_m *ICreateRoomRepository) InsertOneRoom(ctx context.Context, roomDTO mysql.Rooms) (int, error) {
	ret := _m.Called(ctx, roomDTO)

	if len(ret) == 0 {
		panic("no return value specified for InsertOneRoom")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, mysql.Rooms) (int, error)); ok {
		return rf(ctx, roomDTO)
	}
	if rf, ok := ret.Get(0).(func(context.Context, mysql.Rooms) int); ok {
		r0 = rf(ctx, roomDTO)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, mysql.Rooms) error); ok {
		r1 = rf(ctx, roomDTO)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertOneRoomUser provides a mock function with given fields: ctx, roomUserDTO
func (_m *ICreateRoomRepository) InsertOneRoomUser(ctx context.Context, roomUserDTO mysql.RoomUsers) error {
	ret := _m.Called(ctx, roomUserDTO)

	if len(ret) == 0 {
		panic("no return value specified for InsertOneRoomUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, mysql.RoomUsers) error); ok {
		r0 = rf(ctx, roomUserDTO)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewICreateRoomRepository creates a new instance of ICreateRoomRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewICreateRoomRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ICreateRoomRepository {
	mock := &ICreateRoomRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
