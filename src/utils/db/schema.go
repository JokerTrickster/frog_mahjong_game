package db

type Users struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Score    int    `json:"score"`
	State    string `json:"state"`
	RoomID   int    `json:"roomID"`
}
type Rooms struct {
	ID           int    `json:"id"`
	CurrentCount int    `json:"currentCount"` //방 현재 인원
	MaxCount     int    `json:"maxCount"`     //방 최대 인원
	MinCount     int    `json:"minCount"`     //방 최소 인원
	Name         string `json:"name"`         //방 이름
	Password     string `json:"password"`     //방 비밀번호 (옵셔널))
	State        string `json:"state"`        //방 상태 (대기, 진행, 종료)
	Owner        string `json:"owner"`        //방 주인
}

type RoomUsers struct {
	ID          int    `json:"id"`
	UserID      int    `json:"userID"`
	RoomID      int    `json:"roomID"`
	Score       int    `json:"score"`
	CardCount   int    `json:"cardCount"`
	PlayerState string `json:"playerState"` // wait, ready, play, end
}

type Cards struct {
	ID     int    `json:"id"`
	RoomID int    `json:"roomID"` //방 아이디 (어느 방에 있는 카드인지)
	UserID int    `json:"userID"` //유저 아이디 (소유하고 있는 유저 아이디)
	Name   string `json:"name"`   //카드 이름 (1,2,3,4,5,6,7,8,9,중,발)
	Color  string `json:"color"`  //카드 색깔 (레드, 그린, 일반)
	State  string `json:"state"`  //카드 상태 (바닥에 놓여있다, 손에 들고 있다, 유저 앞에 놓여있다.)
}
