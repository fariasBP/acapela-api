package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

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
	fmt.Println(body.Code)
	user, err := models.GetUserByPhone(body.Phone)
	if err != nil {
		return c.JSON(404, config.SetResError(404, "Error: Numero de telefono no registrado.", err.Error()))
	}
	// verficando si el codigo es identico
	if user.Code != body.Code {
		return c.JSON(400, config.SetRes(400, "Error: Codigo incorrecto"))
	}
	// verificar si el codigo se envia antes de 1 hora
	if time.Now().UTC().Before(user.CodeDate.Add(time.Hour)) {
		return c.JSON(400, config.SetRes(400, "Error: Esta enviando un codigo caducado."))
	}
	// crear JWT
	tokenString, expiresJWT, tokenErr := middlewares.CreateToken(user.ID.Hex(), uint8(user.Rol))
	if tokenErr != nil {
		return c.JSON(500, config.SetResError(500, "token no creado", tokenErr.Error()))
	}

	return c.JSON(200, config.SetResToken(200, "Se inicio session correctamente", tokenString, expiresJWT))
}

// ---- REGISTRADORES ----
// ---- registration por WP ----
func RegistrationWp(c echo.Context) error {
	// obteniendo variables
	body := &models.User{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// verificar que la app No este bloqueado, Exista y este activo el usuario
	notblock, exists, _, _, err := models.GetUserAndVerifyNotblockExitsAndActive(body.Phone)
	if !notblock {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "La api de acapela.shop esta en mantenimiento, por favor intentalo mas tarde.")
		return c.JSON(500, config.SetRes(500, "Error: No se completa el proceso por que la api esta en mantenimiento"))
	} else if exists {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "Pero si tu ya estas registrado, no podemos registrarte dos veces.")
		return c.JSON(400, config.SetRes(400, "Error: Ya existe el numero de telefono"))
	} else if err != nil {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "Hubo un problema intentalo de nuevo.")
		return c.JSON(500, config.SetResError(500, "Error: no se pudo terminar la consulta", err.Error()))
	}
	// crear usuario
	err = models.AutoClientRegistrar(body.Phone, body.Name)
	if err != nil {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "No se pudo registrar tu número comunicate por whatsapp al número 69804340.")
		return c.JSON(500, config.SetResError(500, "Error: al crear cliente en la BBDD", err.Error()))
	}
	// enviar el mensaje de bienvenida
	err = middlewares.SendWelcomeMessage(strconv.Itoa(body.Phone), body.Name)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: no se pudo enviar el mensaje de bienvenida", err.Error()))
	}
	return c.JSON(200, config.SetRes(200, "Usuario creado"))
}

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
	existPhone := models.ExistsPhone(body.Phone)
	if existPhone {
		return c.JSON(400, config.SetRes(400, "Error: El telefono ya ha sido registrado."))
	}
	// crear usuario en BBDD
	err := models.ClientRegistrar(body.Name, body.Phone)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: No se ha registrado al cliente", err.Error()))
	}
	// Capitalizando nombre
	name, err := capitalize.Capitalize(body.Name)
	if err != nil {
		name = body.Name
	}
	// enviar el mesaje de bienvenida
	err = middlewares.SendWelcomeMessage(strconv.Itoa(body.Phone), name)
	if err != nil {
		return c.JSON(200, config.SetResError(500, "Error: usario fue registrado en BBDD pero no se envio el mensaje de bienvenida", err.Error()))
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
	existPhone := models.ExistsPhone(body.Phone)
	if existPhone {
		return c.JSON(400, config.SetRes(400, "Error: El telefono ya ha sido registrado."))
	}
	// crear ususario en BBDD
	err := models.EmployRegistrar(body.Name, body.Phone)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: No se ha registrado al empleado.", err.Error()))
	}
	// Capitalizando nombre
	name, err := capitalize.Capitalize(body.Name)
	if err != nil {
		name = body.Name
	}
	// enviar el mesaje de bienvenida
	err = middlewares.SendWelcomeMessage(strconv.Itoa(body.Phone), name)
	if err != nil {
		return c.JSON(200, config.SetResError(500, "Error: usario fue registrado en BBDD pero no se envio el mensaje de bienvenida", err.Error()))
	}

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
	existPhone := models.ExistsPhone(body.Phone)
	if existPhone {
		return c.JSON(400, config.SetRes(400, "Error: El telefono ya ha sido registrado."))
	}
	// crear ususario en BBDD
	err := models.AdminEmployRegistrar(body.Name, body.Phone)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: No se ha registrado al empleado.", err.Error()))
	}
	// Capitalizando nombre
	name, err := capitalize.Capitalize(body.Name)
	if err != nil {
		name = body.Name
	}
	// enviar el mesaje de bienvenida
	err = middlewares.SendWelcomeMessage(strconv.Itoa(body.Phone), name)
	if err != nil {
		return c.JSON(200, config.SetResError(500, "Error: usario fue registrado en BBDD pero no se envio el mensaje de bienvenida", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "Empleado administrativo registrado."))
}
