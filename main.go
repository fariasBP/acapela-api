package main

import (
	"fmt"
	"log"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/fariasBP/acapela-api/src/routes"
	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	// midlewares
	e.Use(middleware.CORS())
	middlewares.LoadEnvLocal()
	// validador
	e.Validator = &CustomValidator{validator: validator.New()}
	// estableciendo rutas
	routes.IndexRoute(e)
	routes.AuthRoute(e)
	routes.User(e)
	routes.Product(e)
	routes.ModelRoute(e)
	routes.KindRoute(e)
	routes.Notification(e)
	routes.WPRoute(e)
	routes.Files(e)
	// iniciando server
	err := godotenv.Load()
	if err == nil {
		fmt.Println("load env successful")
		e.Logger.Fatal(e.Start(":3000"))
	} else {
		log.Fatal(config.SetResError(500, "NO SE PUEDE CARGAR VALORES ENV", err.Error()))
	}
	// e.Logger.Fatal(e.Start(":8080"))
}
