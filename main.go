package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oneoneniaoniao/go_todo/src/domain/models"
	"github.com/oneoneniaoniao/go_todo/src/infrastructure/database"
	repository "github.com/oneoneniaoniao/go_todo/src/infrastructure/database/repositories"
)


func main() {
	engine := gin.Default()
	db, err := database.ConnectionDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// リポジトリの初期化
	todoRepo := repository.NewTodoRepository(db)


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
		todos, err := todoRepo.List(c.Request.Context())
		if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve todos"})
            return
        }
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Hello world",
			"todos": todos,
		})
	})

	//todo edit
	engine.GET("todo/edit", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Query("id"))
		if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "could not process invalid parameter "})
			return
		}
		var todo *models.Todo
		todo, err = todoRepo.GetByID(c.Request.Context(), uint(id))
if err != nil {
    log.Println("Error fetching todo:", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve todo"})
    return
}
		c.HTML(http.StatusOK, "edit.html", gin.H{
			"content": "Todo",
			"todo":  todo,
		})
	})

	engine.GET("/todo/delete", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Query("id"))
		if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "could not process invalid parameter "})
		}
		// uint64型をuintに変換して代入
		todoRepo.Delete(c.Request.Context(), uint(id))
		c.Redirect(http.StatusMovedPermanently, "/index")
	})

	engine.POST("/todo/update", func(c *gin.Context) {
		id, err := strconv.Atoi(c.PostForm("id"))
		if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "could not process invalid parameter "})
		}
		content := c.PostForm("content")
		var todo *models.Todo
		todo, _ = todoRepo.GetByID(c.Request.Context(), uint(id));
		todo.Content = content
		todoRepo.Update(c.Request.Context(), todo)
		c.Redirect(http.StatusMovedPermanently, "/index")
	})

	engine.POST("/todo/create", func(c *gin.Context) {
		content := c.PostForm("content")
		todoRepo.Create(c, &models.Todo{Content: content})
		c.Redirect(http.StatusMovedPermanently, "/index")
	})

	fmt.Println("Database connection and setup successful")
	engine.Run(":8080")
}