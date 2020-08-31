package model

type Category struct {
	ID    uint    `gorm:"primary_key" json:"id"`
	Name  string  `gorm:"varchar(255);not null" json:"name"`
	Posts []*Post `gorm:"many2many:post_cat;" json:"-"`
}
