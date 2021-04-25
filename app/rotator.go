package app

import (
	"context"
	"fmt"
	"log"

	c "github.com/dorneanu/go-key-rotator/configstore"
	"github.com/dorneanu/go-key-rotator/entity"
	k "github.com/dorneanu/go-key-rotator/keymanager"
	s "github.com/dorneanu/go-key-rotator/secretsstore"
)

type AccessKeyRotatorSettings struct {
	CloudProvider string
	IamUser       string
	SecretsStore  string
	RepoOwner     string
	RepoName      string
}

type AccessKeyRotatorApp struct {
	KeyManager   k.KeyManager
	ConfigStore  c.ConfigStore
	SecretsStore s.SecretsStore
}

// AccessKeyRotatorAppFactory will setup an AccessKeyRotatorApp depending on
// the specified cloud provider
func AccessKeyRotatorAppFactory(settings AccessKeyRotatorSettings) *AccessKeyRotatorApp {
	var keyManager k.KeyManager
	var configStore c.ConfigStore
	var secretsStore s.SecretsStore

	// Setup key manager
	switch settings.CloudProvider {
	case "aws":
		keyManager = k.NewAWSKeyManager(settings.IamUser)
		configStore = c.NewAWSConfigStore()
	case "gcp":
		keyManager = k.NewGCPKeyManager()
		configStore = c.NewGCPConfigStore()
	case "azure":
		keyManager = k.NewAzureKeyManager()
		configStore = c.NewAzureConfigStore()
	default:
		panic("Unknown cloud provider")
	}

	// Setup secrets store
	switch settings.SecretsStore {
	case "github":
		// TODO: Put string as ENV variable
		accessToken, err := configStore.GetValue(context.Background(), "github-token")
		if err != nil {
			log.Fatalf("Unable to get value from config store: %s", err)
		}
		githubSecretsClient := s.NewGithubClient(accessToken)
		secretsStore = s.NewGithubSecretsStore(githubSecretsClient, settings.RepoOwner, settings.RepoName)
	default:
		panic("Unknown secrets store")
	}

	return &AccessKeyRotatorApp{
		KeyManager:   keyManager,
		ConfigStore:  configStore,
		SecretsStore: secretsStore,
	}
}

func NewAccessKeyRotatorApp(key_manager k.KeyManager, secrets_store s.SecretsStore, config_store c.ConfigStore) *AccessKeyRotatorApp {
	return &AccessKeyRotatorApp{
		KeyManager:   key_manager,
		SecretsStore: secrets_store,
		ConfigStore:  config_store,
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

func (a *AccessKeyRotatorApp) UploadSecrets(ctx context.Context) error {
	// First get list of keys
	keys, err := a.ListKeys(ctx)
	if err != nil {
		return fmt.Errorf("Couldn't upload secrets: %s", err)
	}

	for _, k := range keys {
		fmt.Printf("%+v\n", k)
		encryptedKey, err := a.SecretsStore.EncryptKey(ctx, k)
		if err != nil {
			return fmt.Errorf("Couldn't encrypt key: %s\n", err)
		}

		fmt.Printf("%+v\n", encryptedKey)
	}
	return nil
}
