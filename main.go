package main

import (
	"fmt"

	"github.com/fariasBP/acapela-api/src/routes"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		fmt.Println("error custom validator")
		// return echo.NewHTTPError(400, err.Error())
	}
	return nil
}

func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	routes.IndexRoute(e)
	routes.AuthRoute(e)
	routes.UserRoute(e)
	routes.ProductsRoute(e)
	routes.ModelRoute(e)
	routes.KindRoute(e)

	e.Logger.Fatal(e.Start(":8080"))
}
