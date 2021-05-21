package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dorneanu/go-key-rotator/app"
	"github.com/kelseyhightower/envconfig"
)

// Config for this lambda
type Config struct {
	CloudProvider        string `envconfig:"CLOUD_PROVIDER" required:"true"`
	IamUser              string `envconfig:"IAM_USER" required:"true"`
	SecretsStore         string `envconfig:"SECRETS_STORE" required:"true"`
	SecretName           string `envconfig:"SECRET_NAME" required:"true"`
	RepoOwner            string `envconfig:"REPO_OWNER" required:"true"`
	RepoName             string `envconfig:"REPO_NAME" required:"true"`
	ConfigStoreTokenPath string `envconfig:"TOKEN_CONFIG_STORE_PATH" required:"true"`
}

var conf Config

func init() {
	err := envconfig.Process("", &conf)
	if err != nil {
		log.Panicf("Could not find all required ENV variables: %s", err)
	}
}

// handler implements the business logic to be executed during Lambda invocation
func handler(ctx context.Context) error {
	rotatorApp := app.AccessKeyRotatorAppFactory(app.AccessKeyRotatorSettings{
		CloudProvider:        conf.CloudProvider,
		SecretsStore:         conf.SecretsStore,
		SecretName:           conf.SecretName,
		IamUser:              conf.IamUser,
		RepoOwner:            conf.RepoOwner,
		RepoName:             conf.RepoName,
		ConfigStoreTokenPath: conf.ConfigStoreTokenPath,
	})
	err := rotatorApp.UploadSecrets(context.Background())
	if err == nil {
		log.Printf("Secret(s) of %s (%s) were successfully rotated and uploaded to %s\n",
			conf.IamUser, conf.CloudProvider, conf.SecretsStore)
	}
	return err
}

func main() {
	lambda.Start(handler)
}
