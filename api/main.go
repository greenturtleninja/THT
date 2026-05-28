package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"THT/eaglebank/models/account"
	"THT/eaglebank/models/transaction"
	"THT/eaglebank/models/user"

	_ "github.com/glebarez/go-sqlite"
)

type errorMessage struct {
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details"`
}

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

var myslog *slog.Logger
var db *sql.DB

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func sqlLiteVersion(w http.ResponseWriter, req *http.Request) {
	var sqliteVersion string
	err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' LIMIT 1").Scan(&sqliteVersion)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Fprintf(w, "%s", sqliteVersion)
}

func getUserHandler(w http.ResponseWriter, req *http.Request) {
	userID := req.PathValue("userId")

	user := user.User{
		UserID: userID,
	}

	if err := user.GetUser(db); err != nil {
		myslog.Error("error getting user", "err", err)

		fmt.Fprintf(w, "")
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", user.Response())
}

func createUserHandler(w http.ResponseWriter, req *http.Request) {
	user := user.User{}
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		myslog.Error("error with decoding user body", "err", err)
		errorMsg := errorMessage{
			Message: "invalid user details",
			Details: []ErrorDetail{
				{
					Field:   "userBody",
					Message: "invalid body",
					Type:    "string",
				},
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", errorResponse(errorMsg))
		return
	}

	err = user.CreateUser(db)
	if err != nil {
		myslog.Error("error with decoding user body", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		errorMsg := errorMessage{
			Message: "error creating user",
		}
		fmt.Fprintf(w, "%s", errorResponse(errorMsg))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", user.Response())
}

func linkUsersToAccount(w http.ResponseWriter, req *http.Request) {
	usersToAccounts := user.UsersToAccounts{}
	err := json.NewDecoder(req.Body).Decode(&usersToAccounts)
	if err != nil {
		myslog.Error("error with decoding user account link", "err", err)
		errorMsg := errorMessage{
			Message: "invalid user-account details",
			Details: []ErrorDetail{
				{
					Field:   "user-accountBody",
					Message: "invalid body",
					Type:    "string",
				},
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", errorResponse(errorMsg))
		return
	}

	err = usersToAccounts.LinkUserToAccount(db)
	if err != nil {
		myslog.Error("error with decoding user body", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		errorMsg := errorMessage{
			Message: "error creating user",
		}
		fmt.Fprintf(w, "%s", errorResponse(errorMsg))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", "{\"status\": \"OK\"}")

}

func getAccountsHandler(w http.ResponseWriter, req *http.Request) {
	accounts, err := account.GetAccounts(db)
	if err != nil {
		myslog.Error("error getting account", "err", err)

		fmt.Fprintf(w, "")
	}

	resp, _ := json.Marshal(accounts)
	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, "%s", resp)
}

func createAccountsHandler(w http.ResponseWriter, req *http.Request) {
	account := account.Account{}
	err := json.NewDecoder(req.Body).Decode(&account)
	if err != nil {
		myslog.Error("error with decoding account body", "err", err)
		errorMsg := errorMessage{
			Message: "invalid account details",
			Details: []ErrorDetail{
				{
					Field:   "accountBody",
					Message: "invalid body",
					Type:    "string",
				},
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", errorResponse(errorMsg))
		return
	}

	err = account.CreateAccount(db)
	if err != nil {
		myslog.Error("error with decoding account body", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		errorMsg := errorMessage{
			Message: "error creating account",
		}
		fmt.Fprintf(w, "%s", errorResponse(errorMsg))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", account.Response())
}

func getTransactionsHandler(w http.ResponseWriter, req *http.Request) {
	accountNumber := req.PathValue("accountNumber")

	txn := transaction.Transactions{
		AccountNumber: accountNumber,
		SortCode:      transaction.DefaultSortCode,
	}

	if err := txn.GetTransactions(db); err != nil {
		myslog.Error("error getting transactions", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "")
	}

	w.Header().Set("Content-Type", "application/json")
	resp, _ := json.Marshal(txn)
	fmt.Fprintf(w, "%s", resp)
}

func createTransactionHandler(w http.ResponseWriter, req *http.Request) {
	transaction := transaction.Transaction{
		AccountNumber: req.PathValue("accountNumber"),
		SortCode:      transaction.DefaultSortCode,
	}
	err := json.NewDecoder(req.Body).Decode(&transaction)
	if err != nil {
		myslog.Error("error with decoding transaction body", "err", err)
		errorMsg := errorMessage{
			Message: "invalid account details",
			Details: []ErrorDetail{
				{
					Field:   "transactionBody",
					Message: "invalid body",
					Type:    "string",
				},
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", errorResponse(errorMsg))
		return
	}

	err = transaction.CreateTransaction(db)
	if err != nil {
		myslog.Error("error with creating transaction", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		errorMsg := errorMessage{
			Message: "error creating transaction",
		}
		fmt.Fprintf(w, "%s", errorResponse(errorMsg))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", transaction.Response())

}

func errorResponse(errorMessage errorMessage) string {
	message, err := json.Marshal(errorMessage)
	if err != nil {
		myslog.Error(fmt.Sprintf("error marshalling: %s", err.Error()))
		return ""
	}

	return string(message)
}

func main() {
	// auth JWT here

	jsonHandler := slog.NewJSONHandler(os.Stderr, nil)
	myslog = slog.New(jsonHandler)

	var err error

	// Connect to the SQLite database
	db, err = sql.Open("sqlite", "./eagle_bank.db")
	if err != nil {
		myslog.Error(fmt.Sprintf("error connecting to db: %s", err.Error()))
	}

	defer db.Close()

	// API calls
	http.HandleFunc("GET /v1/dbcheck", sqlLiteVersion)
	http.HandleFunc("GET /v1/accounts/{accountNumber}", hello)
	// todo
	http.HandleFunc("PATCH /v1/accounts/{accountNumber}", hello)
	http.HandleFunc("DELETE /v1/accounts/{accountNumber}", hello)
	http.HandleFunc("DELETE /v1/users/{userId}", hello)
	http.HandleFunc("GET /v1/accounts/{accountNumber}/transactions/{transactionId}", hello)

	// working
	http.HandleFunc("GET /v1/accounts", getAccountsHandler)
	http.HandleFunc("POST /v1/accounts", createAccountsHandler)

	http.HandleFunc("POST /v1/accounts/{accountNumber}/transactions", createTransactionHandler)
	http.HandleFunc("GET /v1/accounts/{accountNumber}/transactions", getTransactionsHandler)

	http.HandleFunc("POST /v1/users/{userId}/accounts/{accountId}", linkUsersToAccount)

	http.HandleFunc("GET /v1/users/{userId}", getUserHandler)
	http.HandleFunc("POST /v1/users", createUserHandler)

	http.ListenAndServe(":8090", nil)
}
