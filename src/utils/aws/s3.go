package aws

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"image/png"
	"math/rand"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/disintegration/imaging"
)

var imgMeta = map[ImgType]imgMetaStruct{
	ImgTypeProfile: {
		bucket:     func() string { return "board-game-app" },
		domain:     func() string { return "board-game-app.s3.ap-south-1.amazonaws.com" },
		path:       "profiles",
		width:      512,
		height:     512,
		expireTime: 24 * time.Hour,
	},
	ImgTypeCard: {
		bucket:     func() string { return "board-game-app" },
		domain:     func() string { return "board-game-app.s3.ap-south-1.amazonaws.com" },
		path:       "cards",
		width:      512,
		height:     512,
		expireTime: 10 * time.Hour,
	},
	ImgTypeBirdCard: {
		bucket:     func() string { return "board-game-app" },
		domain:     func() string { return "board-game-app.s3.ap-south-1.amazonaws.com" },
		path:       "wingspan/images",
		width:      512,
		height:     512,
		expireTime: 10 * time.Hour,
	},
	ImgTypeMission: {
		bucket:     func() string { return "board-game-app" },
		domain:     func() string { return "board-game-app.s3.ap-south-1.amazonaws.com" },
		path:       "wingspan/missions",
		width:      30,
		height:     30,
		expireTime: 10 * time.Hour,
	},
	ImgTypeFrogCard: {
		bucket:     func() string { return "board-game-app" },
		domain:     func() string { return "board-game-app.s3.ap-south-1.amazonaws.com" },
		path:       "frog/images",
		width:      45,
		height:     45,
		expireTime: 10 * time.Hour,
	},
	ImgTypeFindIt: {
		bucket:     func() string { return "board-game-app" },
		domain:     func() string { return "board-game-app.s3.ap-south-1.amazonaws.com" },
		path:       "find-it/images",
		width:      0,
		height:     0,
		expireTime: 10 * time.Hour,
	},
	ImgTypeBoardGame: {
		bucket:     func() string { return "board-game-app" },
		domain:     func() string { return "board-game-app.s3.ap-south-1.amazonaws.com" },
		path:       "board_game/images",
		width:      0,
		height:     0,
		expireTime: 10 * time.Hour,
	},
	ImgTypeSlimeWar: {
		bucket:     func() string { return "board-game-app" },
		domain:     func() string { return "board-game-app.s3.ap-south-1.amazonaws.com" },
		path:       "slime_war/images",
		width:      0,
		height:     0,
		expireTime: 10 * time.Hour,
	},
}

func ImageUpload(ctx context.Context, file *multipart.FileHeader, filename string, imgType ImgType) error {
	meta, ok := imgMeta[imgType]
	if !ok {
		return fmt.Errorf("not available meta info for imgType - %v", imgType)
	}
	bucket := meta.bucket()

	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("fail to open file - %v", err)
	}
	defer src.Close()

	img, err := imaging.Decode(src)
	if err != nil {
		return fmt.Errorf("fail to load image - %v", err)
	}

	// 비율 유지하며 리사이즈
	if meta.width > 0 && meta.height > 0 {
		img = imaging.Fit(img, meta.width, meta.height, imaging.Lanczos)
	} else if meta.width > 0 {
		img = imaging.Resize(img, meta.width, 0, imaging.Lanczos)
	} else if meta.height > 0 {
		img = imaging.Resize(img, 0, meta.height, imaging.Lanczos)
	}

	buf := new(bytes.Buffer)
	if err := imaging.Encode(buf, img, imaging.PNG, imaging.PNGCompressionLevel(png.BestCompression)); err != nil {
		return fmt.Errorf("fail to encode png image - %v", err)
	}

	_, err = awsClientS3Uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(fmt.Sprintf("%s/%s", meta.path, filename)),
		Body:        buf,
		ContentType: aws.String("image/png"),
	})
	if err != nil {
		return fmt.Errorf("fail to upload image to s3 - bucket:%s / key:%s/%s", bucket, meta.path, filename)
	}
	return nil
}

func ImageGetSignedURL(ctx context.Context, fileName string, imgType ImgType) (string, error) {
	meta, ok := imgMeta[imgType]
	if !ok {
		return "", fmt.Errorf("not available meta info for imgType - %v", imgType)
	}
	presignClient := s3.NewPresignClient(awsClientS3)

	key := fmt.Sprintf("%s/%s", meta.path, fileName)
	presignParams := &s3.GetObjectInput{
		Bucket: aws.String(meta.bucket()),
		Key:    aws.String(key),
	}

	presignResult, err := presignClient.PresignGetObject(ctx, presignParams, s3.WithPresignExpires(meta.expireTime))
	if err != nil {
		return "", err
	}

	return html.UnescapeString(presignResult.URL), nil
}

func ImageDelete(ctx context.Context, fileName string, imgType ImgType) error {
	meta, ok := imgMeta[imgType]
	if !ok {
		return fmt.Errorf("not available meta info for imgType - %v", imgType)
	}

	bucket := meta.bucket()
	key := fmt.Sprintf("%s/%s", meta.path, fileName)

	if _, err := awsClientS3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}); err != nil {
		return fmt.Errorf("fail to delete image from s3 - bucket:%s, key:%s", bucket, key)
	}

	return nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func FileNameGenerateRandom() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, 32)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
