package secretsstore

import (
	"context"

	"github.com/dorneanu/go-key-rotator/entity"
	"github.com/google/go-github/v34/github"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/oauth2"
)

// GithubClient implements GithubSecretsService
type GithubClient struct {
	client *github.Client
}

// GithubSecretsService
type GithubSecretsService interface {
	GetRepoPublicKey(ctx context.Context, owner, repo string) (*github.PublicKey, *github.Response, error)
	CreateOrUpdateRepoSecret(ctx context.Context, owner, repo string, eSecret *github.EncryptedSecret) (*github.Response, error)
	ListRepoSecrets(ctx context.Context, owner, repo string, opts *github.ListOptions) (*github.Secrets, *github.Response, error)
	DeleteRepoSecret(ctx context.Context, owner, repo, name string) (*github.Response, error)
}

// GithubSecretsStore implements a SecretsStore
type GithubSecretsStore struct {
	repo_owner    string
	repo_name     string
	secretsClient GithubSecretsService
}

func NewGithubSecretsStore(secretsService GithubSecretsService, repoOwner, repoName string) *GithubSecretsStore {
	return &GithubSecretsStore{
		secretsClient: secretsService,
		repo_owner:    repoOwner,
		repo_name:     repoName,
	}
}

// ListSecrets
func (s *GithubSecretsStore) ListSecrets(ctx context.Context) ([]entity.AccessKey, error) {
	// Fetch repository secrets
	github_secrets, _, err := s.secretsClient.ListRepoSecrets(
		ctx, s.repo_owner, s.repo_name, &github.ListOptions{},
	)
	if err != nil {
		return nil, err
	}

	access_keys := make([]entity.AccessKey, 0)

	// Convert github secrets to access keys
	for _, secret := range github_secrets.Secrets {
		key := entity.AccessKey{ID: secret.Name}
		access_keys = append(access_keys, key)
	}

	return access_keys, nil
}

// CreateSecret
func (s *GithubSecretsStore) CreateSecret(ctx context.Context, k entity.EncryptedKey) error {
	input := &github.EncryptedSecret{
		Name:           k.ID,
		EncryptedValue: string(k.Secret),
		KeyID:          k.ID,
	}
	_, err := s.secretsClient.CreateOrUpdateRepoSecret(ctx, s.repo_owner, s.repo_name, input)
	return err
}

// DeleteSecret
func (s *GithubSecretsStore) DeleteSecret(ctx context.Context, k entity.EncryptedKey) error {
	_, err := s.secretsClient.DeleteRepoSecret(ctx, s.repo_owner, s.repo_name, k.ID)
	return err
}

// EncryptKey
func (s *GithubSecretsStore) EncryptKey(ctx context.Context, k entity.AccessKey) (*entity.EncryptedKey, error) {
	// First get public key in order to encrypt
	public_key, _, err := s.secretsClient.GetRepoPublicKey(ctx, s.repo_owner, s.repo_name)
	if err != nil {
		return nil, err
	}

	// For a sealed box the public key must be of length 32 bytes
	var pub_key [32]byte
	copy(pub_key[:], *public_key.Key)

	// Create sealed box
	box, err := box.SealAnonymous(nil, []byte(k.Secret), &pub_key, nil)
	if err != nil {
		return nil, err
	}

	return &entity.EncryptedKey{
		ID:     k.ID,
		Secret: box,
	}, nil
}

// NewGithubClient returns an implementation of GithubSecretsService
// TODO: Also support authentication via Github Applications
func NewGithubClient(accessToken string) GithubSecretsService {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client.Actions
}
