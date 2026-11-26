package stat

import (
	"go-api-posts/pkg/db"
	"time"

	"gorm.io/datatypes"
)

type StatRepository struct {
	*db.Db
}

func NewStatRepository(db *db.Db) *StatRepository {
	return &StatRepository{db}
}

func (repo *StatRepository) AddClick(postId uint) {
	var stat Stat
	currentDate := datatypes.Date(time.Now())
	repo.Db.Find(&stat, "post_id = ? and date = ?", postId, currentDate)
	if stat.ID == 0 {
		repo.Db.Create(&Stat{
			PostId: postId,
			Clicks: 1,
			Date:   currentDate,
		})
	} else {
		stat.Clicks++
		repo.Db.Save(&stat)
	}

}

func (repo *StatRepository) GetStats(by string, from time.Time, to time.Time) []GetStatResponse {
	var stats []GetStatResponse
	var selectQuery string
	switch by {
	case FilterByDay:
		selectQuery = "to_char(date, 'YYYY-MM-DD') as period, sum(clicks)"
	case FilterByMonth:
		selectQuery = "to_char(date, 'YYYY-MM') as period, sum(clicks)"
	}
	repo.Db.Table("stats").
		Select(selectQuery).
		Where("date BETWEEN ? AND ?", from, to).
		Group("period").
		Scan(&stats)
	return stats
}
