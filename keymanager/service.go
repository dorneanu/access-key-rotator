package keymanager

import (
	"context"


	"github.com/dorneanu/go-key-rotator/entity"
)

// KeyManager defines methods for rotating an AccessKey
type KeyManager interface {
	ListAccessKeys(ctx context.Context) ([]entity.AccessKey, error)
	CreateAccessKey(ctx context.Context) (entity.AccessKey, error)
	DeleteAccessKey(ctx context.Context, id string) error
}
