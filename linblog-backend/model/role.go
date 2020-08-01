package model

type Role struct {
	ID    uint   `gorm:"primary_key"`
	Name  string `gorm:"type:varchar(25);not null"`
	Users []User `gorm:"foreignkey:RoleID"`
}
