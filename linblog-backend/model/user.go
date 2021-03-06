package model

import "time"

type User struct {
	ID        string    `gorm:"type:varchar(255);" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	NickName  string    `json:"nickname"`
	RoleID    uint      `gorm:"not null" json:"role_id"`
	Posts     []Post    `gorm:"foreignkey:UserID" json:"-"`
}
