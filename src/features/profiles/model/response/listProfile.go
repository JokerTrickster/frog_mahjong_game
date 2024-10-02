package response

type ResListProfile struct {
	ProfileList []Profile `json:"profileList"`
}

type Profile struct {
	ProfileID    int  `json:"profileID"`    //프로필 ID
	IsAchieved   bool `json:"isAchieved"`   //획득 여부
	CurrentCount int  `json:"currentCount"` //현재 획득 포인트
}
