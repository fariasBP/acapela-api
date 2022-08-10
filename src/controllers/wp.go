package controllers

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/labstack/echo/v4"
	"github.com/sethvargo/go-password/password"
)

// ---- AUTENTICACION ----
// ---- registration ----
func RegistrationWp(c echo.Context) error {
	// obteniendo variables
	body := &models.User{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// verificando que existe el usuario
	exists := models.ExistsPhone(body.Phone)
	if exists {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "Pero si tu ya estas registrado, no podemos registrarte dos veces.")
		return c.JSON(400, config.SetRes(400, "Error: no existe el numero de telefono"))
	}
	// crear usuario
	err := models.AutoClientRegistrar(body.Phone, body.Name)
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

// ---- enviar codigo ----
func SendCodeWp(c echo.Context) error {
	// obteniendo variables
	body := &models.User{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// verificando que existe el usuario
	user, err := models.GetUserByPhone(body.Phone)
	if err != nil {
		middlewares.SendDefaultMsgRegistration(strconv.Itoa(body.Phone))
		return c.JSON(400, config.SetRes(400, "Error: no existe el numero de telefono"))
	}
	// verificar si esta inactivo
	active := models.VerifyActiveUserByPhone(body.Phone)
	if !active {
		middlewares.SendReactive(strconv.Itoa(body.Phone))
		return c.JSON(200, config.SetRes(200, "Usuario no esta inactivo enviando mensajes para reactivarse reactive."))
	}
	// verificar si ha pasado 1 hora
	if time.Now().UTC().After(user.CodeDate.Add(time.Minute)) {
		// creando codigo
		cod, err := password.Generate(5, 2, 0, true, false)
		if err != nil {
			middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "No se pudo generar el codigo")
			return c.JSON(500, config.SetResError(500, "Error: al crear codigo", err.Error()))
		}
		// insertando code a user
		_, err = models.SetCodeByPhone(body.Phone, cod)
		if err != nil {
			middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "No se pudo generar el codigo")
			return c.JSON(500, config.SetResError(500, "Error: al insertar codigo a BBDD", err.Error()))
		}
		// enviar mensaje del codigo por whatsapp
		err = middlewares.SendCodeMessage(strconv.Itoa(body.Phone), cod)
		if err != nil {
			return c.JSON(500, config.SetResError(500, "Error: al enviar codigo via whatsapp", err.Error()))
		}

		return c.JSON(200, config.SetRes(200, "Codigo creado"))
	}
	middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "No se puede enviar el código por que ya ha solicitado uno, espere 1 hora para solicitar otro código.")
	return c.JSON(400, config.SetRes(400, "Error: No se puede enviar el codigo por que ya se ha solicitado uno"))
}

