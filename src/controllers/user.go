package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/labstack/echo/v4"
)

type (
	IDuser struct {
		Id string `json:"id"`
	}
	RolUser struct {
		Rol string
	}
)

// func GetUser(c echo.Context) error {
// 	fmt.Println("en getUser controller")
// 	id := fmt.Sprintf("%v", c.Get("id"))
// 	// rol := fmt.Sprintf("%v", c.Get("rol"))

// 	user, err := models.GetUserById(id)
// 	if err != nil {
// 		return c.JSON(400, config.SetResError(400, "error get user from models", err.Error()))
// 	}

// 	return c.JSON(200, config.SetResJson(200, "user recived info", user))
// }
func GetAllUsers(c echo.Context) error {
	users, err := models.GetUsers()
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: al extraer todos los usuarios", err.Error()))
	}
	return c.JSON(200, config.SetResJson(200, "Se obtuvo todos los usuarios", users))
}

// ---- ACTUALIZAR ----
// ---- enviar Confirmacion de inactivarse ----
func SendConfirmInactiveWp(c echo.Context) error {
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
	// enviando el mensaje de la confirmacion de inactivacion
	err = middlewares.SendConfirmInactiveUser(strconv.Itoa(body.Phone))
	if err != nil {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "No se puedo devolver una respuesta, intente de nuevo.")
		return c.JSON(500, config.SetResError(500, "Error: no se puedo enviar la plantilla confirmar inactivacion", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, "Se envio correctamente el mensaje"))
}

// ---- inactivar al usuario ----
func InactiveUserWp(c echo.Context) error {
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
	// consultando - inactivando usuario en BBDD
	err = models.UpdInactiveUserByPhone(body.Phone, body.Sleep)
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
	// verificar que la app No este bloqueado, Exista y este activo el usuario
	notblock, exists, active, _, err := models.GetUserAndVerifyNotblockExitsAndActive(body.Phone)
	if !notblock {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "La api de acapela.shop esta en mantenimiento, por favor intentalo mas tarde.")
		return c.JSON(500, config.SetRes(500, "Error: No se completa el proceso por que la api esta en mantenimiento"))
	} else if !exists {
		middlewares.SendDefaultMsgRegistration(strconv.Itoa(body.Phone))
		return c.JSON(400, config.SetRes(400, "Error: no existe el numero de telefono"))
	} else if active {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "No se puede reactivar por que el usuario ya esta activo.")
		return c.JSON(400, config.SetRes(400, "El usuario ya esta activo"))
	} else if err != nil {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "Hubo un problema intentalo de nuevo.")
		return c.JSON(500, config.SetResError(500, "Error: no se pudo terminar la consulta", err.Error()))
	}
	// consultando - inactivando usuario en BBDD
	err = models.UpdReactiveUserByPhone(body.Phone)
	if err != nil {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "No se te pudo reactivar, posiblemente ya estes activo.")
		return c.JSON(500, config.SetResError(500, "Error: no se pudo reactivar al usuario", err.Error()))
	}
	middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "Se reactivo su cuenta.")

	return c.JSON(200, config.SetRes(200, "Se reactivo al usuario correctamente"))
}

// ---- DELETE ----
// ---- enviar Template wp de Confirmacion para darse de baja ----
func SendConfirmDeleteWp(c echo.Context) error {
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
	// enviando el mensaje de la confirmacion de eliminacion
	err = middlewares.SendConfirmDeleteUser(strconv.Itoa(body.Phone))
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
	// verificar que la app No este bloqueado, Exista y este activo el usuario
	notblock, exists, active, user, err := models.GetUserAndVerifyNotblockExitsAndActive(body.Phone)
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
	// verificando que sea cliente
	if user.Rol != 4 {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "No se puede eliminar tu cuente, porque eres un administrador. Solo un super administrador puede eliminar tu cuenta.")
		return c.JSON(500, config.SetResError(500, "Error: no se pudo eliminar el usuario", err.Error()))
	}
	// consultando - eliminando usuario
	err = models.DelUserByPhone(body.Phone)
	if err != nil {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "No se puedo eliminar los datos, intente de nuevo.")
		return c.JSON(500, config.SetResError(500, "Error: no se pudo eliminar el usuario", err.Error()))
	}
	middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "Se elimino correctamente tus datos.")

	return c.JSON(200, config.SetRes(200, "Se elimino al usuario y envio correctamente el mensaje"))
}
