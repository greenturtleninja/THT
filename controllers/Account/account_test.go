package accountController

import (
	accountModel "THT/eaglebank/models/account"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
)

func TestGetAccount(t *testing.T) {
	// Open new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("error creating mock database")
		return
	}

	columns := []string{"name", "account_number", "sort_code", "account_type", "balance", "currency",
		"created_timestamp", "updated_timestamp", "txn_balance"}

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	testCases := []struct {
		testname    string
		testAccount Account
		request     *http.Request
		response    *httptest.ResponseRecorder
		respAccount accountModel.Account
		expectQuery *sqlmock.ExpectedQuery
	}{
		{
			testname: "successfully get account",
			testAccount: Account{
				UserID:        "test-user-id",
				AccountNumber: "12345678",
				DB:            db,
			},
			request:  req,
			response: w,
			respAccount: accountModel.Account{
				Name:             "Test Account",
				AccountNumber:    "12345678",
				SortCode:         "10-10-10",
				AccountType:      "personal",
				Balance:          1000.00,
				Currency:         "GBP",
				CreatedTimestamp: "2024-01-01T00:00:00Z",
				UpdatedTimestamp: "2024-01-01T00:00:00Z",
				UserID:           "test-user-id",
			},
			expectQuery: func() *sqlmock.ExpectedQuery {
				return mock.ExpectQuery(regexp.QuoteMeta(accountModel.GetAccountByUserIdQuery)).
					WithArgs("test-user-id", "12345678").
					WillReturnRows(sqlmock.NewRows(columns).AddRow(
						"Test Account",
						"12345678",
						"10-10-10",
						"personal",
						0,
						"GBP",
						"2024-01-01T00:00:00Z",
						"2024-01-01T00:00:00Z",
						int64(100000),
					))
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testname, func(t *testing.T) {
			e := echo.New()

			err = tc.testAccount.GetAccount(e.NewContext(req, w))
			if err != nil {
				t.Errorf("test: %s - expected error when passing nil db, got nil", tc.testname)
			}

			expectedResponse := fmt.Sprintf("\"%s\"\n", Response(tc.respAccount))

			outputResponse := strings.ReplaceAll(w.Body.String(), "\\", "")

			if outputResponse != string(expectedResponse) {
				t.Errorf("test: %s - expected response body -%s-, got -%s-", tc.testname, string(expectedResponse), string(outputResponse))
			}

			// ensure all expectations have been met
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("test: %s - unmet expectation error: %s", tc.testname, err)
			}
		})
	}
}

func TestGetAccounts(t *testing.T) {
	// Open new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("error creating mock database")
		return
	}

	columns := []string{"name", "account_number", "sort_code", "account_type", "balance", "currency",
		"created_timestamp", "updated_timestamp", "txn_balance"}

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	testCases := []struct {
		testname    string
		testAccount Account
		request     *http.Request
		response    *httptest.ResponseRecorder
		respAccount accountModel.Accounts
		expectQuery *sqlmock.ExpectedQuery
	}{
		{
			testname: "successfully get all accounts for a user",
			testAccount: Account{
				UserID: "test-user-id",
				DB:     db,
			},
			request:  req,
			response: w,
			respAccount: accountModel.Accounts{
				Accounts: []accountModel.Account{
					{
						Name:             "Test Account",
						AccountNumber:    "12345678",
						SortCode:         "10-10-10",
						AccountType:      "personal",
						Balance:          1000.00,
						Currency:         "GBP",
						CreatedTimestamp: "2024-01-01T00:00:00Z",
						UpdatedTimestamp: "2024-01-01T00:00:00Z",
						UserID:           "test-user-id",
					},
					{
						Name:             "Test Account 2",
						AccountNumber:    "78901234",
						SortCode:         "10-10-10",
						AccountType:      "business",
						Balance:          90.00,
						Currency:         "GBP",
						CreatedTimestamp: "2024-01-02T00:00:00Z",
						UpdatedTimestamp: "2024-01-02T00:00:00Z",
						UserID:           "test-user-id",
					},
				},
			},
			expectQuery: func() *sqlmock.ExpectedQuery {
				return mock.ExpectQuery(regexp.QuoteMeta(accountModel.GetAccountsQuery)).
					WithArgs("test-user-id").
					WillReturnRows(sqlmock.NewRows(columns).AddRow(
						"Test Account",
						"12345678",
						"10-10-10",
						"personal",
						0,
						"GBP",
						"2024-01-01T00:00:00Z",
						"2024-01-01T00:00:00Z",
						int64(100000),
					).AddRow(
						"Test Account 2",
						"78901234",
						"10-10-10",
						"business",
						9000,
						"GBP",
						"2024-01-02T00:00:00Z",
						"2024-01-02T00:00:00Z",
						int64(0),
					))
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testname, func(t *testing.T) {
			e := echo.New()

			err = tc.testAccount.GetAccounts(e.NewContext(req, w))
			if err != nil {
				t.Errorf("test: %s - expected error when passing nil db, got nil", tc.testname)
			}

			expectedResponse := fmt.Sprintf("\"%s\"\n", Response(tc.respAccount))

			outputResponse := strings.ReplaceAll(w.Body.String(), "\\", "")

			if outputResponse != string(expectedResponse) {
				t.Errorf("test: %s - expected response body -%s-, got -%s-", tc.testname, string(expectedResponse), string(outputResponse))
			}

			// ensure all expectations have been met
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("test: %s - unmet expectation error: %s", tc.testname, err)
			}
		})
	}
}

