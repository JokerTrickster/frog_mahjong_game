package request

// 론 카드 ID, 상대방 유저 ID
type ReqWSLoan struct {
	CardID       uint `json:"cardID"`       // 론 카드 ID
	TargetUserID uint `json:"targetUserID"` // 상대방 유저 ID
}
