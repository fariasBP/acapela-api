package controllers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
)

func UploadImage(c echo.Context) error {
	// obteniendo archivo
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(400, config.SetResError(400, "Error: no hay archivo", err.Error()))
	}
	src, err := file.Open()
	if err != nil {
		return c.JSON(400, config.SetResError(400, "Error: No se puede abrir el archivo", err.Error()))
	}
	defer src.Close()

	// creando un id unico para el nombre
	uniqueID := xid.New().String()
	nameImg := uniqueID + "_" + file.Filename

	//creando objeto (imagen)
	nameBucket := os.Getenv("NAME_BUCKET_SPACE")
	object := s3.PutObjectInput{
		Bucket: aws.String(nameBucket),

		Key:         aws.String(nameImg),
		Body:        src,
		ContentType: aws.String("image/jpeg"),
		ACL:         aws.String("public-read"),
	}
	s3Client := myS3Client()

	// Cargando objeto
	_, err = s3Client.PutObject(&object)
	if err != nil {
		return c.JSON(500, config.SetResError(500, "no se pudo subir el archivo a spaces", err.Error()))
	}

	return c.JSON(200, config.SetResJson(200, "Se subio imagen a Spaces correctamente.", map[string]string{
		"name": nameImg,
	}))
}
func UploadImageToCloudinary(c echo.Context) error {
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
func myS3Client() *s3.S3 {
	key := os.Getenv("KEY_ACAPELA_SPACE")
	secret := os.Getenv("SECRET_ACAPELA_SPACE")
	region := os.Getenv("REGION_ACAPELA_SPACE")
	endpoint := os.Getenv("ENDPOINT_ACAPELA_SPACE")

	s3config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:    aws.String("https://" + endpoint),
		Region:      aws.String(region),
	}
	newSession := session.New(s3config)
	s3Client := s3.New(newSession)
	return s3Client
}
