package main

import (
	"github.com/labstack/echo/v4"
	"github.com/marouane-ach/todo-go/controllers"
	"github.com/marouane-ach/todo-go/db"
	_ "github.com/marouane-ach/todo-go/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title           Todo App Backend
// @version         1.0
// @description     Backend for a Todo App.
// @host            localhost:1323
// @BasePath        /

// @securityDefinitions.apiKey	BearerAuth
// @in                          header
// @name                        Authorization
func main() {
	db.CreateDBTables()
	db.SeedColorsTable()

	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.POST("/signup", controllers.Signup)

	e.POST("/login", controllers.Login)

	e.POST("/logout", controllers.Logout)

	e.POST("/todolists", controllers.CreateTodoList)

	e.GET("/todolists", controllers.GetUserTodoLists)

	e.GET("/todolists/:id", controllers.GetTodoListByID)

	e.DELETE("/todolists/:id", controllers.DeleteTodoList)

	e.POST("/todolists/:id/todos", controllers.CreateTodo)

	e.PUT("/todos/:id", controllers.UpdateTodo)

	e.Logger.Fatal(e.Start(":1323"))
}
