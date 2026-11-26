package stat

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Stat struct {
	gorm.Model
	PostId uint           `json:"post_id"`
	Clicks uint           `json:"clicks"`
	Date   datatypes.Date `json:"date"`
}
