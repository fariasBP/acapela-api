package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type (
	signupValidations struct {
		Name     string `json:"name" validate:"required,min=4,"`
		Lastname string `json:"lastname" validate:"min=4"`
		Email    string `json:"email" validate:"email"`
		Password string `json:"password" validate:"required"`
		Code     string `json:"code" validate:"required"`
		Phone    string `json:"phone" validate:"required,number"`
	}
	loginValidations struct {
		Code     string `json:"code" validate:"required"`
		Phone    int    `json:"phone" validate:"required,number"`
		Password string `json:"password" validate:"required"`
	}
	registerValidations struct {
		Name     string `json:"name" validate:"required,lowercase"`
		Lastname string `json:"lastname" validate:"lowercase"`
		Code     string `json:"code" validate:"required"`
		Phone    int    `json:"phone" validate:"required,number"`
	}
	newProductModelValidations struct {
		Name string `json:"email" validate:"required,min=4"`
	}
	errorMsg struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
	}
	LoginValues struct {
		Code     string
		Phone    int
		Password string
	}
	RegisterValues struct {
		Name     string
		Lastname string
		Code     string
		Phone    int
	}
)

// AUTH VALIDATORS
func SingupValidate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// obteniendo body json
		body := make(map[string]string)
		err := json.NewDecoder(c.Request().Body).Decode(&body)
		if err != nil {
			return c.JSON(400, errorMsg{
				Code: "400",
				Msg:  "error: body can't read",
			})
		}
		name, lastname, email, pwd, code, phone := body["name"], body["lastname"], body["email"], body["password"], body["code"], body["phone"]
		// estableciendo los argumentos de validacion
		v := &signupValidations{Name: name, Lastname: lastname, Email: email, Password: pwd, Code: code, Phone: phone}
		// realizando valdacion
		validate := validator.New()
		if errM := validate.Struct(v); errM != nil {
			return c.JSON(400, config.SetResError(400, "error: validator values signup", errM.Error()))
		}
		// continuando con el programa
		return next(c)
	}
}
func LoginValidate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// obteniendo body json
		body := &LoginValues{}
		data, _ := ioutil.ReadAll(c.Request().Body)
		reader := bytes.NewReader(data)
		_ = json.NewDecoder(reader).Decode(body)
		// estableciendo los argumentos de validacion
		v := &loginValidations{Code: body.Code, Phone: body.Phone, Password: body.Password}
		// realizando valdacion
		validate := validator.New()
		if err := validate.Struct(v); err != nil {
			return c.JSON(400, config.SetResError(400, "error: validator values in Login", err.Error()))
		}
		// fin del middleware
		c.Request().Body = ioutil.NopCloser(bytes.NewReader([]byte(data)))
		return next(c)
	}
}
func RegisterValidator(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// obteniendo body json
		body := &RegisterValues{}
		data, _ := ioutil.ReadAll(c.Request().Body)
		reader := bytes.NewReader(data)
		_ = json.NewDecoder(reader).Decode(body)
		fmt.Println(body)
		// estableciendo los argumentos de validacion
		v := &registerValidations{
			Name:     body.Name,
			Lastname: body.Lastname,
			Code:     body.Code,
			Phone:    body.Phone,
		}
		// realizando valdacion
		validate := validator.New()
		if err := validate.Struct(v); err != nil {
			return c.JSON(400, config.SetResError(400, "error: validator values in Register", err.Error()))
		}
		// fin del middleware
		c.Request().Body = ioutil.NopCloser(bytes.NewReader([]byte(data)))
		return next(c)
	}
}

// MODELS VALIDATOR
func NewProductModelValidate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.FormValue("name")

		v := &newProductModelValidations{
			Name: name,
		}

		validate := validator.New()
		if err := validate.Struct(v); err != nil {
			return c.JSON(400, config.SetResError(400, "error: validator values new product model", err.Error()))
		}

		return next(c)
	}
}
