package keymanager

import (
	"context"

	"github.com/dorneanu/go-key-rotator/entity"
)

// TODO
type GCPKeyManager struct{}

func NewGCPKeyManager() *GCPKeyManager {
	return &GCPKeyManager{}
}

func (f *GCPKeyManager) ListAccessKeys(ctx context.Context) ([]entity.AccessKey, error) {
	panic("not implemented") // TODO: Implement
}

func (f *GCPKeyManager) CreateAccessKey(ctx context.Context) (entity.AccessKey, error) {
	panic("not implemented") // TODO: Implement
}

func (f *GCPKeyManager) DeleteAccessKey(ctx context.Context, id string) error {
	panic("not implemented") // TODO: Implement
}

func (f *GCPKeyManager) RotateAccessKey(ctx context.Context, id string) (entity.AccessKey, error) {
	panic("not implemented") // TODO: Implement
}
