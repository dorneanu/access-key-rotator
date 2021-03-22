package keymanager

import (
	"context"
	"errors"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/dorneanu/go-key-rotator/entity"
	"github.com/dorneanu/go-key-rotator/mocks"
	"github.com/stretchr/testify/mock"
)

func TestAWSKeyManager_ListAccessKeys(t *testing.T) {
	mock_iam := mocks.IAMAPI{}

	// Create key manager
	km := AWSKeyManager{
		iam_user:   "test",
		iam_client: &mock_iam,
	}

	expected_keys := []entity.AccessKey{
		{
			ID: "access1",
		},
		{
			ID: "access2",
		},
	}

	mock_iam.On(
		"ListAccessKeys",
		mock.Anything,
		mock.AnythingOfType("*iam.ListAccessKeysInput"),
		mock.Anything).Return(
		func(ctx context.Context, input *iam.ListAccessKeysInput, optFns ...func(*iam.Options)) *iam.ListAccessKeysOutput {
			metaData := []types.AccessKeyMetadata{
				{
					AccessKeyId: aws.String(expected_keys[0].ID),
					Status:      types.StatusTypeActive,
				},
				{
					AccessKeyId: aws.String(expected_keys[1].ID),
					Status:      types.StatusTypeInactive,
				},
			}

			// Construct return struct
			output := iam.ListAccessKeysOutput{
				AccessKeyMetadata: metaData,
			}
			return &output
		},
		func(ctx context.Context, input *iam.ListAccessKeysInput, optFns ...func(*iam.Options)) error {
			return nil
		},
	)

	keys, err := km.ListAccessKeys(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, 2, len(keys))
	assert.Equal(t, keys, expected_keys)
	mock_iam.AssertExpectations(t)
}

func TestAWSKeyManager_CreateAccessKey(t *testing.T) {
	mock_iam := mocks.IAMAPI{}

	// Create key manager
	km := AWSKeyManager{
		iam_user:   "test",
		iam_client: &mock_iam,
	}

	expected_key := entity.AccessKey{
		ID:     "SECRET",
		Secret: "SECRET VALUE",
	}
	access_key := &iam.CreateAccessKeyOutput{
		AccessKey: &types.AccessKey{
			AccessKeyId:     &expected_key.ID,
			SecretAccessKey: &expected_key.Secret,
		},
	}

	mock_iam.On(
		"CreateAccessKey",
		mock.Anything,
		mock.AnythingOfType("*iam.CreateAccessKeyInput"),
		mock.Anything).Return(access_key, nil).Once()

	key, err := km.CreateAccessKey(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, expected_key, key)
	mock_iam.AssertExpectations(t)
}

func TestAWSKeyManager_RotateAccessKey(t *testing.T) {
	// Setup
	access_key := &iam.CreateAccessKeyOutput{
		AccessKey: &types.AccessKey{
			AccessKeyId:     aws.String("SECRET"),
			SecretAccessKey: aws.String("SECRET VALUE"),
		},
	}

	t.Run("RotateAccessKey without errors", func(t *testing.T) {
		mock_iam := mocks.IAMAPI{}

		// Create key manager
		km := AWSKeyManager{
			iam_user:   "test",
			iam_client: &mock_iam,
		}

		mock_iam.On(
			"CreateAccessKey",
			mock.Anything,
			mock.AnythingOfType("*iam.CreateAccessKeyInput"),
			mock.Anything).Return(access_key, nil).Once()

		mock_iam.On(
			"DeleteAccessKey",
			mock.Anything,
			mock.AnythingOfType("*iam.DeleteAccessKeyInput"),
			mock.Anything).Return(&iam.DeleteAccessKeyOutput{}, nil).Once()

		err := km.RotateAccessKey(context.TODO(), "SECRET")
		assert.Nil(t, err)
	})

	t.Run("RotateAccessKey access key not found", func(t *testing.T) {
		mock_iam := mocks.IAMAPI{}

		// Create key manager
		km := AWSKeyManager{
			iam_user:   "test",
			iam_client: &mock_iam,
		}

		mock_iam.On(
			"CreateAccessKey",
			mock.Anything,
			mock.AnythingOfType("*iam.CreateAccessKeyInput"),
			mock.Anything).Return(access_key, nil).Once()

		mock_iam.On(
			"DeleteAccessKey",
			mock.Anything,
			mock.AnythingOfType("*iam.DeleteAccessKeyInput"),
			mock.Anything).Return(&iam.DeleteAccessKeyOutput{}, errors.New("Key not found")).Once()

		err := km.RotateAccessKey(context.TODO(), "SECRET")
		assert.Error(t, err)
	})

	t.Run("RotateAccessKey can't create new key", func(t *testing.T) {
		mock_iam := mocks.IAMAPI{}

		// Create key manager
		km := AWSKeyManager{
			iam_user:   "test",
			iam_client: &mock_iam,
		}

		mock_iam.On(
			"CreateAccessKey",
			mock.Anything,
			mock.AnythingOfType("*iam.CreateAccessKeyInput"),
			mock.Anything).Return(access_key, errors.New("Can't create new key")).Once()

		mock_iam.On(
			"DeleteAccessKey",
			mock.Anything,
			mock.AnythingOfType("*iam.DeleteAccessKeyInput"),
			mock.Anything).Return(&iam.DeleteAccessKeyOutput{}, nil).Once()

		err := km.RotateAccessKey(context.TODO(), "SECRET")
		assert.Error(t, err)
	})
}

func TestAWSKeyManager_DeleteAccessKey(t *testing.T) {
	mock_iam := mocks.IAMAPI{}

	// Create key manager
	km := AWSKeyManager{
		iam_user:   "test",
		iam_client: &mock_iam,
	}

	mock_iam.On(
		"DeleteAccessKey",
		mock.Anything,
		mock.AnythingOfType("*iam.DeleteAccessKeyInput"),
		mock.Anything).Return(&iam.DeleteAccessKeyOutput{}, nil).Once()

	err := km.DeleteAccessKey(context.TODO(), "SECRET")
	assert.Nil(t, err)
}
