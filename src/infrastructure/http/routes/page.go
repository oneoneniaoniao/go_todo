package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/oneoneniaoniao/go_todo/src/interface/controllers"
)

// SetupRouter - ルータの設定を行います
func SetupRouterPage(engine *gin.Engine, todoController *controllers.TodoController) *gin.Engine {

    // HTMLテンプレートのロード
	engine.LoadHTMLGlob("src/infrastructure/http/public/*")

	// 各ルートの設定
	engine.GET("/index", func(c *gin.Context) {
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

	engine.GET("/todo/edit", func(c *gin.Context) {
		todo, err := todoController.GetTodoByID(c)
		if err != nil {
			return
		}
		c.HTML(http.StatusOK, "edit.html", gin.H{
			"title": "Todoの編集",
			"todo":  todo,
		})
	})

	return engine
}