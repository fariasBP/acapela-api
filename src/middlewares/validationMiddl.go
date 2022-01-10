package middlewares

import (
	"fmt"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type (
	signupValidations struct {
		Name     string `json:"name" validate:"required,min=4,"`
		Lastname string `json:"lastname" validate:"required,min=4"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}
	loginValidations struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}
	newProductModelValidations struct {
		Name string `json:"email" validate:"required,min=4"`
	}
)

// AUTH VALIDATORS
func SingupValidate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		name, lastname, email, pwd := c.FormValue("name"), c.FormValue("lastname"), c.FormValue("email"), c.FormValue("password")
		fmt.Println(name, lastname, email, pwd)

		v := &signupValidations{Name: name, Lastname: lastname, Email: email, Password: pwd}

		validate := validator.New()
		if errM := validate.Struct(v); errM != nil {
			return c.JSON(400, config.SetResError(400, "error: validator values signup", errM.Error()))
		}

		return next(c)
	}
}
func LoginValidate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		email, pwd := c.FormValue("email"), c.FormValue("password")

		v := &loginValidations{Email: email, Password: pwd}

		validate := validator.New()
		if err := validate.Struct(v); err != nil {
			return c.JSON(400, config.SetResError(400, "error: validator values login", err.Error()))
		}

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
