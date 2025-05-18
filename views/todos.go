package views

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/xonha/todos/controllers"
)

func todos() {
	group := huma.NewGroup(api, "/todos")
	group.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"Todos"}
	})

	huma.Register(group, huma.Operation{
		OperationID:   "create-todo",
		Method:        http.MethodPost,
		Path:          "/",
		Summary:       "Create a new todo",
		DefaultStatus: http.StatusCreated,
	}, controllers.CreateTodo)

	huma.Register(group, huma.Operation{
		OperationID: "update-todo",
		Method:      http.MethodPut,
		Path:        "/{id}",
		Summary:     "Update a todo by ID",
	}, controllers.UpdateTodoById)

	huma.Register(group, huma.Operation{
		OperationID: "read-todos",
		Method:      http.MethodGet,
		Path:        "/",
		Summary:     "Read all todos",
	}, controllers.ReadTodos)

	huma.Register(group, huma.Operation{
		OperationID: "read-todo",
		Method:      http.MethodGet,
		Path:        "/{id}",
		Summary:     "Read a todo by ID",
	}, controllers.ReadTodoById)

	huma.Register(group, huma.Operation{
		OperationID: "delete-todo",
		Method:      http.MethodDelete,
		Path:        "/{id}",
		Summary:     "Delete a todo by ID",
	}, controllers.DeleteTodoById)
}
