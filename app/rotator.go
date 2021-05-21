package app

import (
	"context"
	"fmt"
	"log"

	c "github.com/dorneanu/go-key-rotator/configstore"
	"github.com/dorneanu/go-key-rotator/entity"
	k "github.com/dorneanu/go-key-rotator/keymanager"
	s "github.com/dorneanu/go-key-rotator/secretsstore"
	"github.com/kelseyhightower/envconfig"
)

// AccessKeyRotatorSettings holds settings for the rotator application
type AccessKeyRotatorSettings struct {
	CloudProvider        string `envconfig:"CLOUD_PROVIDER"`
	IamUser              string `envconfig:"IAM_USER"`
	SecretsStore         string `envconfig:"SECRETS_STORE"`
	RepoOwner            string `envconfig:"REPO_OWNER"`
	RepoName             string `envconfig:"REPO_NAME"`
	SecretName           string `envconfig:"SECRET_NAME"`
	ConfigStoreTokenPath string `envconfig:"TOKEN_CONFIG_STORE_PATH"`
}

// AccessKeyRotatorApp represents the application/business logic to be used in different contexts (CLI, Lambda etc.)
type AccessKeyRotatorApp struct {
	KeyManager   k.KeyManager
	ConfigStore  c.ConfigStore
	SecretsStore s.SecretsStore
}

// AccessKeyRotatorAppFactory will setup an AccessKeyRotatorApp depending on the specified cloud provider
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
		privateKey, err := configStore.GetValue(context.Background(), settings.ConfigStoreTokenPath)
		if err != nil {
			log.Fatalf("Uable to get value from config store: %s", err)
		}

		var githubSettings s.GithubAppSettings
		err = envconfig.Process("", &githubSettings)
		if err != nil {
			log.Fatalf("Couldn't get ENV variables for github settings: %s", err)
		}
		githubSettings.PrivateKey = []byte(privateKey)

		githubSecretsClient := s.NewGithubClientAsApp(githubSettings)
		secretsStore = s.NewGithubSecretsStore(
			githubSecretsClient, settings.RepoOwner, settings.RepoName, settings.SecretName)
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

	_, err := a.KeyManager.RotateAccessKey(ctx, access_key_id)
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
		return fmt.Errorf("Couldn't get list of keys: %s", err)
	}

	// Encrypt each key and upload to secrets store
	for _, k := range keys {
		newKey, err := a.KeyManager.RotateAccessKey(ctx, k.ID)
		if err != nil {
			return fmt.Errorf("Couldn't rotate key: %s", err)
		}

		encryptedKey, err := a.SecretsStore.EncryptKey(ctx, newKey)
		if err != nil {
			return fmt.Errorf("Couldn't encrypt key: %s\n", err)
		}

		err = a.SecretsStore.CreateSecret(ctx, *encryptedKey)
		if err != nil {
			return fmt.Errorf("Couldn't upload secrets: %s", err)
		}
	}
	return nil
}
