package models

import "gorm.io/gorm"

type Label struct {
	gorm.Model
	Name        string `gorm:"type:varchar(255);not null;unique;index;column:name" json:"name"`
	Type        uint8  `gorm:"type:int;column:type" json:"type"`                        //null:无 1:口味 2:做法
	Description string `gorm:"type:varchar(255);column:description" json:"description"` // 标签具体内容：null则没有，1则口味，2则做法
}
