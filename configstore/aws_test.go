package configstore

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/dorneanu/go-key-rotator/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAWSConfigStore(t *testing.T) {
	t.Run("Get value", func(t *testing.T) {
		mock_ssm := mocks.SSMParameterAPI{}

		// Create AWS config store
		ssm_config := AWSConfigStore{
			ssm_client: &mock_ssm,
		}

		param := &types.Parameter{Value: aws.String("some-value")}
		output := &ssm.GetParameterOutput{
			Parameter: param,
		}
		mock_ssm.On("GetParameter",
			mock.Anything,
			mock.AnythingOfType("*ssm.GetParameterInput"),
			mock.Anything,
		).Return(output, nil)

		value, err := ssm_config.GetValue(context.TODO(), "some-value")
		assert.Nil(t, err)
		assert.Equal(t, *param.Value, value)
	})

	t.Run("Get non-existant value", func(t *testing.T) {
		mock_ssm := mocks.SSMParameterAPI{}

		// Create AWS config store
		ssm_config := AWSConfigStore{
			ssm_client: &mock_ssm,
		}

		mock_ssm.On("GetParameter",
			mock.Anything,
			mock.AnythingOfType("*ssm.GetParameterInput"),
			mock.Anything,
		).Return(&ssm.GetParameterOutput{}, errors.New("Not found"))

		_, err := ssm_config.GetValue(context.TODO(), "some-value")
		assert.Error(t, err)
	})
}
