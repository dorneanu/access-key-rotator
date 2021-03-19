package secretsstore

import (
	"context"

	"github.com/dorneanu/go-key-rotator/entity"
)

type SecretsStore interface {
	EncryptKey(context.Context, entity.AccessKey) (*entity.EncryptedKey, error)
	ListSecrets(context.Context) ([]entity.AccessKey, error)
	CreateSecret(context.Context, entity.EncryptedKey) error
	DeleteSecret(context.Context, entity.EncryptedKey) error
}
