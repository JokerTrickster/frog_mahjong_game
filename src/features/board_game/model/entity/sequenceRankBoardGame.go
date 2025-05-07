package entity

type SequenceRankEntity struct {
	UserID int `gorm:"column:user_id"`
	Score  int `gorm:"column:score"`
}
