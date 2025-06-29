package models

import "gorm.io/gorm"

type Food struct {
	gorm.Model
	Name string `gorm:"type:varchar(255);not null;unique;index;column:name" json:"type:name"`
}
