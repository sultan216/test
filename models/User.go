package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	FullName string `json:"fullname"`
	Post     []Post `gorm:"foreignKey:UserRefer"`
}
