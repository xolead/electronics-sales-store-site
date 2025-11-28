package main

import (
	"manage-service/pkg/handlers"
	"manage-service/pkg/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	public := r.Group("/")
	{
		public.POST("/registration", handlers.Registration)
		public.POST("/login", handlers.Login)
		public.GET("/product", handlers.GetAllProduct)
		public.GET("/product/:id", handlers.GetProduct)
	}

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/logout", handlers.Logout)
		protected.POST("/product")
		protected.DELETE("/product/:id")
		protected.PUT("/product/change")
	}

	r.Run(":8080")
}
