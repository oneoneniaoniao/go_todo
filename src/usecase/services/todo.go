package services

import (
	"context"
	"github.com/oneoneniaoniao/go_todo/src/domain/models"
	"github.com/oneoneniaoniao/go_todo/src/domain/repositories"
)

// TodoService - Todoアプリケーションサービス
type TodoService struct {
	repo repositories.TodoRepository
}

// NewTodoService - 新しいTodoServiceを作成
func NewTodoService(repo repositories.TodoRepository) *TodoService {
	return &TodoService{
		repo: repo,
	}
}

// CreateTodo - 新しいTodoを作成
func (s *TodoService) CreateTodo(ctx context.Context, content string) error {
	todo := &models.Todo{Content: content}
	return s.repo.Create(ctx, todo)
}

// GetTodoByID - IDによるTodoの取得
func (s *TodoService) GetTodoByID(ctx context.Context, id uint) (*models.Todo, error) {
	return s.repo.GetByID(ctx, id)
}

// UpdateTodo - Todoの更新
func (s *TodoService) UpdateTodo(ctx context.Context, todo *models.Todo) error {
	if err := todo.Validate(); err != nil {
        return err
    }
	return s.repo.Update(ctx, todo)
}

// DeleteTodo - Todoの削除
func (s *TodoService) DeleteTodo(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// ListTodos - Todoリストの取得
func (s *TodoService) ListTodos(ctx context.Context) ([]*models.Todo, error) {
	return s.repo.List(ctx)
}