package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/oneoneniaoniao/go_todo/src/domain/models"
	"github.com/oneoneniaoniao/go_todo/src/infrastructure/database"
	"github.com/oneoneniaoniao/go_todo/src/infrastructure/database/repositories"
	"github.com/oneoneniaoniao/go_todo/src/usecase/services"
	"github.com/oneoneniaoniao/go_todo/src/interface/controllers"
	"github.com/oneoneniaoniao/go_todo/src/infrastructure/http/routes"
)


func main() {
	// データベース接続の設定
	db, err := database.ConnectionDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// リポジトリの初期化
	todoRepo := repository.NewTodoRepository(db)

	// サービス層の初期化
	todoService := services.NewTodoService(todoRepo)

	// コントローラの初期化
	todoController := controllers.NewTodoController(todoService)

	engine := gin.Default()
	// ルータの設定
	engine = router.SetupRouterTodo(engine, todoController)
	engine = router.SetupRouterPage(engine, todoController)


	// Migrate the schema
	err = db.AutoMigrate(&models.Todo{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	engine.Static("/static", "./static");

	// 下記は開発環境のhtmlのホットリロード用
// 	engine.Use(func(c *gin.Context) {
//     engine.LoadHTMLGlob("src/infrastructure/http/public/*")
//     c.Next()
// })
	
	// サーバを8080ポートで起動
	if err := engine.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}