package repository

import (
	"github.com/go-redis/redis/v7"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
	rc *redis.Client
}

func NewRepository(db *gorm.DB, rc *redis.Client) *Repository {
	return &Repository{db, rc}
}
