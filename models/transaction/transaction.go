package transaction

import (
	"THT/eaglebank/models/account"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/samborkent/uuidv7"
)

type Transactions struct {
	Transactions  []Transaction `json:"transactions"`
	AccountNumber string        `json:"-"`
	SortCode      string        `json:"-"`
}

type Transaction struct {
	TransactionID    string  `json:"id"`
	Amount           float64 `json:"amount"`
	Currency         string  `json:"currency"`
	Type             string  `json:"type"`
	Reference        string  `json:"reference"`
	UserID           string  `json:"userId"`
	CreatedTimestamp string  `json:"createdTimestamp"`
	AccountNumber    string  `json:"-"`
	SortCode         string  `json:"-"`
}

var DefaultSortCode = "10-10-10"

type TransactionHandler interface {
	CreateTransaction(db *sql.DB) error
	GetTransactions(db *sql.DB) error
	Response() string
}

var GetTransactionsSQL = `SELECT 
	txn.transactionID, txn.amount, txn.currency, txn.type, txn.createdTimestamp, txn.reference, txn.userID 
FROM transactions txn
INNER JOIN users_accounts USING (userID)
INNER JOIN accounts acc USING (accountID)
WHERE acc.accountNumber = ? AND acc.sortCode = ?`

func (txns *Transactions) GetTransactions(db *sql.DB) error {
	rows, err := db.Query(GetTransactionsSQL, txns.AccountNumber, txns.SortCode)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var txn Transaction
		var txnAmount int64
		if err := rows.Scan(
			&txn.TransactionID,
			&txnAmount,
			&txn.Currency,
			&txn.Type,
			&txn.CreatedTimestamp,
			&txn.Reference,
			&txn.UserID,
		); err != nil {
			return err
		}

		amount := float64(txnAmount)
		if amount > 0 {
			amount = float64(amount / 100)
		}
		txn.Amount = amount

		txns.Transactions = append(txns.Transactions, txn)
	}

	return nil
}

var CreateTransactionSQL = `INSERT INTO transactions (transactionID, amount, currency, type, reference, userID)
VALUES (?, ?, ?, ?, ?, ?)`

// CreateTransaction
// checks if the passed account id is valid.
// Inserts the transaction details into the database
func (txn *Transaction) CreateTransaction(db *sql.DB) error {
	validAccount, err := account.IsValidAccount(db, txn.AccountNumber, txn.SortCode)
	if err != nil {
		return err
	}

	txn.UserID = validAccount.UserID
	uniqueID := uuidv7.New()
	tanUUID := uniqueID.String()
	txn.TransactionID = fmt.Sprintf("tan-%s", tanUUID)

	amount := int64(txn.Amount * 100)

	_, err = db.Exec(
		CreateTransactionSQL,
		txn.TransactionID,
		amount,
		txn.Currency,
		txn.Type,
		txn.Reference,
		txn.UserID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (txn *Transaction) Response() string {
	resp, _ := json.Marshal(txn)

	return string(resp)
}
