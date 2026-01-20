package cloudinaryutils

import (
	"context"
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryUtil interface {
	UploadImage(ctx context.Context, image any, uploadParams uploader.UploadParams) (string, error)
}

type cloudinaryUtil struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryUtil() *cloudinaryUtil {
	cld, err := cloudinary.New()
	if err != nil {
		log.Fatal(err)
	}
	cld.Config.URL.Secure = true

	return &cloudinaryUtil{
		cld: cld,
	}
}

func (c *cloudinaryUtil) UploadImage(ctx context.Context, image any, uploadParams uploader.UploadParams) (string, error) {
	_, err := c.cld.Upload.Upload(ctx, image, uploadParams)
	if err != nil {
		return "", err
	}

	qsImg, err := c.cld.Image(uploadParams.PublicID)
	if err != nil {
		return "", err
	}

	imgUrl, err := qsImg.String()
	if err != nil {
		return "", err
	}

	return imgUrl, nil
}
