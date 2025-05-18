package databases

import (
	"context"
	"database/sql"
	"log"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/xonha/TodosGo/models"
)

var Todos *bun.DB

func Init() {
	sqldb, err := sql.Open(sqliteshim.ShimName, ":memory:")
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	Todos = bun.NewDB(sqldb, sqlitedialect.New())

	ctx := context.Background()
	_, err = Todos.NewCreateTable().
		Model((*models.Todo)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}
