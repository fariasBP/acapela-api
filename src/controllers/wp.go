package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/labstack/echo/v4"
	"github.com/sethvargo/go-password/password"
)

func SendCode(c echo.Context) error {
	// obteniendo variables
	body := &models.User{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// verificando que existe el usuario
	user, err := models.GetUserByPhone(body.CodePhone, body.Phone)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: no existe el numero de telefono", err.Error()))
	}
	// creando codigo
	cod, err := password.Generate(5, 2, 0, true, false)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: al crear codigo", err.Error()))
	}
	// insertando code a user
	_, err = models.SetCode(user.ID, cod)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: al insertar codigo a BBDD", err.Error()))
	}
	// enviar mensaje del codigo por whatsapp
	err = middlewares.SendCodeMessageWithPhoneString(body.PhoneString, cod)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: al enviar codigo via whatsapp", err.Error()))
	}

	return c.JSON(200, config.SetResJson(200, "Codigo creado", map[string]interface{}{"Code": cod, "date": user.CodeDate.Local()}))
}

func SendMassMessagesFromNewProducts(c echo.Context) error {

	// // obtener todos los kinds y models
	// kinds, err := models.GetAllKinds()
	// if err != nil {
	// 	return c.JSON(500, config.SetResError(500, "Error: no se pudo obtner los Kinds de la BBDD", err.Error()))
	// }
	// models, err = models.GetAllModels()
	// if err != nil {
	// 	return c.JSON(500, config.SetResError(500, "Error: no se pudo obtner los Models de la BBDD", err.Error()))
	// }
	// obtener todos los usuarios
	users, err := models.GetPhoneNameNotificationsFromUsers()
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: no se pudo obtener a los usuarios de la BBDD", err.Error()))
	}

	// enviar mensaje del codigo por whatsapp
	for _, v := range users {
		err = middlewares.SendNewProduct(users[0].CodePhone, users[0].Phone, users[0].Name, "abrigos", "dama y varón")
		if err != nil {
			fmt.Println("Error: no se pudo enviar el mensaje al usuario: ", v.Name,
				", con el numero:", v.CodePhone, v.Phone, ", por la causa de:", err.Error())
			// return c.JSON(500, config.SetResError(500, "Error: al enviar codigo via whatsapp", err.Error()))
		}
	}

	return c.JSON(200, config.SetRes(200, "Se envio los mesajes masivos correctamente"))
}
