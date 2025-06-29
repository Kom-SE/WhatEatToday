package models

import (
	"time"

	"gorm.io/gorm"
)

type Recipe struct {
	gorm.Model
	Name        string    `gorm:"type:varchar(255);not null;column:name" json:"name"`    // Recipe name
	AuthorID    uint      `gorm:"type:varchar(50);not null;column:author" json:"author"` // Author of the recipe
	Description string    `gorm:"type:text;column:description" json:"description"`       // Description of the recipe
	Images      []Image   `gorm:"foreignKey:RecipeID"`
	FoodID      uint      `gorm:"not null;column:food_id;" json:"food_id"`
	CookTime    time.Time `gorm:"type:time;column:cook_time" json:"cook_time"`  // Associated images
	Process     string    `gorm:"type:text;column:process" json:"process"`      // Cooking process
	Likes       uint      `gorm:"type:int;default:0;column:likes" json:"likes"` // Number of likes

	ThisFoodID   Food `gorm:"foreignKey:FoodID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`   // Food associated with the recipe
	ThisAuthorID User `gorm:"foreignKey:AuthorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // User who created the recipe
}
