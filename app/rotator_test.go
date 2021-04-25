package app

import (
	"context"
	"errors"
	"testing"

	"github.com/dorneanu/go-key-rotator/entity"
	"github.com/dorneanu/go-key-rotator/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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
			mock.AnythingOfType("string"),
			mock.Anything).Return(nil).Once()

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
			mock.AnythingOfType("string"),
			mock.Anything).Return(nil).Once()

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
			mock.Anything).Return(errors.New("ROTATE")).Once()

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
			entity.AccessKey{ID: "ID1"},
			entity.AccessKey{ID: "ID2"},
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

func TestUploadSecrets(t *testing.T) {
	mock_key_manager := &mocks.KeyManager{}
	mock_secrets_store := &mocks.SecretsStore{}
	mock_config_store := &mocks.ConfigStore{}

	app := &AccessKeyRotatorApp{
		KeyManager:   mock_key_manager,
		SecretsStore: mock_secrets_store,
		ConfigStore:  mock_config_store,
	}
	err := app.UploadSecrets(context.TODO())
	assert.Nil(t, err)
}
