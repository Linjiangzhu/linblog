package cmd

import (
	"context"
	"fmt"
	"github.com/Linjiangzhu/linblog/linblog-backend/controller"
	"github.com/Linjiangzhu/linblog/linblog-backend/middleware"
	"github.com/Linjiangzhu/linblog/linblog-backend/repository"
	"github.com/Linjiangzhu/linblog/linblog-backend/service"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serveApi(db *gorm.DB, rc *redis.Client, ctrl *controller.Controller) *gin.Engine {
	router := gin.Default()
	api := router.Group("/blog/api")
	{
		api.GET("/post/:pid", ctrl.GetPost)
		api.GET("/posts", ctrl.GetPosts)
		api.POST("/admin/login", ctrl.Login)
	}
	authorized := api.Group("/admin").Use(middleware.NewJWTAuth(db, rc))
	{
		authorized.POST("/post", ctrl.CreatePost)
		authorized.DELETE("/post/:pid", ctrl.DeletePost)
		authorized.PUT("/post/:pid", ctrl.UpdatePost)
		authorized.GET("/post/:pid", ctrl.GetPost)
		authorized.GET("/posts", ctrl.GetPosts)
		authorized.GET("/logout", ctrl.Logout)
	}
	return router
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve blog-api backend server",
	RunE: func(cmd *cobra.Command, args []string) error {
		var config Config
		if err := viper.Unmarshal(&config); err != nil {
			panic(err)
		}
		dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			config.Mysql.Username,
			config.Mysql.Password,
			config.Mysql.Address,
			config.Mysql.DBName)

		// open mysql connection
		db, err := gorm.Open(mysql.New(mysql.Config{DSN: dsn}), &gorm.Config{})
		if err != nil {
			return err
		}
		sqlDB, _ := db.DB()

		// open redis connection
		rc := redis.NewClient(&redis.Options{
			Addr:     config.Redis.Address,
			Password: config.Redis.Password,
			DB:       config.Redis.DB})
		_, err = rc.Ping().Result()
		if err != nil {
			return err
		}

		defer sqlDB.Close()
		defer rc.Close()

		// establish app
		r := repository.NewRepository(db, rc)
		s := service.NewService(r)
		ctrl := controller.NewController(s)

		// establish api
		router := serveApi(db, rc, ctrl)

		portStr := ":" + viper.GetString("PORT")
		srv := &http.Server{
			Addr:    portStr,
			Handler: router,
		}
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server forced to shutdown:", err)
		}
		log.Println("Server exiting")
		return nil
	},
}
