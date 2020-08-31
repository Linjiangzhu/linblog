package model

type Role struct {
	ID    uint   `gorm:"" json:""`
	Name  string `gorm:"" json:""`
	Users []User `gorm:"foreignkey:RoleID" json:"-"`
}