// ---- enviar Confirmacion de inactivarse ----
func SendConfirmInactiveWp(c echo.Context) error {
	// obteniendo variables
	body := &models.User{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// verficar si existe el usuario
	exists := models.ExistsPhone(body.Phone)
	if !exists {
		middlewares.SendDefaultMsgRegistration(strconv.Itoa(body.Phone))
		return c.JSON(400, config.SetRes(400, "Error: no existe el numero de telefono"))
	}
	// verificar si esta inactivo
	verify := models.VerifyActiveUserByPhone(body.Phone)
	if !verify {
		middlewares.SendReactive(strconv.Itoa(body.Phone))
		return c.JSON(200, config.SetRes(200, "Usuario esta inactivo"))
	}
	// enviando el mensaje de la confirmacion de inactivacion
	err := middlewares.SendConfirmInactiveUser(strconv.Itoa(body.Phone))
	if err != nil {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "No se puedo devolver una respuesta, intente de nuevo.")
		return c.JSON(500, config.SetResError(500, "Error: no se puedo enviar la plantilla confirmar inactivacion", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "Se envio correctamente el mensaje"))
}

// ---- enviar Confirmacion de darse de baja ----
func SendConfirmDeleteWp(c echo.Context) error {
	// obteniendo variables
	body := &models.User{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// verficar si existe el usuario
	exists := models.ExistsPhone(body.Phone)
	if !exists {
		middlewares.SendDefaultMsgRegistration(strconv.Itoa(body.Phone))
		return c.JSON(400, config.SetRes(400, "Error: no existe el numero de telefono"))
	}
	// verificar si esta inactivo
	verify := models.VerifyActiveUserByPhone(body.Phone)
	if !verify {
		middlewares.SendReactive(strconv.Itoa(body.Phone))
		return c.JSON(200, config.SetRes(200, "Usuario esta inactivo"))
	}
	// enviando el mensaje de la confirmacion de inactivacion
	err := middlewares.SendConfirmDeleteUser(strconv.Itoa(body.Phone))
	if err != nil {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "No se puedo devolver una respuesta, intente de nuevo.")
		return c.JSON(500, config.SetResError(500, "Error: no se puedo enviar la plantilla confirmar darse de baja", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "Se envio correctamente el mensaje"))
}

// ---- eliminar usuario ----
func DeleteUserWp(c echo.Context) error {
	// obteniendo variables
	body := &models.User{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// verficar si existe el usuario
	exists := models.ExistsPhone(body.Phone)
	if !exists {
		middlewares.SendDefaultMsgRegistration(strconv.Itoa(body.Phone))
		return c.JSON(400, config.SetRes(400, "Error: no existe el numero de telefono"))
	}
	// verificar si esta inactivo
	verify := models.VerifyActiveUserByPhone(body.Phone)
	if !verify {
		middlewares.SendReactive(strconv.Itoa(body.Phone))
		return c.JSON(200, config.SetRes(200, "Usuario esta inactivo"))
	}
	// consultando - eliminando usuario
	err := models.DelUserByPhone(body.Phone)
	if err != nil {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "No se puedo eliminar los datos, intente de nuevo.")
		return c.JSON(500, config.SetResError(500, "Error: no se pudo eliminar el usuario", err.Error()))
	}
	middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "Se elimino correctamente tus datos.")

	return c.JSON(200, config.SetRes(200, "Se elimino al usuario y envio correctamente el mensaje"))
}

// ---- inactivar al usuario ----
func InactiveUserWp(c echo.Context) error {
	// obteniendo variables
	body := &models.User{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// verficar si existe el usuario
	exists := models.ExistsPhone(body.Phone)
	if !exists {
		middlewares.SendDefaultMsgRegistration(strconv.Itoa(body.Phone))
		return c.JSON(400, config.SetRes(400, "Error: no existe el numero de telefono"))
	}
	// verificar si esta inactivo
	verify := models.VerifyActiveUserByPhone(body.Phone)
	if !verify {
		middlewares.SendReactive(strconv.Itoa(body.Phone))
		return c.JSON(200, config.SetRes(200, "Usuario esta inactivo"))
	}
	// consultando - inactivando usuario en BBDD
	err := models.UpdInactiveUserByPhone(body.Phone, body.Sleep)
	if err != nil {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "No se te pudo inactivar.")
		return c.JSON(500, config.SetResError(500, "Error: no se pudo inactivar al usuario", err.Error()))
	}

	middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "Se te inactivo correctamente.")
	return c.JSON(200, config.SetRes(200, "Se inactivo al usuario correctamente"))
}

// ---- reactivar al usuario ----
func ReactiveUserWp(c echo.Context) error {
	// obteniendo variables
	body := &models.User{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// verficar si existe el usuario
	exists := models.ExistsPhone(body.Phone)
	if !exists {
		middlewares.SendDefaultMsgRegistration(strconv.Itoa(body.Phone))
		return c.JSON(400, config.SetRes(400, "Error: no existe el numero de telefono"))
	}
	// consultando - inactivando usuario en BBDD
	err := models.UpdReactiveUserByPhone(body.Phone)
	if err != nil {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "No se te pudo reactivar, posiblemente ya estes activo.")
		return c.JSON(500, config.SetResError(500, "Error: no se pudo reactivar al usuario", err.Error()))
	}
	middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "Se reactivo su cuenta.")

	return c.JSON(200, config.SetRes(200, "Se reactivo al usuario correctamente"))
}

