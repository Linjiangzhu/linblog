package main

import (
	"fmt"
	"github.com/Linjiangzhu/linblog/linblog-backend/cmd"
	"github.com/Linjiangzhu/linblog/linblog-backend/model"
	"github.com/rs/xid"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	var config cmd.Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}
	var db *gorm.DB
	fmt.Println("check if database exist...")
	if hasDatabase(&config) {
		fmt.Println("exist, start migrating schema...")
		db, _ := gorm.Open(mysql.New(mysql.Config{
			DSN: fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
				config.Mysql.Username,
				config.Mysql.Password,
				config.Mysql.Address,
				config.Mysql.DBName),
		}), &gorm.Config{})
		migrateTable(db)
		fmt.Println("finished!")
	} else {
		fmt.Println("not exist, start creating database...")
		db, _ = createDatabase(&config)
		fmt.Println("database created, start creating schema...")
		migrateTable(db)
		fmt.Println("tables created, start insert dummy data...")
		fillDummyData(db)
		fmt.Println("finished!")
	}

}

func hasDatabase(config *cmd.Config) bool {
	dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.Mysql.Username,
		config.Mysql.Password,
		config.Mysql.Address,
		config.Mysql.DBName)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{})
	sqlDB, _ := db.DB()
	sqlDB.Close()
	return err == nil
}

func createDatabase(config *cmd.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.Mysql.Username,
		config.Mysql.Password,
		config.Mysql.Address,
		config.Mysql.DBName)
	// connect to dsn without database
	newDSN := fmt.Sprintf("%s:%s@(%s)/",
		config.Mysql.Username,
		config.Mysql.Password,
		config.Mysql.Address)
	newDB, _ := gorm.Open(mysql.New(mysql.Config{
		DSN: newDSN,
	}), &gorm.Config{})
	newDB.Exec("CREATE DATABASE " + config.Mysql.DBName)
	newSQLDB, _ := newDB.DB()
	newSQLDB.Close()
	return gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})

}

func migrateTable(db *gorm.DB) {
	_ = db.AutoMigrate(
		&model.Role{},
		&model.User{},
		&model.Post{},
		&model.Tag{},
		&model.Category{})
}

func fillDummyData(db *gorm.DB) {
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

	root := "./resource/"
	files, err := ioutil.ReadDir(root)
	if err != nil {
		panic(err)
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
}
