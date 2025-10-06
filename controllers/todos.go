package controllers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/marouane-ach/todo-go/db"
	"github.com/marouane-ach/todo-go/dtos"
	"github.com/marouane-ach/todo-go/models"
	"github.com/marouane-ach/todo-go/utils"
)

func CreateTodo(c echo.Context) error {
	ctx := context.Background()
	db := db.GetDBIntance()

	user, err := utils.ValidateToken(c, db, ctx)
	if err != nil {
		return err
	}

	todoDTO := new(dtos.TodoDTO)
	if err = c.Bind(todoDTO); err != nil {
		return err
	}

	todoListID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, &dtos.ErrorDTO{ErrorCode: 13, Description: "Todo list does not exist."})
	}

	todoList := new(models.TodoList)
	err = db.NewSelect().Model(todoList).Where("id = ?", todoListID).Scan(ctx)
	if err != nil {
		return c.JSON(http.StatusNotFound, &dtos.ErrorDTO{ErrorCode: 13, Description: "Todo list does not exist."})
	}

	if todoList.OwnerID != user.ID {
		return c.JSON(http.StatusUnauthorized, &dtos.ErrorDTO{ErrorCode: 14, Description: "Unauthorized."})
	}

	todo := &models.Todo{Text: todoDTO.Text, TodoListID: todoListID}
	_, err = db.NewInsert().Model(todo).Exec(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &dtos.ErrorDTO{ErrorCode: 15, Description: "Could not create todo."})
	}

	return c.JSON(http.StatusCreated, todo)
}
