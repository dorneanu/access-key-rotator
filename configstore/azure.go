package configstore

import "context"

type AzureConfigStore struct{}

func NewAzureConfigStore() *AzureConfigStore {
	return &AzureConfigStore{}
}
func (f *AzureConfigStore) GetValue(ctx context.Context, key string) (string, error) {
	panic("not implemented") // TODO: Implement
}