// TODO: fix this
func TestCreateAccount(t *testing.T) {
	// Open new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("error creating mock database")
		return
	}

	checkColumns := []string{"acc.accountID", "usr_acc.userID", "accountNumber", "COALESCE(createdTimestamp, '')", "COALESCE(updatedTimestamp, '')"}
	accountColumns := []string{"accountNumber", "sortCode", "name", "accountType", "currency", "status", "accountID", "updatedTimestamp"}
	LinkColumns := []string{"userID", "accountID"}

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	testCases := []struct {
		testname           string
		testAccount        Account
		request            *http.Request
		response           *httptest.ResponseRecorder
		respAccount        accountModel.Account
		expectQuery        func()
		mockAccountDetails func()
	}{
		{
			testname: "successfully get all accounts for a user",
			testAccount: Account{
				UserID:      "test-user-id",
				AccountName: "Test Account",
				Type:        "personal",
				DB:          db,
			},
			request:  req,
			response: w,
			respAccount: accountModel.Account{
				Name:             "Test Account",
				AccountNumber:    "12345678",
				SortCode:         "10-10-10",
				AccountType:      "personal",
				Balance:          0.00,
				Currency:         "GBP",
				CreatedTimestamp: "2024-01-01T00:00:00Z",
				UpdatedTimestamp: "2024-01-01T00:00:00Z",
				UserID:           "test-user-id",
			},
			mockAccountDetails: func() {
				GenerateUUID = func() string {
					return ""
				}
				GenerateAccountNumber = func() string {
					return ""
				}
			},
			expectQuery: func() {
				mock.ExpectQuery(regexp.QuoteMeta(accountModel.AccountExitsSQL)).
					WithArgs("test-user-id").
					WillReturnRows()

				mock.ExpectExec(regexp.QuoteMeta(accountModel.CreateAccountSQL)).
					WithArgs("12345678", "10-10-10", name, accountType, currency, status, "12345678", updatedTimestamp)

				mock.ExpectExec(regexp.QuoteMeta(accountModel.LinkUserToAccountSQL)).
					WithArgs("test-user-id", "12345678")

			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testname, func(t *testing.T) {
			e := echo.New()

			err = tc.testAccount.GetAccounts(e.NewContext(req, w))
			if err != nil {
				t.Errorf("test: %s - expected error when passing nil db, got nil", tc.testname)
			}

			expectedResponse := fmt.Sprintf("\"%s\"\n", Response(tc.respAccount))

			outputResponse := strings.ReplaceAll(w.Body.String(), "\\", "")

			if outputResponse != string(expectedResponse) {
				t.Errorf("test: %s - expected response body -%s-, got -%s-", tc.testname, string(expectedResponse), string(outputResponse))
			}

			// ensure all expectations have been met
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("test: %s - unmet expectation error: %s", tc.testname, err)
			}
		})
	}
}
