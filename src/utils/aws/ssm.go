package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func AwsSsmGetParam(path string) (string, error) {
	ctx := context.TODO()
	param, err := AwsClientSsm.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(path),
		WithDecryption: PointerTrue(),
	})
	if err != nil {
		return "", err
	}

	return aws.ToString(param.Parameter.Value), nil
}

func PointerTrue() *bool {
	t := true
	return &t
}
