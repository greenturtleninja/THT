package account

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/samborkent/uuidv7"
)

type Accounts struct {
	Accounts []Account `json:"accounts"`
}

type Account struct {
	AccountID        string  `json:"-"`
	AccountNumber    string  `json:"accountNumber"`
	SortCode         string  `json:"sortCode"`
	Name             string  `json:"name"`
	AccountType      string  `json:"accountType"`
	Balance          float64 `json:"balance"`
	Currency         string  `json:"currency"`
	CreatedTimestamp string  `json:"createdTimestamp"`
	UpdatedTimestamp string  `json:"updatedTimestamp"`
	Status           string  `json:"-"`
}

type ValidAccount struct {
	AccountID        string
	UserID           string
	Number           string
	CreatedTimestamp string
	UpdatedTimestamp string
}

type AccountsHandler interface {
	GetAccounts(db *sql.DB) error
	CreateAccount(db *sql.DB) error
	Response() string
}

// sum todays transactions and return this as the balance for the account
var GetAccountsQuery = `
	SELECT acc.name, acc.accountNumber, acc.sortCode, acc.accountType, acc.balance, acc.currency,
  		acc.createdTimestamp, acc.updatedTimestamp, COALESCE(SUM(txn.amount), 0) as txnBalance
	FROM accounts acc
	LEFT JOIN users_accounts usr USING (accountID)
	LEFT JOIN transactions txn USING (userID)
	WHERE txn.createdTimestamp IS NULL OR txn.createdTimestamp > date()
	GROUP BY acc.accountID, acc.name, acc.accountNumber, acc.sortCode, acc.accountType, acc.balance, acc.currency,
  		acc.createdTimestamp, acc.updatedTimestamp, acc.status
`

func GetAccounts(db *sql.DB) (Accounts, error) {
	var accounts Accounts

	rows, err := db.Query(GetAccountsQuery)
	if err != nil {
		return accounts, err
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var acc Account
		var currentBalance int64
		var txnBalance int64
		if err := rows.Scan(
			&acc.Name,
			&acc.AccountNumber,
			&acc.SortCode,
			&acc.AccountType,
			&currentBalance,
			&acc.Currency,
			&acc.CreatedTimestamp,
			&acc.UpdatedTimestamp,
			&txnBalance,
		); err != nil {
			return accounts, err
		}

		totalBalance := float64(currentBalance + txnBalance)
		if totalBalance > 0 {
			totalBalance = float64(totalBalance / 100)
		}
		acc.Balance = totalBalance

		accounts.Accounts = append(accounts.Accounts, acc)
	}
	if err = rows.Err(); err != nil {
		return accounts, err
	}

	return accounts, nil
}

func generateAccountNumber() string {
	randNumber := rand.IntN(9999999)
	return fmt.Sprintf("%08d", randNumber)
}

func generateSortCode() string {
	return "10-10-10"
}

var CreateAccountSQL = `INSERT INTO accounts 
(accountNumber, sortCode, name, accountType, currency, status, accountID, updatedTimestamp) VALUES 
(?, ?, ?, ?, ?, ?, ?, ?)`

func (acc *Account) CreateAccount(db *sql.DB) error {
	uniqueID := uuidv7.New()
	accUUID := uniqueID.String()
	acc.AccountID = accUUID
	acc.AccountNumber = generateAccountNumber()
	acc.SortCode = generateSortCode()
	acc.Status = "active"
	acc.Currency = "GBP"

	for {
		acc.AccountNumber = generateAccountNumber()
		_, err := IsValidAccount(db, acc.AccountNumber, acc.SortCode)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if err == sql.ErrNoRows {
			break
		}
	}

	acc.UpdatedTimestamp = time.Now().UTC().Format("2006-01-02 15:04:05")

	_, err := db.Exec(
		CreateAccountSQL,
		acc.AccountNumber,
		acc.SortCode,
		acc.Name,
		acc.AccountType,
		acc.Currency,
		acc.Status,
		acc.AccountID,
		acc.UpdatedTimestamp,
	)
	if err != nil {
		return err
	}

	validAccount, err := IsValidAccount(db, acc.AccountNumber, acc.SortCode)
	if err != nil {
		return err
	}

	acc.CreatedTimestamp = validAccount.CreatedTimestamp
	acc.UpdatedTimestamp = validAccount.UpdatedTimestamp

	return nil
}

var AccountExitsSQL = `SELECT acc.accountID, usr_acc.userID, accountNumber, COALESCE(createdTimestamp, ''), COALESCE(updatedTimestamp, '') 
FROM accounts acc 
INNER JOIN users_accounts usr_acc USING (accountID)
WHERE accountNumber = ? AND sortCode = ?`

func IsValidAccount(db *sql.DB, accountNumber, sortCode string) (ValidAccount, error) {
	var validAccount ValidAccount

	err := db.QueryRow(AccountExitsSQL, accountNumber, sortCode).Scan(
		&validAccount.AccountID,
		&validAccount.UserID,
		&validAccount.Number,
		&validAccount.CreatedTimestamp,
		&validAccount.UpdatedTimestamp,
	)

	return validAccount, err
}

func (acc *Account) Response() string {
	resp, _ := json.Marshal(acc)

	return string(resp)
}
