package entity

type FindItRankEntity struct {
	UserID       int `gorm:"column:user_id"`
	CorrectCount int `gorm:"column:correct_count"`
}
