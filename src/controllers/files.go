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
	fmt.Println("aqui en upload image")

	// OBTENIENDO FUENTE
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(400, config.SetResError(400, "Error: no hay archivo", err.Error()))
	}
	src, err := file.Open()
	if err != nil {
		return c.JSON(400, config.SetResError(400, "Error: No se puede abrir el archivo", err.Error()))
	}
	defer src.Close()
	// CREANDO DESTINO
	path := filepath.Join("assets", "temp", file.Filename)
	dst, err := os.Create(path)
	if err != nil {
		fmt.Println("al crear ", err)
		return c.JSON(500, config.SetResError(500, "no upload image error to create destination", err.Error()))
	}
	// defer dst.Close()
	// COPIAR
	if _, err = io.Copy(dst, src); err != nil {
		fmt.Println("al copiar ", err)
		return c.JSON(500, config.SetResError(500, "no upload image error to copy", err.Error()))
	}
	// enviar a cloudinary y eliminar imagen
	url, err := middlewares.SendImageInCloudinary(path)
	dst.Close() // cerrar archivo para empezar a eliminar
	e := middlewares.DeleteImageForPath(path)
	if e != nil {
		fmt.Println(e)
	}
	if err != nil {
		return c.JSON(500, config.SetResError(500, "no se pudo enviar la imagen a cloudinary", err.Error()))
	}

	return c.JSON(200, config.SetRes(200, url))
}
