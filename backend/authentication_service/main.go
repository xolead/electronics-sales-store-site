package main

import (
	"auth-service/pkg/dbwork"
	"auth-service/pkg/handlers"
	"context"
	"log"

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
	r.POST("/registration", handler.Registration)
	r.POST("/login", handler.Login)

	r.Run(":8081")

}
