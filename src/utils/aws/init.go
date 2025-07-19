package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	AwsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

var AwsClientSsm *ssm.Client
var awsClientSes *sesv2.Client
var awsClientS3 *s3.Client
var awsClientS3Uploader *manager.Uploader
var awsClientS3Downloader *manager.Downloader
var awsS3Signer *s3.PresignClient

type ImgType uint8

const (
	ImgTypeProfile   = ImgType(0)
	ImgTypeCard      = ImgType(1)
	ImgTypeBirdCard  = ImgType(2)
	ImgTypeMission   = ImgType(3)
	ImgTypeFrogCard  = ImgType(4)
	ImgTypeFindIt    = ImgType(5)
	ImgTypeBoardGame = ImgType(6)
	ImgTypeSlimeWar  = ImgType(7)
)

type imgMetaStruct struct {
	bucket     func() string
	domain     func() string
	path       string
	width      int
	height     int
	expireTime time.Duration
}

func InitAws() error {
	var awsConfig aws.Config
	var err error

	awsConfig, err = AwsConfig.LoadDefaultConfig(context.TODO(), AwsConfig.WithRegion("ap-south-1"))
	if err != nil {
		return err
	}
	AwsClientSsm = ssm.NewFromConfig(awsConfig)
	awsClientSes = sesv2.NewFromConfig(awsConfig)
	awsClientS3 = s3.NewFromConfig(awsConfig)
	awsClientS3Uploader = manager.NewUploader(awsClientS3)
	awsClientS3Downloader = manager.NewDownloader(awsClientS3)
	awsClientSes = sesv2.NewFromConfig(awsConfig)
	awsS3Signer = s3.NewPresignClient(awsClientS3)
	err = InitAwsSes()
	if err != nil {
		return err
	}
	return nil
}
