package db

type UserDTO struct {
	ID        string `json:"id`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"` //생성 날짜
	UpdatedAt string `json:"updatedAt"`
	IsDeleted bool   `json:"isDeleted"` //활동 여부
}

type UserAuthDTO struct {
	ID         string `json:"id"`
	Provider   string `json:"provider"`
	UserID     string `json:"userID"`
	Password   string `json:"password`
	LastSignIn string `json:"lastSignIn"`
	IsDeleted  bool   `json:"isDeleted"` //활동 여부
	UpdatedAt  string `json:"updatedAt"`
	CreatedAt  string `json:"createdAt" ` //생성 날짜
}
