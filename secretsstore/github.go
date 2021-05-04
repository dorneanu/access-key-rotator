package secretsstore

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/dorneanu/go-key-rotator/entity"
	"github.com/google/go-github/v34/github"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/oauth2"
)

// GithubAppSettings holds several settings for the github authentication as Github application
type GithubAppSettings struct {
	ApplicationID  int64 `envconfig:"GITHUB_APP_ID" required:"true"`
	InstallationID int64 `envconfig:"GITHUB_INST_ID" required:"true"`
	PrivateKey     []byte
}

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
	repoOwner     string
	repoName      string
	secretName    string
	secretsClient GithubSecretsService
	repoPublicKey *github.PublicKey
}

func NewGithubSecretsStore(secretsService GithubSecretsService, repoOwner, repoName, secretName string) *GithubSecretsStore {
	return &GithubSecretsStore{
		secretsClient: secretsService,
		repoOwner:     repoOwner,
		repoName:      repoName,
		secretName:    secretName,
	}
}

// ListSecrets
func (s *GithubSecretsStore) ListSecrets(ctx context.Context) ([]entity.AccessKey, error) {
	// Fetch repository secrets
	github_secrets, _, err := s.secretsClient.ListRepoSecrets(
		ctx, s.repoOwner, s.repoName, &github.ListOptions{},
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
	b64Encoded := base64.StdEncoding.EncodeToString(k.Secret)
	input := &github.EncryptedSecret{
		Name:           s.secretName,
		EncryptedValue: b64Encoded,
		KeyID:          *s.repoPublicKey.KeyID,
	}
	_, err := s.secretsClient.CreateOrUpdateRepoSecret(ctx, s.repoOwner, s.repoName, input)
	return err
}

// DeleteSecret
func (s *GithubSecretsStore) DeleteSecret(ctx context.Context, k entity.EncryptedKey) error {
	_, err := s.secretsClient.DeleteRepoSecret(ctx, s.repoOwner, s.repoName, k.ID)
	return err
}

// EncryptKey
func (s *GithubSecretsStore) EncryptKey(ctx context.Context, k entity.AccessKey) (*entity.EncryptedKey, error) {
	// First get public key in order to encrypt
	public_key, _, err := s.secretsClient.GetRepoPublicKey(ctx, s.repoOwner, s.repoName)
	if err != nil {
		return nil, err
	}
	s.repoPublicKey = public_key

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

// NewGithubClient returns an implementation of GithubSecretsService using OAUTH tokens
func NewGithubClient(accessToken string) GithubSecretsService {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client.Actions
}

// NewGithubClientAsApp returns an implementation of GithubSecretsServce using a Github Application
func NewGithubClientAsApp(settings GithubAppSettings) GithubSecretsService {
	// Authenticate as Github application
	itr, err := ghinstallation.New(http.DefaultTransport, settings.ApplicationID, settings.InstallationID, settings.PrivateKey)
	if err != nil {
		log.Fatalf("Cannot authenticate as a Github application: %s", err)
	}

	// Use installation transport with client.
	client := github.NewClient(&http.Client{Transport: itr, Timeout: time.Second * 10})
	return client.Actions
}
