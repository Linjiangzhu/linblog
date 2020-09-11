package cmd

import (
	"fmt"
	"github.com/Linjiangzhu/linblog/linblog-backend/model"
	"github.com/rs/xid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize mysql connection and database schema",
	RunE: func(cmd *cobra.Command, args []string) error {
		var config Config
		if err := viper.Unmarshal(&config); err != nil {
			return err
		}
		dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			config.Mysql.Username,
			config.Mysql.Password,
			config.Mysql.Address,
			config.Mysql.DBName)
		log.Println("establish database connection...")
		db, err := gorm.Open(mysql.New(mysql.Config{
			DSN: dsn,
		}), &gorm.Config{})
		if err != nil {
			return err
		}
		log.Println("database connected, start migrating schema...")
		sqlDB, _ := db.DB()
		defer sqlDB.Close()
		isEmpty := !db.Migrator().HasTable("users")
		err = migrateSchema(db)
		if err != nil {
			return err
		}
		if isEmpty {
			log.Println("cannot find user table, start inserting initial data...")
			err = initData(db)
			if err != nil {
				return err
			}
		}
		log.Println("finished!")
		return nil
	},
}

func initData(db *gorm.DB) error {
	role := model.Role{Name: "admin"}
	db.Create(&role)
	username := os.Getenv("BLOG_ADMIN_USERNAME")
	password := os.Getenv("BLOG_ADMIN_PASSWORD")
	if len(username) == 0 {
		username = "admin"
	}
	if len(password) == 0 {
		password = "password"
	}
	uid := xid.New().String()
	pw, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	admin := model.User{
		ID:       uid,
		Username: username,
		Password: string(pw),
		NickName: "default admin",
		RoleID:   role.ID,
	}
	db.Create(&admin)

	root := "./resource"
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}
	for _, f := range files {
		relFilePath := filepath.Join(root, f.Name())
		byteStr, err := ioutil.ReadFile(relFilePath)
		if err != nil {
			continue
		}
		p := model.Post{
			Title:   f.Name(),
			Brief:   "",
			Content: string(byteStr),
			Visible: true,
			UserID:  admin.ID,
		}
		db.Create(&p)
	}
	return nil
}

func migrateSchema(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Role{},
		&model.User{},
		&model.Post{},
		&model.Tag{},
		&model.Category{})
}
