package controllers

import (
	"encoding/json"
	"strings"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/fonini/go-capitalize/capitalize"
	"github.com/labstack/echo/v4"
)

// ---- LOGEADORES ----

// ---- login ----
func Login(c echo.Context) error {
	// obteniendo variables
	body := &models.User{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// buscar usuario por numero
	user, err := models.GetUserByPhone(body.CodePhone, body.Phone)
	if err != nil {
		return c.JSON(404, config.SetResError(404, "Error: Numero de telefono no registrado.", err.Error()))
	}
	if user.Code != body.Code {
		return c.JSON(400, config.SetRes(400, "Error: Codigo incorrecto"))
	}
	// crear JWT
	tokenString, expiresJWT, tokenErr := middlewares.CreateToken(user.ID.Hex(), uint8(user.Rol))
	if tokenErr != nil {
		return c.JSON(500, config.SetResError(500, "token no creado", tokenErr.Error()))
	}

	return c.JSON(200, config.SetResToken(200, "Se inicio session correctamente", tokenString, expiresJWT))
}

// ---- REGISTRADORES ----
// ---- auto registrador ----
// func Signup(c echo.Context) error {
// 	// obteniendo variables
// 	body := &middlewares.SignupValues{}
// 	d := c.Request().Body
// 	_ = json.NewDecoder(d).Decode(body)
// 	defer d.Close()
// 	// verificar que no exista un code + phone iguales
// 	existPhone := models.ExistsPhone(body.Code, body.Phone)
// 	if existPhone {
// 		return c.JSON(400, config.SetRes(400, "Error: El telefono ya ha sido registrado."))
// 	}
// 	// Encriptar contraseña
// 	pwdHashB, errHashing := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
// 	if errHashing != nil {
// 		return c.JSON(500, config.SetResError(500, "Error: La encriptación ha fallado.", errHashing.Error()))
// 	}
// 	pwdH := string(pwdHashB)
// 	// crear ususario en BBDD
// 	err := models.AutoClientRegistrar(body.Name, body.LastName, pwdH, body.Code, body.Phone)
// 	if err != nil {
// 		return c.JSON(500, config.SetResError(500, "Error: No se ha creado al usuario.", err.Error()))
// 	}
// 	// comunicar que se ha creado
// 	return c.JSON(200, config.SetRes(200, "Usuario Registrado"))
// }

// ---- registrador de clientes ----
func ClientRegistrar(c echo.Context) error {
	// obteniendo variables
	body := &models.User{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// verificar que no exista un code + phone iguales
	existPhone := models.ExistsPhone(body.CodePhone, body.Phone)
	if existPhone {
		return c.JSON(400, config.SetRes(400, "Error: El telefono ya ha sido registrado."))
	}
	// crear usuario en BBDD
	err := models.ClientRegistrar(strings.TrimSpace(body.Name), body.CodePhone, body.Phone)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: No se ha registrado al cliente", err.Error()))
	}
	// // enviar el primer mensaje whatsapp
	name, err := capitalize.Capitalize(body.Name)
	if err != nil {
		name = body.Name
	}
	// to := body.Code + strconv.Itoa(body.Phone)
	// err = middlewares.SendFirstMessageWelcome(to, name)
	// if err != nil {
	// 	return c.JSON(200, config.SetResError(500, "user registered but don't send first message whatsapp", err.Error()))
	// }
	// enviar el mesaje de bienvenida
	err = middlewares.SendWelcomeMessage(body.CodePhone, body.Phone, name)
	if err != nil {
		return c.JSON(200, config.SetResError(500, "Error: ususario fue registrado en BBDD pero no se envio el mensaje de bienvenida", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "Cliente registrado."))
}

// ---- registrador de empleados ----
func EmployeRegistrar(c echo.Context) error {
	// obteniendo variables
	body := &models.User{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// verificar que no exista un code + phone iguales
	existPhone := models.ExistsPhone(body.CodePhone, body.Phone)
	if existPhone {
		return c.JSON(400, config.SetRes(400, "Error: El telefono ya ha sido registrado."))
	}
	// crear ususario en BBDD
	err := models.EmployRegistrar(strings.TrimSpace(body.Name), body.CodePhone, body.Phone)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: No se ha registrado al empleado.", err.Error()))
	}
	// // enviar el primer mensaje whatsapp
	// name, err := capitalize.Capitalize(body.Name)
	// if err != nil {
	// 	name = body.Name
	// }
	// to := body.Code + strconv.Itoa(body.Phone)
	// err = middlewares.SendFirstMessageWelcome(to, name)
	// if err != nil {
	// 	return c.JSON(200, config.SetResError(500, "user registered but don't send first message whatsapp", err.Error()))
	// }

	return c.JSON(200, config.SetRes(200, "Empleado registrado."))
}

// ---- registrador de empleados administrativos ----
func AdminEmployeRegistrar(c echo.Context) error {
	// obteniendo variables
	body := &models.User{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// verificar que no exista un code + phone iguales
	existPhone := models.ExistsPhone(body.CodePhone, body.Phone)
	if existPhone {
		return c.JSON(400, config.SetRes(400, "Error: El telefono ya ha sido registrado."))
	}
	// crear ususario en BBDD
	err := models.AdminEmployRegistrar(strings.TrimSpace(body.Name), body.CodePhone, body.Phone)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: No se ha registrado al empleado.", err.Error()))
	}
	// // enviar el primer mensaje whatsapp
	// name, err := capitalize.Capitalize(body.Name)
	// if err != nil {
	// 	name = body.Name
	// }
	// to := body.Code + strconv.Itoa(body.Phone)
	// err = middlewares.SendFirstMessageWelcome(to, name)
	// if err != nil {
	// 	return c.JSON(200, config.SetResError(500, "user registered but don't send first message whatsapp", err.Error()))
	// }

	return c.JSON(200, config.SetRes(200, "Empleado administrativo registrado."))
}
