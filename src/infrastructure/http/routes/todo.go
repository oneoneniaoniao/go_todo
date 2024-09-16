package router

import (
	"github.com/gin-gonic/gin"
	"github.com/oneoneniaoniao/go_todo/src/interface/controllers"
)

// SetupRouter - ルータの設定を行います
func SetupRouterTodo(engine *gin.Engine, todoController *controllers.TodoController) *gin.Engine {
    // HTMLテンプレートのロード

	engine.POST("/todo/create", todoController.CreateTodo)
	engine.POST("/todo/update", todoController.UpdateTodo)
	engine.GET("/todo/delete", todoController.DeleteTodo)

	return engine
}