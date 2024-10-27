package middlewares

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type (
	// validations
	LoginParams struct {
		Phone int    `json:"phone" validate:"required,gt=1000"`
		Code  string `json:"code" validate:"required,min=5,max=5"`
	}
	RegisterParams struct {
		Name  string `json:"name" validate:"required,lowercase,min=3"`
		Phone int    `json:"phone" validate:"required,number,gt=1000"`
	}
	RegisterWpParams struct {
		Phone int `json:"phone" validate:"required,number,gt=1000"`
	}
	GetCodeParams struct {
		// CodePhone int `json:"code_phone" validate:"required,number,gte=1,lte=1000"`
		Phone int `json:"phone" validate:"required,number,gt=1000"`
	}

	//VALUES
	NotificationValues struct {
		To  string
		Msg string
	}
)

// AUTH VALIDATORS
// ---- validaciono para login ----
func LoginValidate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// obteniendo body json
		body := &LoginParams{}
		data, _ := ioutil.ReadAll(c.Request().Body)
		reader := bytes.NewReader(data)
		_ = json.NewDecoder(reader).Decode(body)
		// estableciendo los argumentos de validacion
		v := &LoginParams{
			Phone: body.Phone,
			Code:  strings.TrimSpace(body.Code),
		}
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

// ---- validacion para registraciones ----
func RegisterValidator(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// obteniendo body json
		body := &RegisterParams{}
		data, _ := ioutil.ReadAll(c.Request().Body)
		reader := bytes.NewReader(data)
		_ = json.NewDecoder(reader).Decode(body)
		// estableciendo los argumentos de validacion
		v := &RegisterParams{
			Name:  strings.TrimSpace(body.Name),
			Phone: body.Phone,
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

// ---- validacion para registraciones ----
func RegisterWpValidator(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// obteniendo body json
		body := &RegisterWpParams{}
		data, _ := ioutil.ReadAll(c.Request().Body)
		reader := bytes.NewReader(data)
		_ = json.NewDecoder(reader).Decode(body)
		// estableciendo los argumentos de validacion
		v := &RegisterWpParams{
			Phone: body.Phone,
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

// ---- validacion para codigos ----
func GetCodeValidate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// obteniendo body json
		body := &GetCodeParams{}
		data, _ := ioutil.ReadAll(c.Request().Body)
		reader := bytes.NewReader(data)
		_ = json.NewDecoder(reader).Decode(body)
		// estableciendo los argumentos de validacion
		v := &GetCodeParams{
			Phone: body.Phone,
		}
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
