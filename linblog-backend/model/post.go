package model

import "time"

type Post struct {
	ID         uint `gorm:"primary_key"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
	Title      string `gorm:"type:varchar(255);not null"`
	Brief      string `gorm:"type:varchar(500);not null"`
	Content    string `gorm:"type:text;not null"`
	Visible    bool   `gorm:"type:tinyint(1);not nul"`
	UserID     string
	Tags       []*Tag      `gorm:"many2many:post_tag;"`
	Categories []*Category `gorm:"many2many:post_cat;"`
}
