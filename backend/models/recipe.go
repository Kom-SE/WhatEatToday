package models

import (
	"gorm.io/gorm"
)

type Recipe struct {
	gorm.Model
	Title          string `gorm:"type:varchar(255);not null;column:title" json:"title"` // Recipe name
	AuthorID       uint   `gorm:"type:int;not null;column:author_id" json:"author_id"`  // Author of the recipe
	Description    string `gorm:"type:text;column:description" json:"description"`      // Description of the recipe
	Images         string `gorm:"type:text;column:images" json:"images"`                // Image URLs associated with the recipe
	FoodID         string `gorm:"column:food_id;" json:"food_id"`
	Servings       string `gorm:"type:varchar(255);column:servings" json:"servings"`       // Number of servings
	CookTime       string `gorm:"type:varchar(255);column:cook_time" json:"cook_time"`     // Associated images
	Difficulty     string `gorm:"type:varchar(50);column:difficulty" json:"difficulty"`    // Difficulty level
	Needs          string `gorm:"type:text;column:needs" json:"needs"`                     // Ingredients needed (stored as JSON string)
	Process        string `gorm:"type:text;column:process" json:"process"`                 // Cooking process
	Likes          uint   `gorm:"type:int;default:0;column:likes" json:"likes"`            // Number of likes
	CommentAllowed bool   `gorm:"type:bool;column:comment_allowed" json:"comment_allowed"` // Whether comments are allowed
	Tags           string `gorm:"type:text;column:tags" json:"tags"`                       // Tags associated with the recipe

	//ThisFoodID   Food `gorm:"foreignKey:FoodID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`   // Food associated with the recipe
	ThisAuthorID User `gorm:"foreignKey:AuthorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // User who created the recipe
}

// 添加构造函数
func NewRecipe() *Recipe {
	return &Recipe{
		CommentAllowed: true, // 设置默认值
		Likes:          0,
	}
}
