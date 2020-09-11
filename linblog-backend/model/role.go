package model

type Role struct {
	ID    uint
	Name  string
	Users []User `gorm:"foreignkey:RoleID" json:"-"`
}
