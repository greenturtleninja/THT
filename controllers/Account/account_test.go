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
	}{
		{
			testname: "successfully get account",
			testAccount: Account{
				UserID: "test-user-id",
				DB:     db,
			},
			request:  req,
			response: w,
			respAccount: accountModel.Account{
				Name:             "Test Account",
				AccountNumber:    "12345678",
				SortCode:         "10-10-10",
				AccountType:      "checking",
				Balance:          1000.00,
				Currency:         "GBP",
				CreatedTimestamp: "2024-01-01T00:00:00Z",
				UpdatedTimestamp: "2024-01-01T00:00:00Z",
				UserID:           "test-user-id",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testname, func(t *testing.T) {
			e := echo.New()

			mock.ExpectQuery(regexp.QuoteMeta(accountModel.GetAccountByUserIdQuery)).
				WithArgs(tc.testAccount.UserID).
				WillReturnRows(sqlmock.NewRows(columns).AddRow(
					tc.respAccount.Name,
					tc.respAccount.AccountNumber,
					tc.respAccount.SortCode,
					tc.respAccount.AccountType,
					0,
					tc.respAccount.Currency,
					tc.respAccount.CreatedTimestamp,
					tc.respAccount.UpdatedTimestamp,
					int64(tc.respAccount.Balance*100),
				))

			err = tc.testAccount.GetAccount(e.NewContext(req, w))
			if err != nil {
				t.Error("expected error when passing nil db, got nil")
			}

			expectedResponse := fmt.Sprintf("\"%s\"\n", tc.respAccount.Response())

			outputResponse := strings.ReplaceAll(w.Body.String(), "\\", "")

			if outputResponse != string(expectedResponse) {
				t.Errorf("expected response body -%s-, got -%s-", string(expectedResponse), string(outputResponse))
			}

			// ensure all expectations have been met
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectation error: %s", err)
			}
		})
	}
}
