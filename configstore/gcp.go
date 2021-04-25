package configstore

import "context"

type GCPConfigStore struct{}

func NewGCPConfigStore() *GCPConfigStore {
	return &GCPConfigStore{}
}

func (f *GCPConfigStore) GetValue(ctx context.Context, key string) (string, error) {
	panic("not implemented") // TODO: Implement
}