// ---- enviar template darse de baja al usuario ----
func DelUserTemplateWp(c echo.Context) error {
	// obteniendo variables
	body := &models.User{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// verficar si existe el usuario
	exists := models.ExistsPhone(body.Phone)
	if !exists {
		middlewares.SendDefaultMsgRegistration(strconv.Itoa(body.Phone))
		return c.JSON(400, config.SetRes(400, "Error: no existe el numero de telefono"))
	}
	// verificar si esta inactivo
	verify := models.VerifyActiveUserByPhone(body.Phone)
	if !verify {
		middlewares.SendReactive(strconv.Itoa(body.Phone))
		return c.JSON(200, config.SetRes(200, "Usuario esta inactivo"))
	}
	// consultando - inactivando usuario en BBDD
	err := models.UpdReactiveUserByPhone(body.Phone)
	if err != nil {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "No se te reactivar.")
		return c.JSON(500, config.SetResError(500, "Error: no se pudo reactivar al usuario", err.Error()))
	}
	middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "Se reactivo su cuenta.")

	return c.JSON(200, config.SetRes(200, "Se reactivo al usuario correctamente"))
}

// ---- UTILITARIOS ----
// ---- enviar link de facebook ----
func SendLinkFacebookWp(c echo.Context) error {
	// obteniendo variables
	body := &models.User{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// enviando el link
	err := middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "https://www.facebook.com/Acapela-Dise%C3%B1o-y-Moda-117114304308121")
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: no se puede enviar el mensaje wp", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "se envio el mensaje"))
}

