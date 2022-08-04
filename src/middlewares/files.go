package middlewares

import (
	"context"
	"errors"
	"os"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

// SUBIENDO A CLAUDINARY
func SendImageInCloudinary(path string) (string, error) {
	// otener variables de entorno
	cdyCloudName, defined := os.LookupEnv("CLOUDINARY_CLOUD_NAME")
	if !defined {
		return "", errors.New("no exist var env CLOUDINARY_CLOUD_NAME")
	}
	cdyApiKey, defined := os.LookupEnv("CLOUDINARY_API_KEY")
	if !defined {
		return "", errors.New("no exist var env CLOUDINARY_API_KEY")
	}
	cdyApiSecret, defined := os.LookupEnv("CLOUDINARY_API_SECRET")
	if !defined {
		return "", errors.New("no exist var env CLOUDINARY_API_SECRET")
	}
	// autenticandose en claudinary
	cld, _ := cloudinary.NewFromParams(cdyCloudName, cdyApiKey, cdyApiSecret)
	ctx := context.Background()
	// subiendo imagen
	_, err := cld.Upload.Upload(ctx, path, uploader.UploadParams{})
	if err != nil {
		return "", err
		// return c.JSON(500, config.SetResError(500, "no uploaded image", err.Error()))
	}
	// obteniendo info de la imagen
	my_image, err := cld.Image("docs/sdk/go/apple")
	if err != nil {
		return "", err
		// return c.JSON(500, config.SetResError(500, "dont get image", err.Error()))
	}
	// obteniendo url de la imagen
	url, err := my_image.String()
	if err != nil {
		return "", err
		// return c.JSON(500, config.SetResError(500, "don't get image url", err.Error()))
	}
	// borrando la imagen
	return url, nil
}
func DeleteImageForPath(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}
