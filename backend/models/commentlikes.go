package models

import "gorm.io/gorm"

type CommentLike struct {
	gorm.Model
	UserID    uint `gorm:"not null;index:idx_user_comment,unique" json:"user_id"`    // 点赞用户ID
	CommentID uint `gorm:"not null;index:idx_user_comment,unique" json:"comment_id"` // 被点赞的评论ID

	// 联合唯一索引，确保一个用户对同一评论只能点赞一次
	User    User    `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Comment Comment `gorm:"foreignKey:CommentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
