package post

import (
	"go-api-posts/internal/stat"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Image       string      `json:"image"`
	Description string      `json:"description"`
	UserID      uint        `json:"user_id" gorm:"index"`
	Stats       []stat.Stat `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func NewPost(image string, description string, userID uint) *Post {
	return &Post{
		Image:       image,
		Description: description,
		UserID:      userID,
	}
}
