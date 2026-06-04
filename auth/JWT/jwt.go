package authJwt

import (
	login "THT/eaglebank/models/Login"
	"crypto/rsa"
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type LoginRequest struct {
	DB       *sql.DB
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	RsaPublicKey    *rsa.PublicKey
	RsaPrivateKey   *rsa.PrivateKey
	SignatureMethod = jwt.SigningMethodRS256
)

type JwtCustomClaims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

func Init() {
	pubKey, err := os.ReadFile("./certs/pub_key.pem")
	if err != nil {
		panic(fmt.Sprintf("error reading public key: %v", err))
	}

	RsaPublicKey, err = jwt.ParseRSAPublicKeyFromPEM(pubKey)
	if err != nil {
		panic(fmt.Sprintf("error parsing public key: %v", err))
	}

	privKey, err := os.ReadFile("./certs/priv_key.pem")
	if err != nil {
		panic(fmt.Sprintf("error reading private key: %v", err))
	}

	RsaPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privKey)
	if err != nil {
		panic(fmt.Sprintf("error parsing private key: %v", err))
	}
}

func PasswordHash(password string) string {
	hash := sha512.New()
	hash.Write([]byte(password))
	hashedData := hash.Sum(nil)
	passwordHash := hex.EncodeToString(hashedData)

	return passwordHash
}

func GetUserIdFromToken(e echo.Context) string {
	token := e.Get("user").(*jwt.Token)
	claims := token.Claims.(*JwtCustomClaims)
	userIDFromToken := claims.UserID

	return userIDFromToken
}

func (l *LoginRequest) Login(req echo.Context) error {
	username := req.FormValue("username")
	password := req.FormValue("password")

	if username == "" || password == "" {
		return echo.ErrBadRequest
	}

	passowrdHash := PasswordHash(password)

	login := login.Login{
		Login:        username,
		PasswordHash: passowrdHash,
	}

	err := login.CheckLogin(l.DB)
	if err != nil {
		req.Logger().Error("login failed", "error", err, "username", username, "passwordHash", string(passowrdHash))
		return echo.ErrUnauthorized
	}

	claims := &JwtCustomClaims{
		UserID: login.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   login.Login,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	token := jwt.NewWithClaims(SignatureMethod, claims)

	t, err := token.SignedString(RsaPrivateKey)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return req.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}
