package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	CommentText   string `gorm:"type:text;not null;column:comment_text" json:"comment_text"` // 评论内容
	CommentUserID uint   `gorm:"type:int;not null;column:comment_user" json:"comment_user"`  // 评论用户ID
	RecipeID      uint   `gorm:"type:int;not null;column:recipe_id" json:"recipe_id"`        // 评论对应的食谱ID
	ParentUserID  uint   `gorm:"type:int;column:parent_id" json:"parent_user_id"`            // 父评论ID，null表示没有父评论
	Likes         uint   `gorm:"type:int;default:0;column:likes" json:"likes"`               // 点赞数

	ThisCommentUserID User   `gorm:"foreignKey:CommentUserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // 评论用户
	ThisRecipeID      Recipe `gorm:"foreignKey:RecipeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`      // 评论对应的食谱
	ThisParentUserID  User   `gorm:"foreignKey:ParentUserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`  // 父评论用户
}
