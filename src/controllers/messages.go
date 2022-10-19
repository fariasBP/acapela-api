package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/labstack/echo/v4"
)

func SendMessageToUser(c echo.Context) error {
	// obteniendo variables
	body := &models.Message{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()

	// enviando mensaje wp
	err := middlewares.SendAnyMessageText(strconv.Itoa(body.To), body.Msg)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: No se pudo enviar el mensaje a whatsapp.", err.Error()))
	}
	// creando una respuesta
	err = models.CreateMsgFromAppToUser(body.To, body.Msg)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: No se pudo guardar la respuesta pero se envio el mensaje.", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "Se envio correctamente la respuesta."))
}

func GetUserMessages(c echo.Context) error {
	// obteniendo variables
	body := &models.User{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()

	// consultando
	messages, err := models.GetUserMsgsByPhone(body.Phone)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: No se pudo obtener los mensajes del usuario.", err.Error()))
	}
	return c.JSON(200, config.SetResJson(200, "Se obtuvo los mensajes", messages))
}

func GetUsersMsgNotRead(c echo.Context) error {
	// obteniendo datos
	users, err := models.GetUsersByMailbox()
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: No se pudo obtener los Mailbox's de usuarios.", err.Error()))
	}

	return c.JSON(200, config.SetResJson(200, "Se obtuvieron los Mailbox's de usuarios.", users))
}
