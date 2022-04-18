package services

import (
	"context"
	"net/http"

	"example/todo-go/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TodoServiceImpl struct {
	todoCollection *mongo.Collection
	ctx            context.Context
}

func NewTodoService(todoCollection *mongo.Collection, ctx context.Context) TodoService {
	return &TodoServiceImpl{
		todoCollection: todoCollection,
		ctx:            ctx,
	}
}

func (t *TodoServiceImpl) CreateTodo(todo *models.NewTodo) error {
	_, err := t.todoCollection.InsertOne(t.ctx, todo)
	return err
}

func (t *TodoServiceImpl) GetTodo(id string) (*models.Todo, error) {
	var todo *models.Todo
	objID, _ := primitive.ObjectIDFromHex(id)
	err := t.todoCollection.FindOne(t.ctx, bson.M{"_id": objID}).Decode(&todo)
	return todo, err
}

func (t *TodoServiceImpl) GetAll() ([]*models.Todo, error) {
	cursor, err := t.todoCollection.Find(t.ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var todos []*models.Todo
	for cursor.Next(t.ctx) {
		var todo models.Todo
		err := cursor.Decode(&todo)
		if err != nil {
			return nil, err
		}

		todos = append(todos, &todo)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(todos) == 0 {
		return []*models.Todo{}, nil
	}

	cursor.Close(t.ctx)

	return todos, nil
}

func (t *TodoServiceImpl) UpdateTodo(id string, todo *models.Todo) *models.HttpError {
	objID, _ := primitive.ObjectIDFromHex(id)
	result, err := t.todoCollection.UpdateOne(t.ctx, bson.M{"_id": objID}, bson.D{bson.E{Key: "$set", Value: todo}})
	if result != nil && result.MatchedCount < 1 {
		return &models.HttpError{
			StatusCode: http.StatusNotFound,
			Message:    "No matches found",
		}
	}

	if err != nil {
		return &models.HttpError{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	return nil
}

func (t *TodoServiceImpl) DeleteTodo(id string) *models.HttpError {
	objID, _ := primitive.ObjectIDFromHex(id)
	result, err := t.todoCollection.DeleteOne(t.ctx, bson.M{"_id": objID})

	if result.DeletedCount < 1 {
		return &models.HttpError{
			StatusCode: http.StatusNotFound,
			Message:    "No matches found",
		}
	}

	if err != nil {
		return &models.HttpError{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	return nil
}
