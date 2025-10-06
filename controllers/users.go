package controllers

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"net/mail"
	"unicode/utf8"

	"github.com/labstack/echo/v4"
	"github.com/marouane-ach/todo-go/db"
	"github.com/marouane-ach/todo-go/dtos"
	"github.com/marouane-ach/todo-go/models"
	"github.com/marouane-ach/todo-go/utils"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c echo.Context) error {
	ctx := context.Background()
	db := db.GetDBIntance()

	userDTO := new(dtos.UserDTO)
	if err := c.Bind(userDTO); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(userDTO.Email); err != nil {
		return c.JSON(http.StatusBadRequest, &dtos.ErrorDTO{ErrorCode: 1, Description: "Invalid email address."})
	}

	if utf8.RuneCountInString(userDTO.Password) < 8 {
		return c.JSON(http.StatusBadRequest, &dtos.ErrorDTO{ErrorCode: 2, Description: "Password must contain 8-24 characters."})
	}

	if utf8.RuneCountInString(userDTO.Password) > 24 {
		return c.JSON(http.StatusBadRequest, &dtos.ErrorDTO{ErrorCode: 2, Description: "Password must contain 8-24 characters."})
	}

	hashedPassword := utils.HashPassword(userDTO.Password)

	user := &models.User{Email: userDTO.Email, HashedPassword: hashedPassword}
	_, err := db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return c.JSON(http.StatusConflict, &dtos.ErrorDTO{ErrorCode: 4, Description: "An account with this email already exists."})
	}

	token := &models.Token{Token: utils.GenerateToken(), OwnerID: user.ID}
	_, err = db.NewInsert().Model(token).Exec(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &dtos.ErrorDTO{ErrorCode: 5, Description: "We encoutered a problem while creating your account."})
	}

	return c.JSON(http.StatusCreated, token.Token)
}

func Login(c echo.Context) error {
	ctx := context.Background()
	db := db.GetDBIntance()

	userDTO := new(dtos.UserDTO)
	if err := c.Bind(userDTO); err != nil {
		return err
	}

	user := new(models.User)
	err := db.NewSelect().Model(user).Where("email = ?", userDTO.Email).Scan(ctx)
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusNotFound, &dtos.ErrorDTO{ErrorCode: 6, Description: "An account with this email does not exist."})
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(userDTO.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, &dtos.ErrorDTO{ErrorCode: 7, Description: "Wrong password."})
	}

	b := make([]byte, 32)
	rand.Read(b)
	t := fmt.Sprintf("%x", b)
	token := &models.Token{Token: t, OwnerID: user.ID}
	_, err = db.NewInsert().Model(token).Exec(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &dtos.ErrorDTO{ErrorCode: 8, Description: "We encoutered a problem while logging you in."})
	}

	return c.JSON(http.StatusOK, t)
}

func Logout(c echo.Context) error {
	ctx := context.Background()
	db := db.GetDBIntance()

	_, err := utils.ValidateToken(c, db, ctx)
	if err != nil {
		return err
	}

	headerToken := c.Request().Header.Get("Authorization")[7:]
	token := new(models.Token)
	err = db.NewSelect().Model(token).Where("token = ?", headerToken).Scan(ctx)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.NewDelete().Model(token).WherePK().Exec(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &dtos.ErrorDTO{ErrorCode: 13, Description: "We encoutered a problem while logging you out."})
	}

	return c.String(http.StatusOK, "")
}
