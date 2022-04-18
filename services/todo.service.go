package services

import "example/todo-go/models"

type TodoService interface {
	CreateTodo(*models.NewTodo) error
	GetTodo(string) (*models.Todo, error)
	GetAll() ([]*models.Todo, error)
	UpdateTodo(string, *models.Todo) *models.HttpError
	DeleteTodo(string) *models.HttpError
}
