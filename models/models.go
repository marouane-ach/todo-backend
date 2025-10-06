package models

import (
	"time"

	"github.com/uptrace/bun"
)

type MyBaseModel struct {
	ID        int       `bun:"id,pk,autoincrement"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

type User struct {
	MyBaseModel
	bun.BaseModel `bun:"table:users"`

	Email          string `bun:",unique"`
	HashedPassword string
}

type Token struct {
	MyBaseModel
	bun.BaseModel `bun:"table:tokens"`

	Token   string `bun:",unique"`
	OwnerID int    `bun:",notnull"`
	Owner   *User  `bun:"rel:belongs-to,join:owner_id=id"`
}

type Color struct {
	MyBaseModel
	bun.BaseModel `bun:"table:colors"`

	Name     string `bun:",unique"`
	ColorHex string `bun:",unique,notnull"`
}

type TodoList struct {
	MyBaseModel
	bun.BaseModel `bun:"table:todo_lists"`

	Name    string `bun:",notnull"`
	ColorID int    `bun:",notnull"`
	Color   *Color `bun:"rel:belongs-to,join:color_id=id"`
	OwnerID int    `bun:",notnull"`
	Owner   *User  `bun:"rel:belongs-to,join:owner_id=id"`
	Todos   []Todo `bun:"rel:has-many,join:id=todo_list_id"`
}

type Todo struct {
	MyBaseModel
	bun.BaseModel `bun:"table:todos"`

	Text       string    `bun:",unique"`
	Completed  bool      `bun:"default:false"`
	TodoListID int       `bun:",notnull"`
	TodoList   *TodoList `bun:"rel:belongs-to,join:todo_list_id=id"`
}
