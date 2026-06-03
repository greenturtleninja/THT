package accountController

import (
	accountModel "THT/eaglebank/models/account"
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Account struct {
	DB            *sql.DB
	UserID        string
	AccountName   string
	Type          string
	AccountNumber string
}

func (a *Account) GetAccount(e echo.Context) error {
	accountMod := accountModel.Account{
		UserID: a.UserID,
	}

	if err := accountMod.GetAccountByUserId(a.DB); err != nil {
		e.Logger().Error("error getting account", "err", err)
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get account"})
	}

	return e.JSON(http.StatusOK, accountMod.Response())
}
