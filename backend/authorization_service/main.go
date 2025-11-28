package main

import (
	"authoriz-service/pkg/dbwork"
	"authoriz-service/pkg/handlers"
	"context"
	"log"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func main() {

	ctx := context.Background()
	db, err := dbwork.NewPostgreSQL(ctx, dbwork.LoadPSQLConfig())

	if err != nil {
		log.Fatalf("Ошибка подключения к бд: %v", err)
	}

	defer db.Close()

	if err = db.RunMigrations(ctx, ""); err != nil {
		log.Fatalf("Ошибка миграции бд: %v", err)
	}

	handler := handlers.NewHandler(db)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://manage_service:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/authorization", handler.Authorization)
	r.POST("/refresh", handler.Refresh)
	r.POST("/logout", handler.Logout)
	r.GET("/admin/:access", handler.Admin)
	r.GET("/uuid/:access", handler.GetUUID)

	r.Run(":8083")
}
