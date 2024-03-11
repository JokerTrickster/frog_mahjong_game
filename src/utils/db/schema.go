package db

type Users struct {
	UserID    int    `json:"user_id"`
	UserName  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"` //생성 날짜
	UpdatedAt string `json:"updated_at"`
	IsDeleted bool   `json:"is_deleted"` //활동 여부
}

/*
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN NOT NULL
);
*/
