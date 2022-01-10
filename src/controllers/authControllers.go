package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c echo.Context) error {
	name, lastname, email, pwd := c.FormValue("name"), c.FormValue("lastname"), c.FormValue("email"), c.FormValue("password")
	// verificar si existe un correo similar (si ya se ha creado cuenta)
	existE := models.ExistsEmail(email)
	if existE {
		return c.JSON(400, config.SetResError(400, "Error: Email is already registered.", ""))
	}
	// Encriptar contraseña
	pwdHashB, errHashing := bcrypt.GenerateFromPassword([]byte(pwd), 10)
	if errHashing != nil {
		return c.JSON(500, config.SetResError(500, "Error: Encryption has failed", errHashing.Error()))
	}
	pwdH := string(pwdHashB)
	// crear ususario en BBDD
	errc := models.NewUser(name, lastname, email, pwdH, 3)
	if errc != nil {
		return c.JSON(500, config.SetResError(500, "Error: Do not create user", errc.Error()))
	}
	// comunicar que se ha creado
	return c.JSON(200, config.SetRes(200, "Created user"))
}

func Login(c echo.Context) error {
	email, password := c.FormValue("email"), c.FormValue("password")
	// verificar si ya se hizo login
	boolLog := true // se cree que SI ha iniciado sesion
	tkn, err := c.Cookie("token")
	if err != nil {
		boolLog = false // se cree que NO ha iniciado sesion
	}
	if boolLog {
		fmt.Println("token is:", tkn.Value)
		return c.JSON(400, config.SetRes(400, "user has already logged in"))
	}
	// buscar usuario
	user, errd := models.GetUserByEmail(email)
	if errd != nil {
		return c.JSON(400, config.SetResError(200, "user not found", errd.Error()))
	}
	// desincriptar y comparar la contraseña
	compareErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if compareErr != nil {
		return c.JSON(400, config.SetResError(400, "password incorrect", compareErr.Error()))
	}
	// crear JWT
	tokenString, expiresJWT, tokenErr := middlewares.CreateToken(user.ID.Hex(), uint8(user.Rol))
	if tokenErr != nil {
		return c.JSON(500, config.SetResError(500, "token do not created", tokenErr.Error()))
	}
	// crear cookie
	c.SetCookie(&http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expiresJWT,
	})

	return c.JSON(200, config.SetRes(200, tokenString))
}
func Logout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now().Add(0 * time.Millisecond),
	})

	return c.JSON(200, config.SetRes(200, "user logout"))
}
