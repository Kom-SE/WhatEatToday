package models

import "gorm.io/gorm"

type Image struct {
	gorm.Model
	RecipeID uint   `gorm:"type:int unsigned;not null;index;column:recipe_id"`
	ImageURL string `gorm:"type:varchar(255);not null;column:image_url" json:"image_url"` // URL of the image
}
