package secretsstore

import (
	"context"
	"crypto/rand"
	"testing"
	"time"

	"github.com/alecthomas/assert"
	"github.com/dorneanu/go-key-rotator/entity"
	"github.com/dorneanu/go-key-rotator/mocks"
	"github.com/google/go-github/v33/github"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/nacl/box"
)

func TestGithubSecretsStore_ListRepoSecrets(t *testing.T) {
	// mock client to github secrets service
	mock_secretsservice := &mocks.GithubSecretsService{}

	// Create new secrets store
	github_store := GithubSecretsStore{
		repo_owner:    "dorneanu",
		repo_name:     "test",
		secretsClient: mock_secretsservice,
	}

	expected_keys := []entity.AccessKey{
		{ID: "A"},
		{ID: "B"},
	}

	// Create github secrets
	github_secrets := &github.Secrets{
		TotalCount: 2,
		Secrets: []*github.Secret{
			{Name: "A",
				CreatedAt: github.Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)},
				UpdatedAt: github.Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)},
			},
			{Name: "B",
				CreatedAt: github.Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)},
				UpdatedAt: github.Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)},
			},
		},
	}

	mock_secretsservice.On(
		"ListRepoSecrets",
		mock.Anything,
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("*github.ListOptions")).Return(github_secrets, &github.Response{}, nil).Once()

	// Create new ticket
	secrets, err := github_store.ListSecrets(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, 2, len(secrets))
	assert.Equal(t, expected_keys, secrets)

}
func TestGithubSecretsStore_CreateSecret(t *testing.T) {
	t.Run("Create secret using existing encrypted key", func(t *testing.T) {
		// mock client to github secrets service
		mock_secretsservice := &mocks.GithubSecretsService{}

		// Create new secrets store
		github_store := GithubSecretsStore{
			secretsClient: mock_secretsservice,
		}

		encrypted_secret := &github.EncryptedSecret{
			Name:           "SECRET",
			EncryptedValue: "QIv=",
			KeyID:          "SECRET",
		}

		encrypted_key := entity.EncryptedKey{
			ID:     encrypted_secret.KeyID,
			Secret: []byte(encrypted_secret.EncryptedValue),
		}

		mock_secretsservice.On(
			"CreateOrUpdateRepoSecret",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("*github.EncryptedSecret")).Return(&github.Response{}, nil).Once()

		// Create new ticket
		err := github_store.CreateSecret(context.TODO(), encrypted_key)
		assert.Nil(t, err)
	})
	t.Run("Create encrypted key and then upload", func(t *testing.T) {
		// mock client to github secrets service
		mock_secretsservice := &mocks.GithubSecretsService{}

		// Create new secrets store
		github_store := GithubSecretsStore{
			repo_owner:    "dorneanu",
			repo_name:     "test",
			secretsClient: mock_secretsservice,
		}

		access_key := entity.AccessKey{
			ID:     "SECRET",
			Secret: "SUPER SECRET VALUE",
		}

		// Create pub/priv key pair
		public_key, private_key, _ := box.GenerateKey(rand.Reader)
		assert.Equal(t, 32, len(public_key))

		// Create public key
		pk_id := "secret"
		pk_secret := string(public_key[:])

		// Setup mocks
		mock_secretsservice.On("GetRepoPublicKey",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string")).
			Return(&github.PublicKey{
				KeyID: &pk_id,
				Key:   &pk_secret,
			}, &github.Response{}, nil).Once()

		mock_secretsservice.On(
			"CreateOrUpdateRepoSecret",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("*github.EncryptedSecret")).
			Return(&github.Response{}, nil).Once()

		// First encrypt value
		encrypted_key, err := github_store.EncryptKey(context.TODO(), access_key)
		assert.Nil(t, err)
		assert.Equal(t, access_key.ID, encrypted_key.ID)

		// Decrypt again and check
		decrypted_key, ok := box.OpenAnonymous(nil, encrypted_key.Secret, public_key, private_key)
		assert.Equal(t, true, ok)
		assert.Equal(t, access_key.Secret, string(decrypted_key))

		// Create new ticket
		err = github_store.CreateSecret(context.TODO(), *encrypted_key)
		assert.Nil(t, err)
	})
}

func TestGithubSecretsStore_DeleteSecret(t *testing.T) {
	mock_secretsservice := &mocks.GithubSecretsService{}

	// Create new secrets store
	github_store := GithubSecretsStore{
		repo_owner:    "dorneanu",
		repo_name:     "test",
		secretsClient: mock_secretsservice,
	}
	encrypted_key := entity.EncryptedKey{
		ID:     "SECRET",
		Secret: []byte("SECRET VALUE"),
	}

	mock_secretsservice.On(
		"DeleteRepoSecret",
		mock.Anything,
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string")).
		Return(&github.Response{}, nil).Once()

	err := github_store.DeleteSecret(context.TODO(), encrypted_key)
	assert.Nil(t, err)
}
