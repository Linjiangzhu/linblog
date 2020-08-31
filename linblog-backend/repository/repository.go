package repository

import (
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
)

type Repository struct {
	db *gorm.DB
	rc *redis.Client
}

func NewRepository(db *gorm.DB, rc *redis.Client) *Repository {
	return &Repository{db, rc}
}
