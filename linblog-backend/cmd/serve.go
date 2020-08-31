package cmd

import (
	"context"
	"github.com/Linjiangzhu/linblog/linblog-backend/controller"
	"github.com/Linjiangzhu/linblog/linblog-backend/repository"
	"github.com/Linjiangzhu/linblog/linblog-backend/service"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve blog-api backend server",
	RunE: func(cmd *cobra.Command, args []string) error {

		// get yaml config
		dbMap := viper.Get("DB")
		dbConfig := dbMap.(map[string]interface{})
		dbURL, _ := dbConfig["url"].(string)
		redisMap := viper.Get("REDIS")
		redisConfig := redisMap.(map[string]interface{})
		redisURL := redisConfig["url"].(string)

		// open mysql connection
		db, err := gorm.Open("mysql", dbURL)
		if err != nil {
			return err
		}

		// open redis connection
		rc := redis.NewClient(&redis.Options{Addr: redisURL})
		_, err = rc.Ping().Result()
		if err != nil {
			return err
		}

		defer db.Close()
		defer rc.Close()

		// establish app
		r := repository.NewRepository(db, rc)
		s := service.NewService(r)
		ctrl := controller.NewController(s)

		// establish api

		router := gin.Default()
		api := router.Group("/blog/api")
		{
			api.GET("/post/:pid", ctrl.GetPost)
			api.GET("/posts", ctrl.GetPosts)
		}
		authorized := api.Group("/admin")
		{
			authorized.POST("/post", ctrl.CreatePost)
			authorized.DELETE("/post/:pid", ctrl.DeletePost)
			authorized.PUT("/post/:pid", ctrl.UpdatePost)
			authorized.GET("/post/:pid", ctrl.GetPost)
			authorized.GET("/posts", ctrl.GetPosts)
			authorized.POST("/login", ctrl.Login)
			authorized.GET("/logout", ctrl.Logout)
		}
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
