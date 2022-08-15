package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/labstack/echo/v4"
	"github.com/sethvargo/go-password/password"
)

// ---- enviar codigo ----
func SendCodeWp(c echo.Context) error {
	// obteniendo variables
	body := &models.User{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// verificando si existe el usuario
	user, err := models.GetUserByPhone(body.Phone)
	if err != nil {
		middlewares.SendDefaultMsgRegistration(strconv.Itoa(body.Phone))
		return c.JSON(400, config.SetRes(400, "Error: no existe el numero de telefono"))
	}
	// verificar si ha pasado 1 hora
	if time.Now().UTC().After(user.CodeDate.Add(time.Hour)) {
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
	// verificar que la app No este bloqueado, Exista y este activo el usuario
	notblock, exists, active, _, err := models.GetUserAndVerifyNotblockExitsAndActive(body.Phone)
	if !notblock {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "La api de acapela.shop esta en mantenimiento, por favor intentalo mas tarde.")
		return c.JSON(500, config.SetRes(500, "Error: No se completa el proceso por que la api esta en mantenimiento"))
	} else if !exists {
		middlewares.SendDefaultMsgRegistration(strconv.Itoa(body.Phone))
		return c.JSON(400, config.SetRes(400, "Error: no existe el numero de telefono"))
	} else if !active {
		middlewares.SendReactive(strconv.Itoa(body.Phone))
		return c.JSON(200, config.SetRes(200, "Usuario esta inactivo"))
	} else if err != nil {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "Hubo un problema intentalo de nuevo.")
		return c.JSON(500, config.SetResError(500, "Error: no se pudo terminar la consulta", err.Error()))
	}
	// enviando mensaje default
	err = middlewares.SendDefaultMessageNoCommand(strconv.Itoa(body.Phone))
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
	// verificar que la app No este bloqueado, Exista y este activo el usuario
	notblock, exists, active, _, err := models.GetUserAndVerifyNotblockExitsAndActive(body.Phone)
	if !notblock {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "La api de acapela.shop esta en mantenimiento, por favor intentalo mas tarde.")
		return c.JSON(500, config.SetRes(500, "Error: No se completa el proceso por que la api esta en mantenimiento"))
	} else if !exists {
		middlewares.SendDefaultMsgRegistration(strconv.Itoa(body.Phone))
		return c.JSON(400, config.SetRes(400, "Error: no existe el numero de telefono"))
	} else if !active {
		middlewares.SendReactive(strconv.Itoa(body.Phone))
		return c.JSON(200, config.SetRes(200, "Usuario esta inactivo"))
	} else if err != nil {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "Hubo un problema intentalo de nuevo.")
		return c.JSON(500, config.SetResError(500, "Error: no se pudo terminar la consulta", err.Error()))
	}
	// enviando mensaje location
	err = middlewares.SendLocationMessage(strconv.Itoa(body.Phone))
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
	// verificar que la app No este bloqueado, Exista y este activo el usuario
	notblock, exists, active, _, err := models.GetUserAndVerifyNotblockExitsAndActive(body.Phone)
	if !notblock {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "La api de acapela.shop esta en mantenimiento, por favor intentalo mas tarde.")
		return c.JSON(500, config.SetRes(500, "Error: No se completa el proceso por que la api esta en mantenimiento"))
	} else if !exists {
		middlewares.SendDefaultMsgRegistration(strconv.Itoa(body.Phone))
		return c.JSON(400, config.SetRes(400, "Error: no existe el numero de telefono"))
	} else if !active {
		middlewares.SendReactive(strconv.Itoa(body.Phone))
		return c.JSON(200, config.SetRes(200, "Usuario esta inactivo"))
	} else if err != nil {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "Hubo un problema intentalo de nuevo.")
		return c.JSON(500, config.SetResError(500, "Error: no se pudo terminar la consulta", err.Error()))
	}
	// enviando el mensaje de mas opciones
	err = middlewares.SendMoreOpts(strconv.Itoa(body.Phone))
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
	// verificar que la app No este bloqueado, Exista y este activo el usuario
	notblock, exists, active, _, err := models.GetUserAndVerifyNotblockExitsAndActive(body.Phone)
	if !notblock {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "La api de acapela.shop esta en mantenimiento, por favor intentalo mas tarde.")
		return c.JSON(500, config.SetRes(500, "Error: No se completa el proceso por que la api esta en mantenimiento"))
	} else if !exists {
		middlewares.SendDefaultMsgRegistration(strconv.Itoa(body.Phone))
		return c.JSON(400, config.SetRes(400, "Error: no existe el numero de telefono"))
	} else if !active {
		middlewares.SendReactive(strconv.Itoa(body.Phone))
		return c.JSON(200, config.SetRes(200, "Usuario esta inactivo"))
	} else if err != nil {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "Hubo un problema intentalo de nuevo.")
		return c.JSON(500, config.SetResError(500, "Error: no se pudo terminar la consulta", err.Error()))
	}
	// enviando el mensaje de mas opciones uno
	err = middlewares.SendMoreOptsOne(strconv.Itoa(body.Phone))
	if err != nil {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "No se puedo devolver una respuesta, intente de nuevo.")
		return c.JSON(500, config.SetResError(500, "Error: no se puedo enviar la plantilla mas opciones", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "Se envio correctamente el mensaje"))
}

func SendImgsNewProductsWp(c echo.Context) error {
	to := c.QueryParam("from")
	fmt.Println(to)
	products, err := models.GetNewProducts()

	if err != nil {
		middlewares.SendAnyMessageText(to, "No se pudo enviar los nuevos productos")
		return c.JSON(500, config.SetResError(500, "Error: No se pudo enviar los nuevos productos", err.Error()))
	}

	var x int = 0

	for _, v := range products {
		err = middlewares.SendImageByLink(to, v.Photos[0])
		if err != nil {
			// middlewares.SendAnyMessageText(to, "No se pudo enviar las fotos de los nuevos productos")
			// return c.JSON(500, config.SetResError(500, "Error: No se pudo enviar las fotos de los nuevos productos", err.Error()))
			x++
		}
	}
	if x > 0 {
		middlewares.SendAnyMessageText(to, ("Hubo un problema: " + strconv.Itoa(x) + " fotos no se pudieron enviar de los nuevos productos."))
		return c.JSON(500, config.SetResError(500, "Error: No se pudo enviar algunas fotos de los nuevos productos", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "Se envio correctamente los nuevos productos"))
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
