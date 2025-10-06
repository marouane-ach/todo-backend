package main

import (
	"github.com/labstack/echo/v4"
	"github.com/marouane-ach/todo-go/controllers"
	"github.com/marouane-ach/todo-go/db"
)

func main() {
	db.CreateDBTables()
	db.SeedColorsTable()

	e := echo.New()

	e.POST("/signup", controllers.Signup)

	e.POST("/login", controllers.Login)

	e.POST("/logout", controllers.Logout)

	e.POST("/todolists", controllers.CreateTodoList)

	e.GET("/todolists", controllers.GetUserTodoLists)

	e.POST("/todolists/:id/todos", controllers.CreateTodo)

	e.GET("/todolists/:id", controllers.GetTodoListByID)

	e.Logger.Fatal(e.Start(":1323"))
}
