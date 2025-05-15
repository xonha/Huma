package models

import (
	"github.com/uptrace/bun"
)

type Todo struct {
	bun.BaseModel `bun:"table:todos,alias:t"`

	ID        int64  `bun:",pk,autoincrement" json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}
