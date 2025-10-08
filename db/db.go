package db

import (
	"context"
	"database/sql"

	"github.com/marouane-ach/todo-go/models"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

var dbInstance *bun.DB
var dbContext context.Context

func GetDBIntance() *bun.DB {
	if dbInstance == nil {
		sqldb, err := sql.Open(sqliteshim.ShimName, "file:dev.db?cache=shared&mode=rwc")
		if err != nil {
			panic(err)
		}
		dbInstance = bun.NewDB(sqldb, sqlitedialect.New())
	}

	return dbInstance
}

func GetAppContext() context.Context {
	if dbContext == nil {
		dbContext = context.Background()
	}

	return dbContext
}

func CreateDBTables() {
	ctx := GetAppContext()
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
	ctx := GetAppContext()
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
