package model

type Tag struct {
	ID    uint    `gorm:"primary_key"`
	Name  string  `gorm:"varchar(255);not null"`
	Posts []*Post `gorm:"many2many:post_tag;"`
}