// ---- Enviar mensaje por defecto ----
func SendDefaultMessageWp(c echo.Context) error {
	// obteniendo variables
	body := &models.User{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// verficar si existe el usuario
	exists := models.ExistsPhone(body.Phone)
	if !exists {
		middlewares.SendDefaultMsgRegistration(strconv.Itoa(body.Phone))
		return c.JSON(400, config.SetRes(400, "Error: no existe el numero de telefono"))
	}
	// verificar si esta inactivo
	active := models.VerifyActiveUserByPhone(body.Phone)
	if !active {
		middlewares.SendReactive(strconv.Itoa(body.Phone))
		return c.JSON(200, config.SetRes(200, "Usuario esta inactivo"))
	}
	// enviando mensaje default
	err := middlewares.SendDefaultMessageNoCommand(strconv.Itoa(body.Phone))
	if err != nil {
		return c.JSON(400, config.SetResError(400, "Error: al enviar mensaje whatsapp", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "Se envio correctamente el mensaje"))
}

// ---- enviar ubicacion ----
func SendLocationMessageWp(c echo.Context) error {
	// obteniendo variables
	body := &models.User{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// verficar si existe el usuario
	exists := models.ExistsPhone(body.Phone)
	if !exists {
		middlewares.SendDefaultMsgRegistration(strconv.Itoa(body.Phone))
		return c.JSON(400, config.SetRes(400, "Error: no existe el numero de telefono"))
	}
	// verificar si esta inactivo
	verify := models.VerifyActiveUserByPhone(body.Phone)
	if !verify {
		middlewares.SendReactive(strconv.Itoa(body.Phone))
		return c.JSON(200, config.SetRes(200, "Usuario esta inactivo"))
	}
	// enviando mensaje location
	err := middlewares.SendLocationMessage(strconv.Itoa(body.Phone))
	if err != nil {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "No se pudo enviar la ubicacion, intente de nuevo.")
		return c.JSON(400, config.SetResError(400, "Error: al enviar la localizacion whatsapp", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "Se envio correctamente la localizacion"))
}

// ---- MAS OPCIONES ----
// ---- enviar mas opciones cero ----
func SendMoreOptWp(c echo.Context) error {
	// obteniendo variables
	body := &models.User{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// verficar si existe el usuario
	exists := models.ExistsPhone(body.Phone)
	if !exists {
		middlewares.SendDefaultMsgRegistration(strconv.Itoa(body.Phone))
		return c.JSON(400, config.SetRes(400, "Error: no existe el numero de telefono"))
	}
	// verificar si esta inactivo
	verify := models.VerifyActiveUserByPhone(body.Phone)
	if !verify {
		middlewares.SendReactive(strconv.Itoa(body.Phone))
		return c.JSON(200, config.SetRes(200, "Usuario esta inactivo"))
	}
	// enviando el mensaje de mas opciones
	err := middlewares.SendMoreOpts(strconv.Itoa(body.Phone))
	if err != nil {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "No se puedo devolver una respuesta, intente de nuevo.")
		return c.JSON(500, config.SetResError(500, "Error: no se puedo enviar la plantilla mas opciones", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "Se envio correctamente el mensaje"))
}

// ---- enviar mas opciones uno ----
func SendMoreOptOneWp(c echo.Context) error {
	// obteniendo variables
	body := &models.User{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// verficar si existe el usuario
	exists := models.ExistsPhone(body.Phone)
	if !exists {
		middlewares.SendDefaultMsgRegistration(strconv.Itoa(body.Phone))
		return c.JSON(400, config.SetRes(400, "Error: no existe el numero de telefono"))
	}
	// verificar si esta inactivo
	active := models.VerifyActiveUserByPhone(body.Phone)
	if !active {
		middlewares.SendReactive(strconv.Itoa(body.Phone))
		return c.JSON(200, config.SetRes(200, "Usuario esta inactivo"))
	}
	// enviando el mensaje de mas opciones uno
	err := middlewares.SendMoreOptsOne(strconv.Itoa(body.Phone))
	if err != nil {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "No se puedo devolver una respuesta, intente de nuevo.")
		return c.JSON(500, config.SetResError(500, "Error: no se puedo enviar la plantilla mas opciones", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "Se envio correctamente el mensaje"))
}

// func SendMassMessagesFromNewProducts(c echo.Context) error {

// 	// // obtener todos los kinds y models
// 	// kinds, err := models.GetAllKinds()
// 	// if err != nil {
// 	// 	return c.JSON(500, config.SetResError(500, "Error: no se pudo obtner los Kinds de la BBDD", err.Error()))
// 	// }
// 	// models, err = models.GetAllModels()
// 	// if err != nil {
// 	// 	return c.JSON(500, config.SetResError(500, "Error: no se pudo obtner los Models de la BBDD", err.Error()))
// 	// }
// 	// obtener todos los usuarios
// 	users, err := models.GetPhoneNameNotificationsFromUsers()
// 	if err != nil {
// 		return c.JSON(500, config.SetResError(500, "Error: no se pudo obtener a los usuarios de la BBDD", err.Error()))
// 	}

// 	// enviar mensaje del codigo por whatsapp
// 	for _, v := range users {
// 		err = middlewares.SendNewProduct(users[0].CodePhone, users[0].Phone, users[0].Name, "abrigos", "dama y varón")
// 		if err != nil {
// 			fmt.Println("Error: no se pudo enviar el mensaje al usuario: ", v.Name,
// 				", con el numero:", v.CodePhone, v.Phone, ", por la causa de:", err.Error())
// 			// return c.JSON(500, config.SetResError(500, "Error: al enviar codigo via whatsapp", err.Error()))
// 		}
// 	}

// 	return c.JSON(200, config.SetRes(200, "Se envio los mesajes masivos correctamente"))
// }
