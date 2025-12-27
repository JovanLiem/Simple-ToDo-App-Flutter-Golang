package routes

import (
	"todo-api/controllers"
	"todo-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Health check
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Todo API is running"})
	})

	// Auth routes
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.GET("/profile", middleware.AuthMiddleware(), controllers.GetProfile)
	}

	// Todo routes (protected)
	todos := r.Group("/api/todos")
	todos.Use(middleware.AuthMiddleware())
	{
		todos.GET("", controllers.GetTodos)
		todos.GET("/:id", controllers.GetTodo)
		todos.POST("", controllers.CreateTodo)
		todos.PUT("/:id", controllers.UpdateTodo)
		todos.PATCH("/:id", controllers.ToggleTodo)
		todos.DELETE("/:id", controllers.DeleteTodo)
	}
}
