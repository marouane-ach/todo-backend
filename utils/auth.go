package utils

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/marouane-ach/todo-go/dtos"
	"github.com/marouane-ach/todo-go/models"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPasswordBytes)
}

func GenerateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func ValidateToken(c echo.Context, db *bun.DB, ctx context.Context) (*models.User, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if len(authHeader) != 72 {
		return &models.User{}, c.JSON(http.StatusUnauthorized, &dtos.ErrorDTO{ErrorCode: 9, Description: "Invalid token."})
	}

	headerToken := authHeader[7:]
	user := new(models.User)

	token := new(models.Token)
	err := db.NewSelect().Model(token).Where("token = ?", headerToken).Scan(ctx)
	if err != nil {
		return user, c.JSON(http.StatusUnauthorized, &dtos.ErrorDTO{ErrorCode: 9, Description: "Invalid token."})
	}

	if err := db.NewSelect().Model(user).Where("id = ?", token.OwnerID).Scan(ctx); err != nil {
		return user, c.JSON(http.StatusInternalServerError, &dtos.ErrorDTO{ErrorCode: 10, Description: "Could not fetch user data."})
	}

	return user, nil
}
