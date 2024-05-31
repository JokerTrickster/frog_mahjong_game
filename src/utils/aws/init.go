package aws

import (
	"context"

	AwsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"

	"github.com/aws/aws-sdk-go-v2/aws"
)

var AwsClientSsm *ssm.Client

func InitAws() error {
	var awsConfig aws.Config
	var err error

	awsConfig, err = AwsConfig.LoadDefaultConfig(context.TODO(), AwsConfig.WithRegion("ap-northeast-2"))
	if err != nil {
		return err
	}
	AwsClientSsm = ssm.NewFromConfig(awsConfig)

	return nil
}
