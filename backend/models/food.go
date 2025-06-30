package models

import "gorm.io/gorm"

type Food struct {
	gorm.Model
	Name        string `gorm:"type:varchar(255);not null;index;column:name" json:"type:name"`
	Description string `gorm:"type:varchar(255)" json:"description"`
	Image       string `gorm:"type:varchar(500);column:image;default:''" json:"image"` // 食材图片
}
