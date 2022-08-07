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

// ---- registration ----
func RegistrationWp(c echo.Context) error {
	// obteniendo variables
	body := &models.User{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// verificando que existe el usuario
	exists := models.ExistsPhone(body.Phone)
	if !exists {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "No esta registrado en Acapela.Shop, registrese envian un mensaje con la palabra 'registrarse'.")
		return c.JSON(400, config.SetRes(400, "Error: no existe el numero de telefono"))
	}
	// crear usuario
	err := models.AutoClientRegistrarWithWP(body.Phone)
	if err != nil {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "No se pudo registrar tu número comunicate por whatsapp al número 69804340.")
		return c.JSON(500, config.SetResError(500, "Error: al crear cliente en la BBDD", err.Error()))
	}
	// enviar mensaje del codigo por whatsapp
	err = middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "Su número se registro con exito, solo falta ingresar su nombre para completar con el registro...")
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: al enviar codigo via whatsapp", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "Codigo creado"))
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
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "No esta registrado en Acapela.Shop, registrese envian un mensaje con la palabra 'registrarse'.")
		return c.JSON(400, config.SetRes(400, "Error: no existe el numero de telefono"))
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
	user, err := models.GetUserByPhone(body.Phone)
	if err != nil {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "No se reconoce el comando y no esta registrado en Acapela.Shop, registrese envian un mensaje con la palabra 'registrarse'.")
		return c.JSON(400, config.SetRes(400, "Error: no existe el numero de telefono"))
	}
	// verficar si es un auto registration
	if user.WpRegistration && (time.Now().UTC().Before(user.WpRegistrationDate.Add((time.Minute))) || time.Time{} == user.WpRegistrationDate) {
		err = models.UpdUserNameByPhone(body.Phone, body.Name)
		if err != nil {
			middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "No se pudo registrar el nombre. Si tiene alguna duda o consulta envie un mensaje via whatsapp al 69804340.")
		}
	}
	// enviando mensaje default
	err = middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "No se reconoce el comando. Si tiene alguna duda o consulta envie un mensaje via whatsapp al 69804340.")
	if err != nil {
		return c.JSON(400, config.SetResError(400, "Error: al enviar mensaje whatsapp", err.Error()))
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
