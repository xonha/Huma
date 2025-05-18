package controllers

import (
	"context"

	"github.com/xonha/todos/schemas"
	"github.com/xonha/todos/services"
)

func CreateTodo(ctx context.Context, input *schemas.TodoInput) (*schemas.TodoOutput, error) {
	todo, err := services.CreateTodo(ctx, input.Body.Title, input.Body.Completed)
	if err != nil {
		return nil, err
	}
	return &schemas.TodoOutput{Body: todo}, nil
}

func UpdateTodoById(ctx context.Context, input *schemas.UpdateTodoInput) (*schemas.TodoOutput, error) {
	todo, err := services.UpdateTodoById(ctx, input.ID, input.Body.Title, input.Body.Completed)
	if err != nil {
		return nil, err
	}
	return &schemas.TodoOutput{Body: todo}, nil
}

func ReadTodos(ctx context.Context, _ *struct{}) (*schemas.TodoListOutput, error) {
	todos, err := services.GetTodos(ctx)
	if err != nil {
		return nil, err
	}
	return &schemas.TodoListOutput{Body: todos}, nil
}

func ReadTodoById(ctx context.Context, input *struct {
	ID string `path:"id"`
},
) (*schemas.TodoOutput, error) {
	todo, err := services.GetTodoById(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	return &schemas.TodoOutput{Body: todo}, nil
}

func DeleteTodoById(ctx context.Context, input *struct {
	ID string `path:"id"`
},
) (*struct{}, error) {
	err := services.DeleteTodoById(ctx, input.ID)
	return &struct{}{}, err
}
