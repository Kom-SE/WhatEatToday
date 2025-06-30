package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(50);not null;unique;index;column:username"`
	Password string `gorm:"type:varchar(255);not null;column:password"`
	Name     string `gorm:"type:varchar(50);not null;column:name"`
	Gender   string `gorm:"type:enum('M','F');not null;column:gender"`
	Phone    string `gorm:"type:varchar(11);not null;column:phone"`
	Address  string `gorm:"type:varchar(255);column:address"`
	Avatar   string `gorm:"type:varchar(500);column:avatar;default:''" json:"avatar"`
	UserType uint8  `gorm:"type:int;not null;default:1;column:user_type" json:"user_type"` // 0: Admin, 1: Normal User
}
