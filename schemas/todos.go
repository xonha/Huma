package schemas

import "github.com/xonha/huma/models"

type TodoInput struct {
	Body struct {
		Title     string `json:"title" maxLength:"100"`
		Completed bool   `json:"completed"`
	}
}

type TodoOutput struct {
	Body *models.Todo `json:"todo"`
}

type TodoListOutput struct {
	Body []models.Todo `json:"todos"`
}

type UpdateTodoInput struct {
	ID   string `path:"id"`
	Body struct {
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
}
