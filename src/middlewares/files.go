package middlewares

import (
	"context"
	"errors"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/rs/xid"
)

// SUBIENDO A CLAUDINARY
func SendImageInCloudinary(file multipart.File) (*uploader.UploadResult, error) {
	// otener variables de entorno
	pathCld, _ := os.LookupEnv("FOLDER_PRODUCTS_CLOUDINARY")

	// creando cliente cloudinary
	ctx, cld := cloudinaryClient()

	// subiendo imagen
	rss, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: pathCld + "/",
	})
	if err != nil {
		return nil, err
		// return c.JSON(500, config.SetResError(500, "no uploaded image", err.Error()))
	}

	return rss, nil
}

func SendImageInSpaces(file *multipart.FileHeader) (string, error) {
	// abriendo el archivo
	src, err := file.Open()
	if err != nil {
		return "", errors.New("No se puede abrir el archivo")
	}
	defer src.Close()

	// creando un id unico para el nombre
	uniqueID := xid.New().String()
	nameImg := uniqueID + "_" + file.Filename

	//creando objeto (imagen)
	nameBucket := os.Getenv("NAME_BUCKET_SPACE")
	object := s3.PutObjectInput{
		Bucket:      aws.String(nameBucket),
		Key:         aws.String(nameImg),
		Body:        src,
		ContentType: aws.String("image/jpeg"),
		ACL:         aws.String("public-read"),
	}

	// creando cliente
	s3Client := myS3Client()

	// Cargando objeto
	_, err = s3Client.PutObject(&object)
	if err != nil {
		return "", errors.New("no se pudo subir el archivo a spaces")
	}

	// creando la url de la imagen
	endpoint := os.Getenv("ENDPOINT_ACAPELA_SPACE")
	url := "https://" + nameBucket + "." + endpoint + "/" + nameImg

	return url, nil
}

// func DeleteImageInCloudinary(publicId string) error {
// 	// autenticandose en claudinary
// 	cld, _ := cloudinary.NewFromParams(cdyCloudName, cdyApiKey, cdyApiSecret)
// 	ctx := context.Background()

// }
func DeleteImageForPath(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

// CREAR CLIENTE DE S3 PARA SPACES
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
	newSession, _ := session.NewSession(s3config)
	// session.NewSession()
	s3Client := s3.New(newSession)
	return s3Client
}

// CREAR CLIENTE DE CLOUDINARY
func cloudinaryClient() (context.Context, *cloudinary.Cloudinary) {
	// otener variables de entorno
	cdyCloudName, _ := os.LookupEnv("CLOUDINARY_CLOUD_NAME")
	cdyApiKey, _ := os.LookupEnv("CLOUDINARY_API_KEY")
	cdyApiSecret, _ := os.LookupEnv("CLOUDINARY_API_SECRET")

	// autenticandose en claudinary
	cld, _ := cloudinary.NewFromParams(cdyCloudName, cdyApiKey, cdyApiSecret)
	ctx := context.Background()

	return ctx, cld
}
