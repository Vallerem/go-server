package models

import (
	"github.com/jinzhu/gorm"
)

type TodoModel struct {
	gorm.Model
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	UserID    uint
}

type TransformedTodo struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}
