package controllers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/labstack/echo/v4"
)

func UploadImage(c echo.Context) error {
	cld, define := os.LookupEnv("CLOUD_IMG")
	if !define {
		cld = "cloudinary"
	}
	if cld == "cloudinary" {
		return UploadImageToCloudinary(c)
	} else if cld == "spaces" {
		return UploadImageToSpaces(c)
	}
	return c.JSON(500, config.SetRes(500, "Error: No se sabe a que cloud se enviara la images."))
}

func UploadImageToSpaces(c echo.Context) error {
	// obteniendo archivo
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(400, config.SetResError(400, "Error: no hay archivo", err.Error()))
	}

	// enviar a spaces
	url, err := middlewares.SendImageInSpaces(file)
	if err != nil {
		return c.JSON(400, config.SetResError(400, "Errors: No se pudo subir la imagen a Spaces.", err.Error()))
	}

	return c.JSON(200, config.SetResJson(200, "Se subio imagen a Spaces correctamente.", map[string]string{
		"url": url,
	}))
}
func UploadImageToCloudinary(c echo.Context) error {
	// obteniendo fuente
	file, _, err := c.Request().FormFile("file")
	if err != nil {
		return c.JSON(400, config.SetResError(400, "Error: no hay archivo", err.Error()))
	}

	// enviar a cloudinary
	url, err := middlewares.SendImageInCloudinary(file)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "no se pudo enviar la imagen a cloudinary", err.Error()))
	}

	// return c.JSON(200, config.SetRes(200, url))
	return c.JSON(200, config.SetResJson(200, "Se subio imagen a Cloudinary correctamente.", map[string]string{
		"url": url.URL,
		"id":  url.PublicID,
	}))
}

func UploadImageToHost(c echo.Context) error {
	// obteniendo funte
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(400, config.SetResError(400, "Error: no hay archivo", err.Error()))
	}
	// abriendo el archivo
	src, err := file.Open()
	if err != nil {
		return c.JSON(400, config.SetResError(400, "Error: No se puede abrir el archivo", err.Error()))
	}
	defer src.Close() // indicar que se cierre despues
	// creando destino
	path := filepath.Join("assets", "img", "products", file.Filename)
	dst, err := os.Create(path)
	if err != nil {
		fmt.Println("al crear ", err)
		return c.JSON(500, config.SetResError(500, "no upload image error to create destination", err.Error()))
	}
	defer dst.Close()
	// se copia el archivo al destino
	if _, err = io.Copy(dst, src); err != nil {
		fmt.Println("al copiar ", err)
		return c.JSON(500, config.SetResError(500, "no upload image error to copy", err.Error()))
	}
	// para eliminar imagen
	// e := middlewares.DeleteImageForPath(path)
	// if e != nil {
	// 	fmt.Println(e)
	// }

	// obteniendo url imagen
	mainUrl, _ := os.LookupEnv("APP_NAME")
	url := "https://" + mainUrl + "/assets/img/products/" + file.Filename

	return c.JSON(200, config.SetResJson(200, "Se subio imagen a Cloudinary correctamente.", map[string]string{
		"url": url,
	}))
}
