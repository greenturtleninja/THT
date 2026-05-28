package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/samborkent/uuidv7"
)

type User struct {
	UserID           string  `json:"id"`
	Name             string  `json:"name"`
	PhoneNumber      string  `json:"phoneNumber"`
	Email            string  `json:"email"`
	Status           string  `json:"-"`
	Address          Address `json:"address"`
	CreatedTimestamp string  `json:"CreatedTimestamp"`
	UpdatedTimestamp string  `json:"UpdatedTimestamp"`
}

type Address struct {
	Line1    string `json:"line1"`
	Line2    string `json:"line2"`
	Line3    string `json:"line3"`
	Town     string `json:"town"`
	County   string `json:"county"`
	Postcode string `json:"postcode"`
	Status   string `json:"-"`
}

type UsersToAccounts struct {
	UserID    string `json:"userId"`
	AccountID string `json:"accountId"`
}

type UserHandler interface {
	GetUser(db *sql.DB) error
	UpdateUser() error
	CreateUser(db *sql.DB) error
	DeleteUser() error
	Response() string
}

var GetUserQuery = `
	SELECT 
		user.userId, user.name, user.phoneNumber, user.email, user.createdTimestamp, user.updatedTimestamp,
		address.line1, address.line2, address.line3, address.town, address.county, address.postcode
	FROM users user
	LEFT JOIN addresses address USING (userId)
	WHERE user.userId = ?`

func (user *User) GetUser(db *sql.DB) error {
	err := db.QueryRow(GetUserQuery, user.UserID).Scan(
		&user.UserID,
		&user.Name,
		&user.PhoneNumber,
		&user.Email,
		&user.CreatedTimestamp,
		&user.UpdatedTimestamp,
		&user.Address.Line1,
		&user.Address.Line2,
		&user.Address.Line3,
		&user.Address.Town,
		&user.Address.County,
		&user.Address.Postcode,
	)

	return err
}

var LinkUserToAccountSQL = `INSERT INTO users_accounts (userID, accountID) VALUES (?, ?)`

func (usrToAcc *UsersToAccounts) LinkUserToAccount(db *sql.DB) error {
	_, err := db.Exec(
		LinkUserToAccountSQL,
		usrToAcc.UserID,
		usrToAcc.AccountID,
	)

	return err
}

func (user *User) UpdateUser() error {
	return nil
}

var CreateUserSQL = `INSERT INTO users (userID, name, email, phoneNumber, createdTimestamp, updatedTimestamp, status) 
VALUES (?, ?, ?, ?, ?, ?, ?)`

var CreateAddressSQL = `INSERT INTO addresses (addressID, userID, line1, line2, line3, town, county, postcode, status)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

func (user *User) CreateUser(db *sql.DB) error {
	currentTime := time.Now().UTC().Format("2006-01-02 15:04:05")
	user.UpdatedTimestamp = currentTime
	user.CreatedTimestamp = currentTime

	uniqueID := uuidv7.New()
	userUUID := uniqueID.String()
	user.UserID = fmt.Sprintf("usr-%s", userUUID)
	user.Status = "active"

	_, err := db.Exec(
		CreateUserSQL,
		user.UserID,
		user.Name,
		user.Email,
		user.PhoneNumber,
		user.CreatedTimestamp,
		user.UpdatedTimestamp,
		user.Status,
	)
	if err != nil {
		return err
	}

	uniqueID = uuidv7.New()
	addUUID := uniqueID.Short()
	addressId := fmt.Sprintf("usr-%s", addUUID)

	user.Address.Status = "active"

	_, err = db.Exec(
		CreateAddressSQL,
		addressId,
		user.UserID,
		user.Address.Line1,
		user.Address.Line2,
		user.Address.Line3,
		user.Address.Town,
		user.Address.County,
		user.Address.Postcode,
		user.Address.Status,
	)
	if err != nil {
		return err
	}

	return nil
}

func (user *User) DeleteUser() error {
	return nil
}

func (user *User) Response() string {
	resp, _ := json.Marshal(user)

	return string(resp)
}
