package controllers

import (
	"net/http"
	"time"
	"todo-api/config"
	"todo-api/models"

	"github.com/gin-gonic/gin"
)

type CreateTodoRequest struct {
	Title    string `json:"title" binding:"required"`
	Deadline string `json:"deadline"`
}

type UpdateTodoRequest struct {
	Title       *string `json:"title"`
	IsCompleted *bool   `json:"is_completed"`
	Deadline    *string `json:"deadline"`
}

func GetTodos(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var todos []models.Todo
	if err := config.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"todos": todos})
}

func GetTodo(c *gin.Context) {
	userID, _ := c.Get("user_id")
	todoID := c.Param("id")

	var todo models.Todo
	if err := config.DB.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"todo": todo})
}

func CreateTodo(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var deadline *time.Time
	if req.Deadline != "" {
		parsed, err := time.Parse("2006-01-02T15:04:05.000", req.Deadline)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deadline format"})
			return
		}
		deadline = &parsed
	}

	todo := models.Todo{
		Title:       req.Title,
		Date:        time.Now(),
		Deadline:    deadline,
		UserID:      userID.(uint),
		IsCompleted: false,
	}

	if err := config.DB.Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"todo": todo})
}

func ToggleTodo(c *gin.Context) {
	userID, _ := c.Get("user_id")
	todoID := c.Param("id")

	var todo models.Todo
	if err := config.DB.
		Where("id = ? AND user_id = ?", todoID, userID).
		First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	todo.IsCompleted = !todo.IsCompleted

	if err := config.DB.Save(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"todo": todo})
}

func UpdateTodo(c *gin.Context) {
	userID, _ := c.Get("user_id")
	todoID := c.Param("id")

	var todo models.Todo
	if err := config.DB.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	var req UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update title
	if req.Title != nil {
		todo.Title = *req.Title
	}

	// Update completed
	if req.IsCompleted != nil {
		todo.IsCompleted = *req.IsCompleted
	}

	// ✅ Update deadline (STRING → TIME)
	if req.Deadline != nil {
		// kalau string kosong → hapus deadline
		if *req.Deadline == "" {
			todo.Deadline = nil
		} else {
			parsed, err := time.Parse("2006-01-02T15:04:05.000", *req.Deadline)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deadline format"})
				return
			}
			todo.Deadline = &parsed
		}
	}

	if err := config.DB.Save(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"todo": todo})
}

func DeleteTodo(c *gin.Context) {
	userID, _ := c.Get("user_id")
	todoID := c.Param("id")

	var todo models.Todo
	if err := config.DB.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	if err := config.DB.Delete(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}
