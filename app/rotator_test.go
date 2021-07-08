package app

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/dorneanu/go-key-rotator/entity"
	"github.com/dorneanu/go-key-rotator/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGenerator makes mock generation an easy task
type MockGenerator struct {
	MockKeyManager   *mocks.KeyManager
	MockSecretsStore *mocks.SecretsStore
	MockConfigStore  *mocks.ConfigStore
}

func NewMockGenerator() *MockGenerator {
	mockGenerator := &MockGenerator{}
	mockGenerator.Init()
	return mockGenerator
}

// InitDefaultKeyManager initializes and sets up a default key manager
func (m *MockGenerator) InitDefaultKeyManager() {
	mockKeyManager := &mocks.KeyManager{}

	// ListAccessKeys mock
	mockKeyManager.On(
		"ListAccessKeys",
		mock.Anything,
	).Return([]entity.AccessKey{}, nil).Once()

	// RotateAccessKey mock
	mockKeyManager.On(
		"RotateAccessKey",
		mock.Anything,
		mock.AnythingOfType("string"),
	).Return(entity.AccessKey{ID: "ID1", Secret: "secret"}, nil).Once()

	m.MockKeyManager = mockKeyManager
}

// NewKeyManager initilizes a new key manager without any mock expectations
func (m *MockGenerator) NewKeyManager() {
	m.MockKeyManager = &mocks.KeyManager{}
}

// InitDefaultSecretsStore initializes and sets up a default secrets store
func (m *MockGenerator) InitDefaultSecretsStore() {
	mockSecretsStore := &mocks.SecretsStore{}

	// CreateSecret mock
	mockSecretsStore.On(
		"CreateSecret",
		mock.Anything,
		mock.AnythingOfType("entity.EncryptedKey"),
	).Return(nil).Once()

	// EncryptKey mock
	mockSecretsStore.On(
		"EncryptKey",
		mock.Anything,
		mock.AnythingOfType("entity.AccessKey"),
	).Return(&entity.EncryptedKey{ID: "ID1", Secret: []byte{0x1, 0x2, 0x3}}, nil).Once()

	m.MockSecretsStore = mockSecretsStore
}

// NewSecretsStore initilizes a new secrets store without any mock expectations
func (m *MockGenerator) NewSecretsStore() {
	m.MockSecretsStore = &mocks.SecretsStore{}
}

// InitDefaultConfigStore initializes and sets up a default config store
func (m *MockGenerator) InitDefaultConfigStore() {
	m.NewConfigStore()
}

// NewConfigStore initializes a new config store without any mock expectations
func (m *MockGenerator) NewConfigStore() {
	m.MockConfigStore = &mocks.ConfigStore{}
}

func (m *MockGenerator) GetRotatorApp() *AccessKeyRotatorApp {
	return &AccessKeyRotatorApp{
		KeyManager:   m.MockKeyManager,
		SecretsStore: m.MockSecretsStore,
		ConfigStore:  m.MockConfigStore,
	}
}

func (m *MockGenerator) Init() {
	m.InitDefaultKeyManager()
	m.InitDefaultSecretsStore()
	m.InitDefaultSecretsStore()
}

func TestRotate(t *testing.T) {
	t.Run("Check for errors", func(t *testing.T) {
		mock_key_manager := &mocks.KeyManager{}
		mock_secrets_store := &mocks.SecretsStore{}
		mock_config_store := &mocks.ConfigStore{}

		app := &AccessKeyRotatorApp{
			KeyManager:   mock_key_manager,
			SecretsStore: mock_secrets_store,
			ConfigStore:  mock_config_store,
		}

		mock_key_manager.On(
			"RotateAccessKey",
			mock.Anything,
			mock.AnythingOfType("string")).Return(entity.AccessKey{}, nil).Once()

		err := app.Rotate(context.TODO(), "SECRET")
		assert.Nil(t, err)
	})

	t.Run("Check for empty access key id", func(t *testing.T) {
		mock_key_manager := &mocks.KeyManager{}
		mock_secrets_store := &mocks.SecretsStore{}
		mock_config_store := &mocks.ConfigStore{}

		app := &AccessKeyRotatorApp{
			KeyManager:   mock_key_manager,
			SecretsStore: mock_secrets_store,
			ConfigStore:  mock_config_store,
		}

		mock_key_manager.On(
			"RotateAccessKey",
			mock.Anything,
			mock.AnythingOfType("string")).Return(entity.AccessKey{}, nil).Once()

		err := app.Rotate(context.TODO(), "")
		assert.Error(t, err)
	})

	t.Run("Check for RotateAccessKey error", func(t *testing.T) {
		mock_key_manager := &mocks.KeyManager{}
		mock_secrets_store := &mocks.SecretsStore{}
		mock_config_store := &mocks.ConfigStore{}

		app := &AccessKeyRotatorApp{
			KeyManager:   mock_key_manager,
			SecretsStore: mock_secrets_store,
			ConfigStore:  mock_config_store,
		}

		mock_key_manager.On(
			"RotateAccessKey",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.Anything).Return(entity.AccessKey{}, errors.New("ROTATE")).Once()

		err := app.Rotate(context.TODO(), "SECRET")
		assert.Error(t, err)
	})
}
func TestListKeys(t *testing.T) {
	t.Run("ListKeys non empty", func(t *testing.T) {
		mock_key_manager := &mocks.KeyManager{}
		mock_secrets_store := &mocks.SecretsStore{}
		mock_config_store := &mocks.ConfigStore{}

		app := &AccessKeyRotatorApp{
			KeyManager:   mock_key_manager,
			SecretsStore: mock_secrets_store,
			ConfigStore:  mock_config_store,
		}

		mock_key_manager.On(
			"ListAccessKeys",
			mock.Anything).Return([]entity.AccessKey{
			entity.AccessKey{ID: "ID1", Secret: ""},
			entity.AccessKey{ID: "ID2", Secret: ""},
		}, nil).Once()

		keys, err := app.ListKeys(context.TODO())
		assert.Nil(t, err)
		assert.Equal(t, 2, len(keys))
	})

	t.Run("ListKeys on error", func(t *testing.T) {
		mock_key_manager := &mocks.KeyManager{}
		mock_secrets_store := &mocks.SecretsStore{}
		mock_config_store := &mocks.ConfigStore{}

		app := &AccessKeyRotatorApp{
			KeyManager:   mock_key_manager,
			SecretsStore: mock_secrets_store,
			ConfigStore:  mock_config_store,
		}

		mock_key_manager.On(
			"ListAccessKeys",
			mock.Anything).Return([]entity.AccessKey{}, errors.New("SOME ERROR")).Once()

		_, err := app.ListKeys(context.TODO())
		assert.Error(t, err)
	})

}

// TODO: Implement me
func TestUploadSecrets(t *testing.T) {

	t.Run("Test for normal behaviour", func(t *testing.T) {
		mockGenerator := NewMockGenerator()
		rotatorApp := mockGenerator.GetRotatorApp()
		err := rotatorApp.UploadSecrets(context.TODO())
		assert.NoError(t, err)
	})

	t.Run("Test for missing keys ", func(t *testing.T) {
		mockGenerator := NewMockGenerator()
		mockGenerator.NewKeyManager()
		mockGenerator.MockKeyManager.On(
			"ListAccessKeys",
			mock.Anything).Return(nil, fmt.Errorf("Error")).Once()

		rotatorApp := mockGenerator.GetRotatorApp()
		err := rotatorApp.UploadSecrets(context.TODO())
		assert.Error(t, err)
	})
}
