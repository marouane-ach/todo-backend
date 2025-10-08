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

// Create Todo List godoc
// @Summary      Create a new todo list
// @Description  Accepts `name` and `color_id` as JSON and returns the created todo list.
// @Tags         Todo Lists
// @Param        todo_list body dtos.TodoListDTO true "The todo list's name and color ID"
// @Accept       json
// @Produce      json
// @Success      201  {object}	models.TodoList
// @Failure      400  {object}  dtos.ErrorDTO
// @Failure      401  {object}  dtos.ErrorDTO
// @Failure      500  {object}  dtos.ErrorDTO
// @Security	 BearerAuth
// @Router       /todolists [post]
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

// Get User Todo Lists godoc
// @Summary      Get user's todo lists.
// @Description  Returns a JSON array of the user's todo lists along with their associated todos.
// @Tags         Todo Lists
// @Accept       json
// @Produce      json
// @Success      200  {array}	models.TodoList
// @Failure      401  {object}  dtos.ErrorDTO
// @Failure      500  {object}  dtos.ErrorDTO
// @Security	 BearerAuth
// @Router       /todolists [get]
func GetUserTodoLists(c echo.Context) error {
	ctx := context.Background()
	db := db.GetDBIntance()

	user, err := utils.ValidateToken(c, db, ctx)
	if err != nil {
		fmt.Println("err:")
		fmt.Println(err)
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

// Get a Todo List by ID godoc
// @Summary      Get a single todo list by ID
// @Description  Returns a JSON object of a todo list along with its associated todos.
// @Tags         Todo Lists
// @Param        id path int true "Todo List ID"
// @Accept       json
// @Produce      json
// @Success      200  {array}	models.TodoList
// @Failure      401  {object}  dtos.ErrorDTO
// @Failure      500  {object}  dtos.ErrorDTO
// @Security	 BearerAuth
// @Router       /todolists/{id} [get]
func GetTodoListByID(c echo.Context) error {
	ctx := context.Background()
	db := db.GetDBIntance()

	user, err := utils.ValidateToken(c, db, ctx)
	if err != nil {
		return err
	}

	todoListID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, &dtos.ErrorDTO{ErrorCode: 13, Description: "Todo list does not exist."})
	}

	todoList := new(models.TodoList)
	err = db.NewSelect().
		Model(todoList).
		Where("id = ?", todoListID).
		Relation("Todos").
		Order("created_at DESC").
		Scan(ctx)
	if err != nil {
		return c.JSON(http.StatusNotFound, &dtos.ErrorDTO{ErrorCode: 13, Description: "Todo list does not exist."})
	}

	if todoList.OwnerID != user.ID {
		return c.JSON(http.StatusUnauthorized, &dtos.ErrorDTO{ErrorCode: 14, Description: "Unauthorized."})
	}

	return c.JSON(http.StatusOK, todoList)
}

// Delete Todo List godoc
// @Summary      Delete a todo list by ID
// @Description  Deletes a todo list and return JSON object of deleted todo list.
// @Tags         Todo Lists
// @Param        id path int true "Todo List ID"
// @Accept       json
// @Produce      json
// @Success      200  {array}	models.TodoList
// @Failure      401  {object}  dtos.ErrorDTO
// @Failure      500  {object}  dtos.ErrorDTO
// @Security	 BearerAuth
// @Router       /todolists/{id} [delete]
func DeleteTodoList(c echo.Context) error {
	ctx := context.Background()
	db := db.GetDBIntance()

	user, err := utils.ValidateToken(c, db, ctx)
	if err != nil {
		return err
	}

	todoListID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, &dtos.ErrorDTO{ErrorCode: 13, Description: "Todo list does not exist."})
	}

	todoList := new(models.TodoList)
	err = db.NewSelect().
		Model(todoList).
		Where("id = ?", todoListID).
		Relation("Todos").
		Order("created_at DESC").
		Scan(ctx)
	if err != nil {
		return c.JSON(http.StatusNotFound, &dtos.ErrorDTO{ErrorCode: 13, Description: "Todo list does not exist."})
	}

	if todoList.OwnerID != user.ID {
		return c.JSON(http.StatusUnauthorized, &dtos.ErrorDTO{ErrorCode: 14, Description: "Unauthorized."})
	}

	_, err = db.NewDelete().Model(todoList).Where("id = ?", todoListID).Exec(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &dtos.ErrorDTO{ErrorCode: 14, Description: "Could not delete todo."})
	}

	return c.JSON(http.StatusOK, todoList)
}
