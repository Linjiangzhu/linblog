package main

import (
	"context"
	"github.com/Linjiangzhu/linblog/linblog-backend/controller"
	"github.com/Linjiangzhu/linblog/linblog-backend/middleware"
	"github.com/Linjiangzhu/linblog/linblog-backend/repository"
	"github.com/Linjiangzhu/linblog/linblog-backend/service"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	db, err := gorm.Open("mysql", "root:password@/blogdb_v2?charset=utf8&parseTime=True&loc=Local")
	db.LogMode(true)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rc := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rc.Close()
	r := repository.NewRepository(db, rc)
	s := service.NewService(r)
	ctrl := controller.NewController(s)

	router := gin.Default()

	router.GET("/posts", ctrl.GetPosts)
	router.GET("/post/:pid", ctrl.GetPost)
	router.POST("/admin/login", ctrl.Login)

	authorized := router.Group("/admin")
	authorized.Use(middleware.NewJWTAuth(db, rc))
	{
		authorized.GET("/posts", ctrl.GetPosts)
		authorized.GET("/post/:pid", ctrl.GetPost)
		authorized.POST("/post", ctrl.CreatePost)
		authorized.DELETE("/post/:pid", ctrl.DeletePost)
		authorized.PUT("/post/:pid", ctrl.UpdatePost)
		authorized.GET("/logout", ctrl.Logout)
	}

	srv := &http.Server{
		Addr:    ":9090",
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
}
