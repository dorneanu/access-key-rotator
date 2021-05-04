package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dorneanu/go-key-rotator/app"
	"github.com/kelseyhightower/envconfig"
)

// Config for this lambda
type config struct {
	CLOUD_PROVIDER string `required:"true"`
	IAM_USER       string `required:"true"`
	SECRETS_STORE  string `required:"true"`
}

var conf = new(config)

func init() {
	err := envconfig.Process("", &conf)
	if err != nil {
		log.Panic("Could not find all required ENV variables")
	}
}

// handler implements the business logic to be executed during Lambda invocation
func handler(ctx context.Context) error {
	rotatorApp := app.AccessKeyRotatorAppFactory(app.AccessKeyRotatorSettings{
		CloudProvider: conf.CLOUD_PROVIDER,
		IamUser:       conf.IAM_USER,
		SecretsStore:  conf.SECRETS_STORE,
	})
	err := rotatorApp.UploadSecrets(context.Background())
	return err
}

func main() {
	lambda.Start(handler)
}
