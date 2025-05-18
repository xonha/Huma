package services

import (
	"context"
	"strconv"

	"github.com/xonha/TodosGo/databases"
	"github.com/xonha/TodosGo/models"
)

func CreateTodo(ctx context.Context, title string, completed bool) (*models.Todo, error) {
	todo := &models.Todo{
		Title:     title,
		Completed: completed,
	}
	_, err := databases.Todos.NewInsert().Model(todo).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func UpdateTodoById(ctx context.Context, idStr string, title string, completed bool) (*models.Todo, error) {
	id, _ := strconv.ParseInt(idStr, 10, 64)
	todo := &models.Todo{
		ID:        id,
		Title:     title,
		Completed: completed,
	}
	_, err := databases.Todos.NewUpdate().Model(todo).WherePK().Exec(ctx)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func GetTodos(ctx context.Context) ([]models.Todo, error) {
	var todos []models.Todo
	err := databases.Todos.NewSelect().Model(&todos).Order("id ASC").Scan(ctx)
	return todos, err
}

func GetTodoById(ctx context.Context, idStr string) (*models.Todo, error) {
	id, _ := strconv.ParseInt(idStr, 10, 64)
	todo := new(models.Todo)
	err := databases.Todos.NewSelect().Model(todo).Where("id = ?", id).Scan(ctx)
	return todo, err
}

func DeleteTodoById(ctx context.Context, idStr string) error {
	id, _ := strconv.ParseInt(idStr, 10, 64)
	_, err := databases.Todos.NewDelete().Model(&models.Todo{}).Where("id = ?", id).Exec(ctx)
	return err
}
