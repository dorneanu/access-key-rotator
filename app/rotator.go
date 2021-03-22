package app

import (
	"context"
	"fmt"

	"github.com/dorneanu/go-key-rotator/entity"
	k "github.com/dorneanu/go-key-rotator/keymanager"
	s "github.com/dorneanu/go-key-rotator/secretsstore"
)

type AccessKeyRotatorApp struct {
	KeyManager   k.KeyManager
	SecretsStore s.SecretsStore
}

func NewAccessKeyRotatorApp(key_manager k.KeyManager, secrets_store s.SecretsStore) *AccessKeyRotatorApp {
	return &AccessKeyRotatorApp{
		KeyManager:   key_manager,
		SecretsStore: secrets_store,
	}
}

// Rotate will rotate a specified access key
func (a *AccessKeyRotatorApp) Rotate(ctx context.Context, access_key_id string) error {
	if access_key_id == "" {
		return fmt.Errorf("access_key_id is empty")
	}

	err := a.KeyManager.RotateAccessKey(ctx, access_key_id)
	if err != nil {
		return fmt.Errorf("Key rotation failed: %s", err)
	}
	return nil
}

// ListKeys will list all available keys within the key manager
func (a *AccessKeyRotatorApp) ListKeys(ctx context.Context) ([]entity.AccessKey, error) {
	keys, err := a.KeyManager.ListAccessKeys(ctx)
	if err != nil {
		return []entity.AccessKey{}, fmt.Errorf("Couldn't fetch access keys: %s", err)
	}
	return keys, nil
}
