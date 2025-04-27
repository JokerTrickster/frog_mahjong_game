package entity

type SlimeWarRankEntity struct {
	UserID int `gorm:"column:user_id"`
	Score  int `gorm:"column:score"`
}
