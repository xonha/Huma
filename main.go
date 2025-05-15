package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humabunrouter"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bunrouter"
)

type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"3000"`
}

// Todo model
type Todo struct {
	bun.BaseModel `bun:"table:todos,alias:t"`

	ID        int64  `bun:",pk,autoincrement" json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type TodoInput struct {
	Body struct {
		Title     string `json:"title" maxLength:"100"`
		Completed bool   `json:"completed"`
	}
}

type TodoOutput struct {
	Body *Todo `json:"todo"`
}

type TodoListOutput struct {
	Body []Todo `json:"todos"`
}

type UpdateTodoInput struct {
	ID   string `path:"id"`
	Body struct {
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
}

var db *bun.DB

func main() {
	ctx := context.Background()

	sqldb, err := sql.Open(sqliteshim.ShimName, "todos.db")
	if err != nil {
		panic(err)
	}
	defer sqldb.Close()

	db = bun.NewDB(sqldb, sqlitedialect.New())

	_, err = db.NewCreateTable().
		Model((*Todo)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		panic(err)
	}

	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		router := bunrouter.New()
		config := huma.DefaultConfig("My API", "1.0.0")

		config.DocsPath = ""
		router.GET("/docs", bunrouter.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`<!doctype html>
			<html>
				<head>
					<title>API Reference</title>
					<meta charset="utf-8" />
					<meta
						name="viewport"
						content="width=device-width, initial-scale=1" />
				</head>
				<body>
					<script
						id="api-reference"
						data-url="/openapi.json"></script>
					<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
				</body>
			</html>`))
		}))
		config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
			"Bearer": {
				Type: "oauth2",
				Flows: &huma.OAuthFlows{
					AuthorizationCode: &huma.OAuthFlow{
						AuthorizationURL: "https://example.com/oauth/authorize",
						TokenURL:         "https://example.com/oauth/token",
						Scopes: map[string]string{
							"scope1": "Scope 1 description...",
							"scope2": "Scope 2 description...",
						},
					},
				},
			},
		}

		api := humabunrouter.New(router, config)
		grp := huma.NewGroup(api, "/todos")
		grp.UseSimpleModifier(func(op *huma.Operation) {
			op.Tags = []string{"Todos"}
			op.Security = []map[string][]string{
				{"myAuth": {"scope1"}},
			}
		})

		huma.Register(grp, huma.Operation{
			OperationID:   "create-todo",
			Method:        http.MethodPost,
			Path:          "/",
			Summary:       "Create a new todo",
			DefaultStatus: http.StatusCreated,
		}, func(ctx context.Context, input *TodoInput) (*TodoOutput, error) {
			todo := &Todo{
				Title:     input.Body.Title,
				Completed: input.Body.Completed,
			}
			_, err := db.NewInsert().Model(todo).Exec(ctx)
			if err != nil {
				return nil, err
			}
			return &TodoOutput{Body: todo}, nil
		})

		huma.Register(grp, huma.Operation{
			OperationID: "list-todos",
			Method:      http.MethodGet,
			Path:        "/",
			Summary:     "Get all todos",
		}, func(ctx context.Context, _ *struct{}) (*TodoListOutput, error) {
			var todos []Todo
			err := db.NewSelect().Model(&todos).Order("id ASC").Scan(ctx)
			if err != nil {
				return nil, err
			}
			return &TodoListOutput{Body: todos}, nil
		})

		huma.Register(grp, huma.Operation{
			OperationID: "get-todo",
			Method:      http.MethodGet,
			Path:        "/{id}",
			Summary:     "Get a todo by ID",
		}, func(ctx context.Context, input *struct {
			ID string `path:"id"`
		},
		) (*TodoOutput, error) {
			id, _ := strconv.ParseInt(input.ID, 10, 64)
			todo := new(Todo)
			err := db.NewSelect().Model(todo).Where("id = ?", id).Scan(ctx)
			if err != nil {
				return nil, err
			}
			return &TodoOutput{Body: todo}, nil
		})

		huma.Register(grp, huma.Operation{
			OperationID: "update-todo",
			Method:      http.MethodPut,
			Path:        "/{id}",
			Summary:     "Update a todo by ID",
		}, func(ctx context.Context, input *UpdateTodoInput) (*TodoOutput, error) {
			id, _ := strconv.ParseInt(input.ID, 10, 64)
			todo := &Todo{
				ID:        id,
				Title:     input.Body.Title,
				Completed: input.Body.Completed,
			}
			_, err := db.NewUpdate().Model(todo).WherePK().Exec(ctx)
			if err != nil {
				return nil, err
			}
			return &TodoOutput{Body: todo}, nil
		})

		huma.Register(grp, huma.Operation{
			OperationID: "delete-todo",
			Method:      http.MethodDelete,
			Path:        "/{id}",
			Summary:     "Delete a todo by ID",
		}, func(ctx context.Context, input *struct {
			ID string `path:"id"`
		},
		) (*struct{}, error) {
			id, _ := strconv.ParseInt(input.ID, 10, 64)
			_, err := db.NewDelete().Model(&Todo{}).Where("id = ?", id).Exec(ctx)
			return &struct{}{}, err
		})

		hooks.OnStart(func() {
			fmt.Printf("🚀 Server listening on port %d...\n", options.Port)
			err := http.ListenAndServe(fmt.Sprintf(":%d", options.Port), router)
			if err != nil {
				panic(err)
			}
		})
	})

	cli.Run()
}
