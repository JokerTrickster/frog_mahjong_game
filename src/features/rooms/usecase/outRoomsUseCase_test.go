package usecase

/*
TODO 트랜잭션 테스트 코드 어떻게 처리할지 고민 필요

func TestOutRoomsUseCase_Out(t *testing.T) {
	testCases := []struct {
		name        string
		uID         uint
		roomDTO     mysql.Rooms
		roomUserDTO mysql.RoomUsers
		userDTO     mysql.Users
		req         *request.ReqOut
		err         error
	}{
		{
			name: "Test Case 1 Success 방 인원이 0명인 경우",
			uID:  1,
			roomDTO: mysql.Rooms{
				Model: gorm.Model{
					ID: 2,
				},
				CurrentCount: 0,
			},
			roomUserDTO: mysql.RoomUsers{},
			userDTO:     mysql.Users{},
			req:         &request.ReqOut{},
			err:         nil,
		},
		{
			name: "Test Case 2 Success 방 인원이 1명인 경우",
			uID:  2,
			roomDTO: mysql.Rooms{
				Model: gorm.Model{
					ID: 3,
				},
				CurrentCount: 1,
			},
			roomUserDTO: mysql.RoomUsers{
				UserID: 2,
				RoomID: 3,
			},
			userDTO: mysql.Users{
				Model: gorm.Model{
					ID: 2,
				},
				RoomID: 3,
				Email:  "test@gmail.com",
			},
			req: &request.ReqOut{},
			err: nil,
		},
		{
			name: "Test Case 2 Success 방 인원이 3명인 경우",
			uID:  2,
			roomDTO: mysql.Rooms{
				Model: gorm.Model{
					ID: 3,
				},
				CurrentCount: 1,
			},
			roomUserDTO: mysql.RoomUsers{
				UserID: 2,
				RoomID: 3,
			},
			userDTO: mysql.Users{
				Model: gorm.Model{
					ID: 2,
				},
				RoomID: 3,
				Email:  "test@gmail.com",
			},
			req: &request.ReqOut{},
			err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//given
			mockOutRoomRepository := new(mocks.IOutRoomsRepository)

			mockOutRoomRepository.On("FindOneAndDeleteRoomUser", mock.Anything, mock.Anything, mock.Anything).Return(tc.err) //mock
			mockOutRoomRepository.On("FindOneAndUpdateRoom", mock.Anything, mock.Anything).Return(tc.roomDTO, nil)
			mockOutRoomRepository.On("FindOneAndUpdateUser", mock.Anything, mock.Anything).Return(nil)
			if tc.roomDTO.CurrentCount == 0 {
				mockOutRoomRepository.On("FindOneAndDeleteRoom", mock.Anything, mock.Anything).Return(nil)
			} else {
				mockOutRoomRepository.On("FindOneRoomUser", mock.Anything, mock.Anything).Return(tc.roomUserDTO, nil)
				mockOutRoomRepository.On("FindOneUser", mock.Anything, mock.Anything).Return(tc.userDTO, nil)
				mockOutRoomRepository.On("ChangeRoomOnwer", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			}

			us := NewOutRoomsUseCase(mockOutRoomRepository, 8*time.Second)

			//when
			err := us.Out(context.TODO(), tc.uID, tc.req)

			//then
			assert.Equal(t, tc.err, err)
		})
	}
}
*/
