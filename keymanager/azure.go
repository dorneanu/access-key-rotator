package keymanager

import (
	"context"

	"github.com/dorneanu/go-key-rotator/entity"
)

// TODO
type AzureKeyManager struct{}

func NewAzureKeyManager() *AzureKeyManager {
	return &AzureKeyManager{}
}

func (a *AzureKeyManager) ListAccessKeys(ctx context.Context) ([]entity.AccessKey, error) {
	panic("not implemented") // TODO: Implement
}

func (a *AzureKeyManager) CreateAccessKey(ctx context.Context) (entity.AccessKey, error) {
	panic("not implemented") // TODO: Implement
}

func (a *AzureKeyManager) DeleteAccessKey(ctx context.Context, id string) error {
	panic("not implemented") // TODO: Implement
}

func (a *AzureKeyManager) RotateAccessKey(ctx context.Context, id string) error {
	panic("not implemented") // TODO: Implement
}
