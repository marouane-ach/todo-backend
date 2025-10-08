package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/marouane-ach/todo-go/db"
	"github.com/marouane-ach/todo-go/dtos"
	"github.com/marouane-ach/todo-go/models"
	"github.com/marouane-ach/todo-go/utils"
)

// Create Todo godoc
// @Summary      Create a new todo in this todo list
// @Description  Accepts `text` as a JSON object and returns the created todo.
// @Tags         Todos
// @Param        todo body dtos.TodoDTO true "The todo list's name and color ID"
// @Param        id path int true "Todo List ID"
// @Accept       json
// @Produce      json
// @Success      201  {object}	models.Todo
// @Failure      401  {object}  dtos.ErrorDTO
// @Failure      404  {object}  dtos.ErrorDTO
// @Failure      500  {object}  dtos.ErrorDTO
// @Security	 BearerAuth
// @Router       /todolists/{id}/todos [post]
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

// Update todo godoc
// @Summary      Update this todo
// @Description  Accepts `text` and `completed` as a JSON object and returns the updated todo.
// @Tags         Todos
// @Param        todo body dtos.TodoDTO true "The todo's text and completed status"
// @Param        id path int true "Todo ID"
// @Accept       json
// @Produce      json
// @Success      200  {object}	models.Todo
// @Failure      401  {object}  dtos.ErrorDTO
// @Failure      404  {object}  dtos.ErrorDTO
// @Failure      500  {object}  dtos.ErrorDTO
// @Security	 BearerAuth
// @Router       /todos/{id} [put]
func UpdateTodo(c echo.Context) error {
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

	todoID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, &dtos.ErrorDTO{ErrorCode: 16, Description: "Todo does not exist."})
	}

	todo := new(models.Todo)
	err = db.NewSelect().Model(todo).Relation("TodoList").Where("todo.id = ?", todoID).Scan(ctx)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusNotFound, &dtos.ErrorDTO{ErrorCode: 17, Description: "Todo does not exist."})
	}

	if todo.TodoList.OwnerID != user.ID {
		return c.JSON(http.StatusUnauthorized, &dtos.ErrorDTO{ErrorCode: 18, Description: "Unauthorized."})
	}

	todo.Completed = todoDTO.Completed
	todo.Text = todoDTO.Text

	_, err = db.NewUpdate().Model(todo).Where("id = ?", todoID).Exec(ctx)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, &dtos.ErrorDTO{ErrorCode: 19, Description: "We encoutered a problem while updating the todo."})
	}

	return c.JSON(http.StatusOK, todo)
}
