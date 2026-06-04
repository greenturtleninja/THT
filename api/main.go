package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	authJwt "THT/eaglebank/auth/JWT"
	accountController "THT/eaglebank/controllers/Account"
	userController "THT/eaglebank/controllers/User"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samborkent/uuidv7"

	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	dbuser = "eagle_bank"
	dbname = "eagle_bank"
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

func dbVersion(req echo.Context) error {
	var timeNow string
	err := db.QueryRow("SELECT now()").Scan(&timeNow)
	if err != nil {
		myslog.Error("error getting database version", "err", err)
		return echo.ErrInternalServerError
	}

	return req.JSON(http.StatusOK, echo.Map{
		"version": timeNow,
	})
}

func newUserId(req echo.Context) error {
	uniqueID := uuidv7.New()
	usrUUID := uniqueID.String()
	usrID := fmt.Sprintf("usr-%s", usrUUID)

	return req.JSON(http.StatusOK, echo.Map{
		"user-id": usrID,
	})
}

func passwordHash(req echo.Context) error {
	password := req.QueryParam("password")
	if password == "" {
		return req.JSON(http.StatusBadRequest, echo.Map{
			"error": "password query parameter is required",
		})
	}

	passHash := authJwt.PasswordHash(password)

	return req.JSON(http.StatusOK, echo.Map{
		"hash": passHash,
	})
}

func errorResponse(errorMessage errorMessage) string {
	message, err := json.Marshal(errorMessage)
	if err != nil {
		myslog.Error(fmt.Sprintf("error marshalling: %s", err.Error()))
		return ""
	}

	return string(message)
}

func init() {
	authJwt.Init()
	jsonHandler := slog.NewJSONHandler(os.Stderr, nil)
	myslog = slog.New(jsonHandler)

	var err error

	dbPassword := os.Getenv("DBPASSWORD")
	if dbPassword == "" {
		myslog.Error("error connecting to db: no password provided")
		return
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, dbuser, dbPassword, dbname)

	// Connect to the PostgreSQL database
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		myslog.Error(fmt.Sprintf("error connecting to db: %s", err.Error()))
		return
	}

	myslog.Info("successfully connected to db")
}

func loginRoute(parent *echo.Group, db *sql.DB) {
	u := authJwt.LoginRequest{
		DB: db,
	}

	// register all routes on controller
	parent.POST("/", u.Login)
}

func registerApiRoutes(parent *echo.Group, db *sql.DB) {
	// return user details if you are the current user
	// todo: or you are admin
	parent.GET("/users/:userId", func(c echo.Context) error {
		user := userController.User{
			DB:          db,
			UserID:      c.Param("userId"),
			TokenUserID: authJwt.GetUserIdFromToken(c),
		}
		return user.GetUser(c)
	})

	parent.POST("/users", func(c echo.Context) error {
		user := userController.User{
			DB:          db,
			TokenUserID: authJwt.GetUserIdFromToken(c),
		}

		return user.CreateUser(c)
	})

	parent.GET("/accounts/:accountId", func(c echo.Context) error {
		account := accountController.Account{
			DB:            db,
			AccountNumber: c.Param("accountId"),
			UserID:        authJwt.GetUserIdFromToken(c),
		}
		return account.GetAccount(c)
	})
	parent.POST("/accounts", func(c echo.Context) error {
		account := accountController.Account{
			UserID: authJwt.GetUserIdFromToken(c),
			DB:     db,
		}
		if err := c.Bind(&account); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "incoorect details to create account"})
		}

		return account.CreateAccount(c)
	})

}

func tools(parent *echo.Group) {
	parent.GET("/dbversion", dbVersion)
	parent.GET("/userID", newUserId)
	parent.GET("/passwordHash", passwordHash)
}

func main() {
	e := echo.New()
	e.Use(middleware.RequestLogger())

	toolsGroup := e.Group("/tools")
	tools(toolsGroup)

	login := e.Group("/login")
	loginRoute(login, db)

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(authJwt.JwtCustomClaims)
		},
		SigningKey:    authJwt.RsaPublicKey,
		SigningMethod: authJwt.SignatureMethod.Name,
	}

	api := e.Group("/v1")
	api.Use(echojwt.WithConfig(config))
	registerApiRoutes(api, db)

	if err := e.Start("127.0.0.1:8090"); err != nil {
		myslog.Error(fmt.Sprintf("error starting server: %s", err.Error()))
	}

	defer db.Close()
}
