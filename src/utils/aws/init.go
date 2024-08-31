package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	AwsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

var AwsClientSsm *ssm.Client
var awsClientSes *sesv2.Client

func InitAws() error {
	var awsConfig aws.Config
	var err error

	awsConfig, err = AwsConfig.LoadDefaultConfig(context.TODO(), AwsConfig.WithRegion("ap-northeast-2"))
	if err != nil {
		return err
	}
	AwsClientSsm = ssm.NewFromConfig(awsConfig)
	awsClientSes = sesv2.NewFromConfig(awsConfig)

	err = InitAwsSes()
	if err != nil {
		return err
	}
	return nil
}
