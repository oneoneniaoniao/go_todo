package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oneoneniaoniao/go_todo/src/domain/models"
	"github.com/oneoneniaoniao/go_todo/src/infra/database/repositories"
	"github.com/oneoneniaoniao/go_todo/src/infra/database"
)


func errorDB(db *gorm.DB, c *gin.Context) bool {
	if db.Error != nil {
		log.Printf("Error todos: %v", db.Error)
		c.AbortWithStatus(http.StatusInternalServerError)
		return true // エラーがあったことを示す
	}
	return false // エラーがなかったことを示す
}

func listeners(r *gin.Engine, db *gorm.DB) {
	r.GET("/todo/delete", func(c *gin.Context) {
		id, _ := c.GetQuery("id")
		result := db.Delete(&models.Todo{}, id)
		if errorDB(result, c) { return }
		c.Redirect(http.StatusMovedPermanently, "/index")
	})

	r.POST("/todo/update", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.PostForm("id"))
		content := c.PostForm("content")
		var todo *models.Todo
		result := db.Where("id = ?", id).Take(&todo)
		if errorDB(result, c) { return }
		todo.Content = content
		result = db.Save(&todo)
		if errorDB(result, c) { return }
		c.Redirect(http.StatusMovedPermanently, "/index")
	})

	r.POST("/todo/create", func(c *gin.Context) {
		content := c.PostForm("content")
		fmt.Println(c.Request.PostForm, content)
		result := db.Create(&models.Todo{Content: content})
		if errorDB(result, c) { return }
		c.Redirect(http.StatusMovedPermanently, "/index")
	})

	r.GET("/todo/list", func(c *gin.Context) {
		var todos []*models.Todo
		// Get all records
		result := db.Find(&todos)
		if errorDB(result, c) { return }
		fmt.Println(json.NewEncoder(os.Stdout).Encode(todos))
		c.JSON(http.StatusOK, todos)
	})

	r.GET("/todo/get", func(c *gin.Context) {
		var todo *models.Todo
		id, _ := c.GetQuery("id")
		result := db.First(&todo, id)
		if errorDB(result, c) { return }
		// JSON形式でレスポンスを返す
		fmt.Println(json.NewEncoder(os.Stdout).Encode(todo))
		c.JSON(http.StatusOK, todo)
	})

	r.GET("/index", func(c *gin.Context) {
		var todos []*models.Todo
		result := db.Find(&todos)
		if errorDB(result, c) { return }
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "やることリスト",
			"todos": todos,
		})
	})

	//todo edit
	r.GET("/edit", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Query("id"))
		if err != nil {
			log.Fatalln(err)
		}
		var todo *models.Todo
		db.Where("id = ?", id).Take(&todo)
		c.HTML(http.StatusOK, "edit.html", gin.H{
			"title": "Todoの編集",
			"todo":  todo,
		})
	})
}

func main() {
	r := gin.Default()
	db, err := connectionDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.Todo{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	r.LoadHTMLGlob("src/infrastructure/http/public/*")
	listeners(r, db)

	fmt.Println("Database connection and setup successful")
	r.Run(":8080")
}