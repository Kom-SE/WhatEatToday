package models

import "gorm.io/gorm"

type CollectedRecipes struct {
	gorm.Model
	UserID   uint `gorm:"not null;index"`
	RecipeID uint `gorm:"not null;index"`

	ThisUser   User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ThisRecipe Recipe `gorm:"foreignKey:RecipeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// 构造函数
func NewCollectedRecipes() *CollectedRecipes {
	return &CollectedRecipes{}
}
