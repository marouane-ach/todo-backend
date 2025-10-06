package db

import (
	"context"
	"database/sql"

	"github.com/marouane-ach/todo-go/models"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var dbInstance *bun.DB
var dbContext context.Context

func GetDBIntance() *bun.DB {
	if dbInstance == nil {
		sqldb := sql.OpenDB(pgdriver.NewConnector(
			pgdriver.WithDSN("postgres://todoapp:todoapp@localhost:5432/todoapp?sslmode=disable"),
		))
		dbInstance = bun.NewDB(sqldb, pgdialect.New())
	}

	return dbInstance
}

func GetDBContext() context.Context {
	if dbContext == nil {
		dbContext = context.Background()
	}

	return dbContext
}

func CreateDBTables() {
	ctx := GetDBContext()
	db := GetDBIntance()

	_, err := db.NewCreateTable().Model((*models.User)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		panic(err)
	}

	_, err = db.NewCreateTable().Model((*models.Token)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		panic(err)
	}

	_, err = db.NewCreateTable().Model((*models.Color)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		panic(err)
	}

	_, err = db.NewCreateTable().Model((*models.TodoList)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		panic(err)
	}

	_, err = db.NewCreateTable().Model((*models.Todo)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		panic(err)
	}
}

func SeedColorsTable() {
	ctx := GetDBContext()
	db := GetDBIntance()

	colors := [8]models.Color{
		{ColorHex: "#FF6B6B", Name: "coral"},
		{ColorHex: "#DA77F2", Name: "orchid"},
		{ColorHex: "#9775FA", Name: "amethyst"},
		{ColorHex: "#5C7CFA", Name: "cobalt"},
		{ColorHex: "#66D9E8", Name: "aqua"},
		{ColorHex: "#8CE99A", Name: "lime"},
		{ColorHex: "#FFD43B", Name: "canary"},
		{ColorHex: "#FF922B", Name: "tangerine"},
	}

	for _, c := range colors {
		color := new(models.Color)
		err := db.NewSelect().Model(color).Where("name = ?", c.Name).Scan(ctx)
		if err != nil {
			db.NewInsert().Model(&c).Exec(ctx)
		}
	}
}
