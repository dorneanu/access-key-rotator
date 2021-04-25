package configstore

import "context"

// ConfigStore reads configuration values from a paramater store (e.g. AWS SSM, GCP Secret Manager)
type ConfigStore interface {
	GetValue(ctx context.Context, key string) (string, error)
}
