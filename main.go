package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oneoneniaoniao/go_todo/src/domain/models"
	"github.com/oneoneniaoniao/go_todo/src/infrastructure/database"
	"github.com/oneoneniaoniao/go_todo/src/infrastructure/database/repositories"
	"github.com/oneoneniaoniao/go_todo/src/usecase/services"
	"github.com/oneoneniaoniao/go_todo/src/interface/controllers"

)


func main() {
	engine := gin.Default()
	db, err := database.ConnectionDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// リポジトリの初期化
	todoRepo := repository.NewTodoRepository(db)
	todoService := services.NewTodoService(todoRepo)
	todoController := controllers.NewTodoController(todoService)

	// Migrate the schema
	err = db.AutoMigrate(&models.Todo{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	engine.Static("/static", "./static");
	engine.LoadHTMLGlob("src/infrastructure/http/public/*")
	// 下記は開発環境のhtmlのホットリロード用
	engine.Use(func(c *gin.Context) {
    engine.LoadHTMLGlob("src/infrastructure/http/public/*")
    c.Next()
})
	engine.GET("/index", func(c *gin.Context) {
		var todos []*models.Todo

		// Get all records
		todos, err := todoController.ListTodos(c)
		if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve todos"})
				return
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
				"title": "やることリスト",
				"todos": todos,
		})
})

	//todo edit
	engine.GET("todo/edit", func(c *gin.Context) {
		todo, err := todoController.GetTodoByID(c)
		if err != nil {
			return
		}
		c.HTML(http.StatusOK, "edit.html", gin.H{
			"content": "Todoの編集",
			"todo":  todo,
		})
	})

	engine.GET("/todo/delete", func(c *gin.Context) {
		todoController.DeleteTodo(c)
		c.Redirect(http.StatusMovedPermanently, "/index")
	})

	engine.POST("/todo/update", func(c *gin.Context) {
		todoController.UpdateTodo(c)
		c.Redirect(http.StatusMovedPermanently, "/index")
	})

	engine.POST("/todo/create", func(c *gin.Context) {
		todoController.CreateTodo(c)
		c.Redirect(http.StatusMovedPermanently, "/index")
	})

	fmt.Println("Database connection and setup successful")
	engine.Run(":8080")
}