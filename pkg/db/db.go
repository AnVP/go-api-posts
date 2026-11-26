package db

import (
	"go-api-posts/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct {
	*gorm.DB
}

func NewDb(conf *configs.Config) *Db {
	db, err := gorm.Open(postgres.Open(conf.Db.Host), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &Db{db}
}
