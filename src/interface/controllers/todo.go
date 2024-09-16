package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oneoneniaoniao/go_todo/src/usecase/services"
	"github.com/oneoneniaoniao/go_todo/src/domain/models"
)

// TodoController - Todoに関するHTTPリクエストを処理するコントローラ
type TodoController struct {
	service *services.TodoService
}

// NewTodoController - 新しいTodoControllerを作成します
func NewTodoController(service *services.TodoService) *TodoController {
	return &TodoController{
		service: service,
	}
}

// GetTodo - Todoを一つ取得します
func (tc *TodoController) GetTodoByID(c *gin.Context) (*models.Todo, error) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return nil, err
	}
	return tc.service.GetTodoByID(c.Request.Context(), uint(id))
}

// CreateTodo - Todoを作成します
func (tc *TodoController) CreateTodo(c *gin.Context) {
    // フォームデータからcontentを取得
    content := c.PostForm("content")
    if content == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "content is required"})
        return
    }
    err := tc.service.CreateTodo(c.Request.Context(), content)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create todo"})
        return
    }
    c.Status(http.StatusCreated)
	c.Redirect(http.StatusMovedPermanently, "/index")
}

// UpdateTodo - Todoを更新します
func (tc *TodoController) UpdateTodo(c *gin.Context) {
	idStr := c.PostForm("id")
    id, err := strconv.ParseUint(idStr, 10, 64) // idの安全な取得とエラーチェック
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }

	content := c.PostForm("content")
	todo, err := tc.service.GetTodoByID(c.Request.Context(), uint(id));
	if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Content is not found"})
	}
	todo.Content = content
	if err := tc.service.UpdateTodo(c.Request.Context(), todo); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
	c.Status(http.StatusOK)
	c.Redirect(http.StatusSeeOther, "/index")
}


// DeleteTodo - Todoを削除します
func (tc *TodoController) DeleteTodo(c *gin.Context) {
    // クエリからidを取得する
    idStr, ok := c.GetQuery("id")
    if !ok {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID query parameter is required"})
        return
    }
	// 文字列のidをuintに変換
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	err = tc.service.DeleteTodo(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
		return
	}
	c.Status(http.StatusOK)
	c.Redirect(http.StatusSeeOther, "/index")
}


// ListTodos - Todoのリストを取得します
func (tc *TodoController) ListTodos(c *gin.Context) ([]*models.Todo, error) {
    return tc.service.ListTodos(c.Request.Context())
}