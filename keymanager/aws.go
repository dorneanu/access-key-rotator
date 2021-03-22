package keymanager

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/dorneanu/go-key-rotator/entity"
)

// We'll define an interface fot the IAM API in order to make testing easy
// This interface will be extended as we go through the different steps
type IAMAPI interface {
	CreateAccessKey(ctx context.Context, params *iam.CreateAccessKeyInput, optFns ...func(*iam.Options)) (*iam.CreateAccessKeyOutput, error)
	ListAccessKeys(ctx context.Context, params *iam.ListAccessKeysInput, optFns ...func(*iam.Options)) (*iam.ListAccessKeysOutput, error)
	DeleteAccessKey(ctx context.Context, params *iam.DeleteAccessKeyInput, optFns ...func(*iam.Options)) (*iam.DeleteAccessKeyOutput, error)
}

type AWSKeyManager struct {
	iam_user   string
	iam_client IAMAPI
}

func NewAWSKeyManager(iam_user string) *AWSKeyManager {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	// Create new IAM client
	iam_client := iam.NewFromConfig(cfg)

	return &AWSKeyManager{
		iam_user:   iam_user,
		iam_client: iam_client,
	}
}

// ListAccessKeys retrieves the IAM access keys for an user
func (m *AWSKeyManager) ListAccessKeys(ctx context.Context) ([]entity.AccessKey, error) {
	var keys []entity.AccessKey
	input := &iam.ListAccessKeysInput{
		MaxItems: aws.Int32(int32(10)),
		UserName: &m.iam_user,
	}

	res, err := m.iam_client.ListAccessKeys(ctx, input)
	if err != nil {
		return nil, err
	}

	// Create slice of AccessKey
	for _, key := range res.AccessKeyMetadata {
		k := entity.AccessKey{
			ID:     *key.AccessKeyId,
			Secret: "",
		}
		keys = append(keys, k)
	}

	return keys, nil
}

// CreateAccessKey
func (m *AWSKeyManager) CreateAccessKey(ctx context.Context) (entity.AccessKey, error) {
	input := &iam.CreateAccessKeyInput{
		UserName: &m.iam_user,
	}
	key, err := m.iam_client.CreateAccessKey(ctx, input)

	if err != nil {
		return entity.AccessKey{}, err
	}

	return entity.AccessKey{
		ID:     *key.AccessKey.AccessKeyId,
		Secret: *key.AccessKey.SecretAccessKey,
	}, nil
}

// RotateAccessKey
func (m *AWSKeyManager) RotateAccessKey(ctx context.Context, id string) error {
	// First delete access key specified by id
	err := m.DeleteAccessKey(ctx, id)
	if err != nil {
		return fmt.Errorf("Couldn't delete key (id = %s): %s", id, err)
	}

	// Create new one
	_, err = m.CreateAccessKey(ctx)
	if err != nil {
		return fmt.Errorf("Couldn't create new key: %s", err)
	}

	return nil
}

// DeleteAccessKey
func (m *AWSKeyManager) DeleteAccessKey(ctx context.Context, id string) error {
	input := &iam.DeleteAccessKeyInput{
		AccessKeyId: &id,
		UserName:    &m.iam_user,
	}
	_, err := m.iam_client.DeleteAccessKey(ctx, input)
	return err
}
