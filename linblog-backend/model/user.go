package model

import "time"

type User struct {
	ID        string `gorm:"type:varchar(50);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	UserName  string `gorm:"type:varchar(255);not null"`
	Email     string `gorm:"type:varchar(255);not null"`
	NickName  string `gorm:"type:varchar(255);not null"`
	Password  string `gorm:"type:varchar(255);not null"`
	RoleID    uint   `gorm:"not null"`
	LastLogin *time.Time
	Posts     []Post `gorm:"foreignkey:UserID"`
}
