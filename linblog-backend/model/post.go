package model

import (
	"time"
)

type Post struct {
	ID         uint        `gorm:"primary_key" json:"id"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	DeletedAt  *time.Time  `json:"-"`
	Title      string      `gorm:"type:varchar(255);not null" json:"title"`
	Brief      string      `gorm:"type:varchar(500);not null" json:"brief,omitempty"`
	Content    string      `gorm:"type:text;not null" json:"content,omitempty"`
	Visible    bool        `gorm:"type:tinyint(1);not nul" json:"visible"`
	UserID     string      `json:"uid"`
	Tags       []*Tag      `gorm:"many2many:post_tag;" json:"tags,omitempty"`
	Categories []*Category `gorm:"many2many:post_cat;" json:"cats,omitempty"`
}

func (p *Post) ViewClass(role string) interface{} {
	switch role {
	case "admin":
		return p.AdminView()
	case "visitor":
		return p.VisitorView()
	default:
		return struct{}{}
	}
}

func (p *Post) AdminView() interface{} {
	return Post{
		ID:         p.ID,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
		Title:      p.Title,
		Brief:      p.Brief,
		Content:    p.Content,
		Visible:    p.Visible,
		UserID:     p.UserID,
		Tags:       p.Tags,
		Categories: p.Categories,
	}
}

func (p *Post) VisitorView() interface{} {
	return struct {
		ID         uint        `json:"id"`
		CreatedAt  time.Time   `json:"created_at"`
		UpdatedAt  time.Time   `json:"updated_at"`
		Title      string      `json:"title"`
		Brief      string      `json:"brief,omitempty"`
		Content    string      `json:"content,omitempty"`
		Tags       []*Tag      `json:"tags,omitempty"`
		Categories []*Category `json:"cats,omitempty"`
	}{
		p.ID,
		p.CreatedAt,
		p.UpdatedAt,
		p.Title,
		p.Brief,
		p.Content,
		p.Tags,
		p.Categories,
	}
}
