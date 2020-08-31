package main

import (
	"fmt"
	"github.com/Linjiangzhu/blog-v2/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func initDB() *gorm.DB {
	db, err := gorm.Open("mysql", "root:password@/blogdb_v2?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return db
}

func main() {
	db := initDB()
	defer db.Close()
	db.Exec("DROP TABLE IF EXISTS post_tag")
	db.Exec("DROP TABLE IF EXISTS post_cat")
	db.Exec("DROP TABLE IF EXISTS tags")
	db.Exec("DROP TABLE IF EXISTS categories")
	db.Exec("DROP TABLE IF EXISTS posts")
	db.Exec("DROP TABLE IF EXISTS users")
	db.Exec("DROP TABLE IF EXISTS roles")

	// create roles
	db.AutoMigrate(&model.Role{})
	roleAdmin := model.Role{Name: "admin"}
	roleUser := model.Role{Name: "user"}
	db.Create(&roleAdmin)
	db.Create(&roleUser)
	// create users
	db.AutoMigrate(&model.User{})
	db.Model(&model.User{}).AddForeignKey("role_id", "roles(id)", "RESTRICT", "RESTRICT")

	pw1, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	userAdmin := model.User{
		ID:       "0001",
		Username: "admin01",
		NickName: "admin 01",
		Password: string(pw1),
		RoleID:   roleAdmin.ID,
	}
	pw2, _ := bcrypt.GenerateFromPassword([]byte("user"), bcrypt.DefaultCost)
	userUser := model.User{
		ID:       "0002",
		Username: "user01",
		NickName: "user 01",
		Password: string(pw2),
		RoleID:   roleUser.ID,
	}
	db.Create(&userAdmin)
	db.Create(&userUser)

	// create posts
	db.AutoMigrate(&model.Post{})
	db.Model(&model.Post{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	p1 := model.Post{
		Title:   "p01",
		Brief:   "p01",
		Content: "p01",
		Visible: true,
		UserID:  userAdmin.ID,
	}
	p2 := model.Post{
		Title:   "p02",
		Brief:   "p02",
		Content: "p02",
		Visible: true,
		UserID:  userAdmin.ID,
	}
	p3 := model.Post{
		Title:   "p03",
		Brief:   "p03",
		Content: "p03",
		Visible: true,
		UserID:  userAdmin.ID,
	}
	p4 := model.Post{
		Title:   "p04",
		Brief:   "p04",
		Content: "p04",
		Visible: true,
		UserID:  userUser.ID,
	}
	p5 := model.Post{
		Title:   "p05",
		Brief:   "p05",
		Content: "p05",
		Visible: true,
		UserID:  userUser.ID,
	}
	p6 := model.Post{
		Title:   "p06",
		Brief:   "p06",
		Content: "p06",
		Visible: true,
		UserID:  userUser.ID,
	}
	db.Create(&p1)
	db.Create(&p2)
	db.Create(&p3)
	db.Create(&p4)
	db.Create(&p5)
	db.Create(&p6)

	db.AutoMigrate(&model.Tag{})
	db.Table("post_tag").AddForeignKey("post_id", "posts(id)", "RESTRICT", "RESTRICT")
	db.Table("post_tag").AddForeignKey("tag_id", "tags(id)", "RESTRICT", "RESTRICT")
	tag1 := model.Tag{Name: "tag01"}
	tag2 := model.Tag{Name: "tag02"}
	db.Model(&p1).Association("Tags").Append(&tag1)
	db.Model(&p2).Association("Tags").Append(&tag1, &tag2)
	db.Model(&p4).Association("Tags").Append(&tag1)
	db.Model(&p5).Association("Tags").Append(&tag1, &tag2)

	db.AutoMigrate(&model.Category{})
	cat1 := model.Category{Name: "cat01"}
	cat2 := model.Category{Name: "cat02"}
	db.Model(&p1).Association("Categories").Append(&cat1)
	db.Model(&p2).Association("Categories").Append(&cat2)
	db.Model(&p4).Association("Categories").Append(&cat1)
	db.Model(&p5).Association("Categories").Append(&cat2)
	db.Table("post_cat").AddForeignKey("post_id", "posts(id)", "RESTRICT", "RESTRICT")
	db.Table("post_cat").AddForeignKey("category_id", "categories(id)", "RESTRICT", "RESTRICT")

}
