package request

type ReqWSFailedLoan struct {
	CardID       uint `json:"cardID"`       //카드 정보
	TargetUserID uint `json:"targetUserID"` //카드 주인
	PlayTurn     int  `json:"playTurn"`     // 플레이턴
}
