package controllers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/labstack/echo/v4"
)

func FirstConfirmWithApiWhatsappForReciveMessage(c echo.Context) error {
	tokenVerifyWP, _ := os.LookupEnv("WP_VERIFY_TOKEN")

	// 	let mode = req.query["hub.mode"];
	// let token = req.query["hub.verify_token"];
	// let challenge = req.query["hub.challenge"];

	mode := c.QueryParam("hub.mode")
	token := c.QueryParam("hub.verify_token")
	challange := c.QueryParam("hub.challenge")

	fmt.Println(mode, token, challange)

	// Check the mode and token sent are correct
	if mode == "subscribe" && token == tokenVerifyWP {
		// Respond with 200 OK and challenge token from the request
		fmt.Println("WEBHOOK_VERIFIED")
		return c.String(200, challange)
	} else {
		// Responds with '403 Forbidden' if verify tokens do not match
		return c.String(404, challange)
	}
}
func RecivedMessagesWhatsapp(c echo.Context) error {
	body := make(map[string]interface{})
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()

	fmt.Println(body[""])

	return c.JSON(200, config.SetResJson(200, "se ha recibido el mensaje", body))
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
