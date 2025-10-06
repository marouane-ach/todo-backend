package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/marouane-ach/todo-go/db"
	"github.com/marouane-ach/todo-go/dtos"
	"github.com/marouane-ach/todo-go/models"
	"github.com/marouane-ach/todo-go/utils"
)

func CreateTodoList(c echo.Context) error {
	ctx := context.Background()
	db := db.GetDBIntance()

	user, err := utils.ValidateToken(c, db, ctx)
	if err != nil {
		return err
	}

	todoListDTO := new(dtos.TodoListDTO)
	if err = c.Bind(todoListDTO); err != nil {
		return err
	}

	color := new(models.Color)
	err = db.NewSelect().Model(color).Where("id = ?", todoListDTO.ColorID).Scan(ctx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &dtos.ErrorDTO{ErrorCode: 11, Description: "Invalid color ID."})
	}

	todoList := &models.TodoList{Name: todoListDTO.Name, ColorID: todoListDTO.ColorID, OwnerID: user.ID}
	_, err = db.NewInsert().Model(todoList).Exec(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &dtos.ErrorDTO{ErrorCode: 12, Description: "Could not create todo list."})
	}

	return c.JSON(http.StatusCreated, todoList)
}

func GetUserTodoLists(c echo.Context) error {
	ctx := context.Background()
	db := db.GetDBIntance()

	user, err := utils.ValidateToken(c, db, ctx)
	if err != nil {
		return err
	}

	var todoLists []models.TodoList
	err = db.NewSelect().
		Model(&todoLists).
		Where("owner_id = ?", user.ID).
		Relation("Todos").
		Order("created_at DESC").
		Scan(ctx)

	if err != nil {
		fmt.Println(err)
	}

	return c.JSON(http.StatusOK, todoLists)
}
