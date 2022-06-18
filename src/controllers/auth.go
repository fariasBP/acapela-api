package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/labstack/echo/v4"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c echo.Context) error {
	// obteniendo body json
	body := make(map[string]string)
	json.NewDecoder(c.Request().Body).Decode(&body)
	name, lastname, email, pwd, code, phoneString := body["name"], body["lastname"], body["email"], body["password"], body["code"], body["phone"]
	// convertir el phone a tipo int
	phone, _ := strconv.Atoi(phoneString)
	// verificar que no exista un email igual
	existEmail := models.ExistsEmail(email) // true si existe
	if existEmail {
		return c.JSON(400, config.SetResError(400, "Error: Email is already registered.", ""))
	}
	// verificar que no exista un code + phone iguales
	existPhone := models.ExistsPhone(code, phone)
	if existPhone {
		return c.JSON(400, config.SetResError(400, "Error: Phone is already registered.", ""))
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
	// obteniendo variables
	body := &middlewares.LoginValues{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// buscar usuario por numero
	user, err := models.GetUserByPhone(body.Code, body.Phone)
	if err != nil {
		return c.JSON(400, config.SetResError(200, "user not found", err.Error()))
	}
	// desincriptar y comparar la contraseña
	compareErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if compareErr != nil {
		return c.JSON(400, config.SetResError(400, "password incorrect", compareErr.Error()))
	}
	// crear JWT
	tokenString, expiresJWT, tokenErr := middlewares.CreateToken(user.ID.Hex(), uint8(user.Rol))
	if tokenErr != nil {
		return c.JSON(500, config.SetResError(500, "token do not created", tokenErr.Error()))
	}

	return c.JSON(200, config.SetResToken(400, "token created, expires in: ", tokenString, expiresJWT))
}
func RegisterUser(c echo.Context) error {
	// obteniendo variables
	body := &middlewares.RegisterValues{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	fmt.Println(body)
	// verificar que no exista un code + phone iguales
	existPhone := models.ExistsPhone(body.Code, body.Phone)
	if existPhone {
		return c.JSON(400, config.SetRes(400, "Error: Phone is already registered."))
	}
	//Generar primera contraseña
	// Generate a password that is 8 characters long with 10 digits, 10 symbols,
	// allowing upper and lower case letters, disallowing repeat characters.
	firstPwd, err := password.Generate(8, 3, 0, false, true)
	if err != nil {
		return c.JSON(400, config.SetResError(200, "first password don't create", err.Error()))
	}
	// crear ususario en BBDD
	errc := models.NewUserRegister(body.Name, body.Lastname, firstPwd, body.Code, body.Phone)
	if errc != nil {
		return c.JSON(500, config.SetResError(500, "Error: Do not create user register", errc.Error()))
	}
	// enviar el primer mensaje whatsapp

	return c.JSON(200, config.SetRes(200, "user registered"))
}
func Logout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now().Add(0 * time.Millisecond),
	})

	return c.JSON(200, config.SetRes(200, "user logout"))
}
