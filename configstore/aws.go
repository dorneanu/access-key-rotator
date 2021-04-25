package configstore

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type SSMParameterAPI interface {
	GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
}

// AWSConfigStore implements a ConfigStore
type AWSConfigStore struct {
	ssm_client SSMParameterAPI
}

func NewAWSConfigStore() *AWSConfigStore {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	// Create new SSM client
	ssm_client := ssm.NewFromConfig(cfg)

	return &AWSConfigStore{
		ssm_client: ssm_client,
	}
}

// GetValue fetches a value from the SSM parameter store
func (s *AWSConfigStore) GetValue(ctx context.Context, key string) (string, error) {
	input := &ssm.GetParameterInput{Name: &key}
	results, err := s.ssm_client.GetParameter(ctx, input)
	if err != nil {
		return "", err
	}

	return *results.Parameter.Value, nil
}
