package controllers

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/fonini/go-capitalize/capitalize"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/sethvargo/go-password/password"
)

// ---- LOGEADORES ----
// ---- login ----
func Login(c echo.Context) error {
	// obteniendo variables
	body := &middlewares.LoginParams{}
	d := c.Request().Body
	_ = json.NewDecoder(d).Decode(body)
	defer d.Close()
	// buscar usuario por numero
	user, err := models.GetUserByPhone(body.Phone)
	if err != nil {
		return c.JSON(404, config.SetResError(404, "Error: Numero de telefono no registrado.", err.Error()))
	}
	// verficando si el codigo es identico
	if user.Code != body.Code {
		return c.JSON(400, config.SetRes(400, "Error: Codigo incorrecto"))
	}
	// verificar si el codigo se envia antes de 1 hora
	if time.Now().UTC().After(user.CodeDate.Add(time.Hour)) {
		return c.JSON(400, config.SetRes(400, "Error: Esta enviando un codigo caducado."))
	}
	// crear JWT
	tokenString, expiresJWT, tokenErr := middlewares.CreateToken(user.ID.Hex(), uint8(user.Rol))
	if tokenErr != nil {
		return c.JSON(500, config.SetResError(500, "Error: token no creado", tokenErr.Error()))
	}

	return c.JSON(200, config.SetResToken(200, "Se inicio session correctamente", tokenString, expiresJWT))
}

// ---- enviar codigo ----
func SendCodeWpAndEmail(c echo.Context) error {
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
		// enviar mensaje del codigo por email
		if user.Email != "" {
			// err := middlewares.SendEmailCodeTemplate(user.Name, user.Email, "assets/templates/emailTemplate.html", cod)
			err := middlewares.SendEmailBody(user.Name, user.Email, "Hola "+user.Name+". Tu codigo es: "+cod)
			if err != nil {
				log.Printf("No se envio el correo electronio con el codigo por el siguiente motivo: %s\n", err)
			}
		}
		// enviar mensaje del codigo por whatsapp
		err = middlewares.SendCodeMessage(strconv.Itoa(body.Phone), cod)
		if err != nil {
			return c.JSON(500, config.SetResError(500, "Error: al enviar codigo via whatsapp", err.Error()))
		}

		return c.JSON(200, config.SetRes(200, "Codigo creado"))
	}
	middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "No se puede enviar el código por que ya ha solicitado uno, espere 1 hora para solicitar otro código.")
	err = middlewares.SendEmailBody(user.Name, user.Email, "Hola "+user.Name+". Tu no puedes recibir otro codigo por que ya has solicitado uno")
	if err != nil {
		log.Printf("No se envio el correo electronio con el codigo por el siguiente motivo: %s\n", err)
	}
	return c.JSON(400, config.SetRes(400, "Error: No se puede enviar el codigo por que ya se ha solicitado uno"))
}

// verificar token
func ValidateToken(c echo.Context) error {
	// obteniedo el header access-token
	var tkn string = ""
	for name, values := range c.Request().Header {
		if name == "Access-Token" {
			tkn = string(values[0])
		}
	}
	// obteniendo secret de variable de entorno
	secretVal, defined := os.LookupEnv("SECRET_JWT")
	if !defined {
		secretVal = "secreto"
	}
	// verificando token
	secret := []byte(secretVal)
	claims := &middlewares.JwtCustomClaims{}
	token, err := jwt.ParseWithClaims(tkn, claims, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return c.JSON(400, config.SetResError(400, "Error: firma del token, no autorizado", err.Error()))
		}
		return c.JSON(400, config.SetResError(400, "Error: token invalido", err.Error()))
	}
	if !token.Valid {
		return c.JSON(400, config.SetResError(400, "Error: token no autorizado", ""))
	}
	// verificando que el usuario existe
	_, err = models.GetUserByIDStr(claims.Id)
	if err != nil {
		return c.JSON(400, config.SetResError(400, "Error: Id del token incorrecto", err.Error()))
	}
	return c.JSON(200, config.SetRes(200, "Token validado correctamente"))
}

// verificar si es desarrollador
func IsSuperUser(c echo.Context) error {
	// obteniendo variables
	id := c.Get("id").(string)
	// consultando
	user, err := models.GetUserByIDStr(id)
	if err != nil {
		return c.JSON(400, config.SetResError(400, "Error: el usuario no existe.", err.Error()))
	}
	// obteniendo variables .env
	phone, defined := os.LookupEnv("INIT_PHONE_ADMIN")
	if !defined {
		return c.JSON(500, config.SetResError(500, "Error: no se pudo obtener variables .env", ""))
	}
	email, defined := os.LookupEnv("INIT_EMAIL_ADMIN")
	if !defined {
		return c.JSON(500, config.SetResError(500, "Error: no se pudo obtener variables .env", ""))
	}
	// convirtiendo valores
	phoneInt, err := strconv.Atoi(phone)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "Error: no se pudo convertir a entero", err.Error()))
	}
	if user.Phone != phoneInt {
		return c.JSON(400, config.SetResError(400, "Error: no es un super usuario", ""))
	}
	if user.Email != email {
		return c.JSON(400, config.SetResError(400, "Error: no es un super usuario", ""))

	}
	return c.JSON(200, config.SetRes(200, "Es superusuario"))
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
	notblock, exists, _, _, _ := models.GetUserAndVerifyNotblockExitsAndActive(body.Phone)
	if !notblock {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "La api de acapela.shop esta en mantenimiento, por favor intentalo mas tarde.")
		return c.JSON(500, config.SetRes(500, "Error: No se completa el proceso por que la api esta en mantenimiento"))
	} else if exists {
		middlewares.SendAnyMessageText(strconv.Itoa(body.Phone), "Pero si tu ya estas registrado, no podemos registrarte dos veces.")
		return c.JSON(400, config.SetRes(400, "Error: Ya existe el numero de telefono"))
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
// func ClientRegistrar(c echo.Context) error {
// 	// obteniendo variables
// 	body := &models.User{}
// 	d := c.Request().Body
// 	_ = json.NewDecoder(d).Decode(body)
// 	defer d.Close()
// 	// verificar que no exista un code + phone iguales
// 	existPhone := models.ExistsPhone(body.Phone)
// 	if existPhone {
// 		return c.JSON(400, config.SetRes(400, "Error: El telefono ya ha sido registrado."))
// 	}
// 	// crear usuario en BBDD
// 	err := models.ClientRegistrar(body.Name, body.Phone)
// 	if err != nil {
// 		return c.JSON(500, config.SetResError(500, "Error: No se ha registrado al cliente", err.Error()))
// 	}
// 	// Capitalizando nombre
// 	name, err := capitalize.Capitalize(body.Name)
// 	if err != nil {
// 		name = body.Name
// 	}
// 	// enviar el mesaje de bienvenida
// 	err = middlewares.SendWelcomeMessage(strconv.Itoa(body.Phone), name)
// 	if err != nil {
// 		return c.JSON(200, config.SetResError(500, "Error: usario fue registrado en BBDD pero no se envio el mensaje de bienvenida", err.Error()))
// 	}

// 	return c.JSON(200, config.SetRes(200, "Cliente registrado."))
// }

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
