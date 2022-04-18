package controllers

import (
	"net/http"

	"example/todo-go/models"
	"example/todo-go/services"

	"github.com/gin-gonic/gin"
)

type TodoController struct {
	TodoService services.TodoService
}

func New(todoService services.TodoService) TodoController {
	return TodoController{
		TodoService: todoService,
	}
}

func (tc *TodoController) CreateTodo(ctx *gin.Context) {
	var todo models.NewTodo
	if err := ctx.ShouldBindJSON(&todo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := tc.TodoService.CreateTodo(&todo); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "success"})
}

func (tc *TodoController) GetTodo(ctx *gin.Context) {
	id := ctx.Param("id")
	todo, err := tc.TodoService.GetTodo(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, todo)
}

func (tc *TodoController) GetAll(ctx *gin.Context) {
	todos, err := tc.TodoService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, todos)
}

func (tc *TodoController) UpdateTodo(ctx *gin.Context) {
	var todo models.Todo
	if err := ctx.ShouldBindJSON(&todo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	id := ctx.Param("id")
	if err := tc.TodoService.UpdateTodo(id, &todo); err != nil {
		ctx.JSON(int(err.StatusCode), gin.H{"message": err.Message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (tc *TodoController) DeleteTodo(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := tc.TodoService.DeleteTodo(id); err != nil {
		ctx.JSON(int(err.StatusCode), gin.H{"message": err.Message})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"message": "success"})
}

func (tc *TodoController) RegisterTodoRoutes(rg *gin.RouterGroup) {
	todoRoute := rg.Group("/todos")
	todoRoute.POST("/", tc.CreateTodo)
	todoRoute.GET("/:id", tc.GetTodo)
	todoRoute.GET("/", tc.GetAll)
	todoRoute.PATCH("/:id", tc.UpdateTodo)
	todoRoute.DELETE("/:id", tc.DeleteTodo)
}
