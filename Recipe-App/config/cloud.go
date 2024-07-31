package config

import (
	"bytes"
	"context"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func CloudConnect(imageBytes []byte) (string, string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("error load .env file")
		panic(err)
	}

	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	id := uuid.New()
	idString := id.String()

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		log.Println("not created cloudinary file")
		panic(err)
	}

	var ctx = context.Background()
	resp, err := cld.Upload.Upload(ctx, bytes.NewReader(imageBytes), uploader.UploadParams{PublicID: idString})
	if err != nil {
		log.Println("error loading image file")
		panic(err)
	}
	url := GetPhoto(idString)
	log.Println(resp)
	return idString, url, nil

}

func GetPhoto(image string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("error load .env file")
		panic(err)
	}

	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		log.Println("not created cloudinary file")
		panic(err)
	}

	var ctx = context.Background()
	resp, err := cld.Admin.Asset(ctx, admin.AssetParams{PublicID: image})
	if err != nil {
		log.Println("error getting photo:", err) // Log the specific error
		panic(err)
	}

	return resp.SecureURL
}