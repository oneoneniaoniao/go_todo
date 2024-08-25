package models

import (
	"errors"
	"gorm.io/gorm"
)

type Todo struct {
	*gorm.Model
	Content string
}

// Validate - Todoのバリデーションルール
func (t *Todo) Validate() error {
    if t.Content == "" {
        return errors.New("content cannot be empty")
    }
    return nil
}