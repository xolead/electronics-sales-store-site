package main

import (
	"manage-service/pkg/handlers"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	public := r.Group("/")
	{
		public.POST("/registration", handlers.Registration)
		public.POST("/login", handlers.Login)
		public.GET("/product")
		public.GET("/product/:id")
	}

	protected := r.Group("/")
	{
		protected.POST("/logout")
		protected.POST("/product")
		protected.DELETE("/product/:id")
		protected.PUT("/product/change")
	}

	r.Run(":8080")
}
