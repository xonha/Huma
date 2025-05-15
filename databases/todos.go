package databases

import (
	"context"
	"database/sql"
	"log"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/xonha/huma/models"
)

var Todos *bun.DB

func Init() {
	sqldb, err := sql.Open(sqliteshim.ShimName, "todos.db")
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	Todos = bun.NewDB(sqldb, sqlitedialect.New())

	// Create the table if it doesn't exist
	ctx := context.Background()
	_, err = Todos.NewCreateTable().
		Model((*models.Todo)(nil)). // <- this requires Todo model here
		IfNotExists().
		Exec(ctx)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}
