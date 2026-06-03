package login

import (
	"database/sql"
)

type Login struct {
	UserID       string `json:"userId"`
	DisplayName  string `json:"displayName"`
	Login        string `json:"login"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

var GetLoginQuery = `
	SELECT 
		userID, displayName, login, email
	FROM logins
	WHERE login = $1 AND passwordHash = $2`

func (l *Login) CheckLogin(db *sql.DB) error {
	err := db.QueryRow(GetLoginQuery, l.Login, l.PasswordHash).Scan(
		&l.UserID,
		&l.DisplayName,
		&l.Login,
		&l.Email,
	)

	return err
}
