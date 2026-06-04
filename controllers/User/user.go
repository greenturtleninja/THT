package userController

import (
	"database/sql"
	"net/http"

	userModel "THT/eaglebank/models/user"

	"github.com/labstack/echo/v4"
)

type User struct {
	DB          *sql.DB
	UserID      string
	TokenUserID string
}

func (u *User) GetUser(e echo.Context) error {
	userID := e.Param("userId")
	if u.TokenUserID != userID {
		return e.JSON(http.StatusForbidden, map[string]string{"error": "Forbidden"})
	}

	user := userModel.User{
		UserID: userID,
	}

	if err := user.GetUser(u.DB); err != nil {
		e.Logger().Error("error getting user", "err", err)

		return e.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get user"})
	}

	return e.JSON(http.StatusOK, user.Response())
}

// Create a new user with address details from login account
func (u *User) CreateUser(e echo.Context) error {
	userMod := userModel.User{}
	if err := e.Bind(&userMod); err != nil {
		e.Logger().Error("error binding user", "err", err)
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	userMod.UserID = u.TokenUserID
	if err := userMod.CreateUser(u.DB); err != nil {
		e.Logger().Error("error with creating user", "err", err)
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create user"})
	}
	return e.JSON(http.StatusCreated, userMod.Response())
}
