package post

import (
	"go-api-posts/pkg/db"

	"gorm.io/gorm/clause"
)

type PostRepositoryDeps struct {
	Database *db.Db
}

type PostRepository struct {
	Database *db.Db
}

func NewPostRepository(db *db.Db) *PostRepository {
	return &PostRepository{
		Database: db,
	}
}

func (repo *PostRepository) Create(posts *Post) (*Post, error) {
	result := repo.Database.DB.Create(posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}

func (repo *PostRepository) GetById(id string) (*Post, error) {
	var post Post
	result := repo.Database.DB.First(&post, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &post, nil
}

func (repo *PostRepository) Update(posts *Post) (*Post, error) {
	result := repo.Database.DB.Clauses(clause.Returning{}).Updates(posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}

func (repo *PostRepository) Delete(post *Post) error {
	return repo.Database.DB.Delete(post).Error
}

func (repo *PostRepository) Count() int64 {
	var count int64
	repo.Database.
		Table("posts").
		Where("deleted_at is null").
		Count(&count)
	return count
}

func (repo *PostRepository) GetPosts(limit, offset int) []Post {
	var posts []Post
	repo.Database.
		Table("posts").
		Where("deleted_at is null").
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Scan(&posts)
	return posts
}
