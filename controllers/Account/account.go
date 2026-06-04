package accountController

import (
	accountModel "THT/eaglebank/models/account"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samborkent/uuidv7"
)

type Account struct {
	DB            *sql.DB
	UserID        string
	AccountName   string `json:"name"`
	Type          string `json:"accountType"`
	AccountNumber string
}

func (a *Account) GetAccount(e echo.Context) error {
	accountMod := accountModel.Account{
		UserID:        a.UserID,
		AccountNumber: a.AccountNumber,
	}

	if err := accountMod.GetAccountByUserId(a.DB); err != nil {
		e.Logger().Error("error getting account", "err", err)
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get account"})
	}

	return e.JSON(http.StatusOK, Response(accountMod))
}

func (a *Account) GetAccounts(e echo.Context) error {
	accounts, err := accountModel.GetAccountsByUserId(a.DB, a.UserID)
	if err != nil {
		e.Logger().Error("error getting accounts", "err", err)
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get accounts"})
	}

	return e.JSON(http.StatusOK, Response(accounts))
}

var GenerateUUID = func() string {
	uniqueID := uuidv7.New() // move to function so can override for unit tests
	return uniqueID.String()
}

var GenerateAccountNumber = func() string {
	randNumber := rand.IntN(9999999)
	return fmt.Sprintf("%08d", randNumber)
}

var GenerateSortCode = func() string {
	return "10-10-10"
}

func (a *Account) CreateAccount(e echo.Context) error {
	newAccount := accountModel.Account{
		UserID:      a.UserID,
		AccountType: a.Type,
		Name:        a.AccountName,
	}

	newAccount.AccountID = GenerateUUID()
	newAccount.AccountNumber = GenerateAccountNumber()
	newAccount.SortCode = GenerateSortCode()
	newAccount.Status = "active"
	newAccount.Currency = "GBP"

	for {
		newAccount.AccountNumber = GenerateAccountNumber()
		_, err := accountModel.IsValidAccount(a.DB, newAccount.AccountNumber, newAccount.SortCode)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if err == sql.ErrNoRows {
			break
		}
	}

	//CreateAccount
	if err := newAccount.CreateAccount(a.DB); err != nil {
		e.Logger().Error("error creating account", "err", err)
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to creating account"})
	}

	//LinkAccount
	if err := newAccount.LinkUserToAccount(a.DB); err != nil {
		e.Logger().Error("error linking account to user", "err", err)
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to linking account to user"})
	}

	return e.JSON(http.StatusCreated, Response(newAccount))
}

func Response(accounts interface{}) string {
	resp, _ := json.Marshal(accounts)

	return string(resp)
}
