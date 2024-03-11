package mysql

import "gorm.io/gorm"

type Users struct {
	gorm.Model 
	Name     string `json:"name"  gorm:"column:name"`
	Email    string `json:"email"  gorm:"column:email"`
	Password string `json:"password"  gorm:"column:password"`
}

/*
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);
*/
