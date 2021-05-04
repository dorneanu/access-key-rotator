package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dorneanu/go-key-rotator/app"
	"github.com/urfave/cli/v2"
)

var (
	cloudProvider string
	secretsStore  string
	iamUser       string
	accessKeyID   string
	repoOwner     string
	repoName      string
	tokenPath     string
	secretName    string
)

func main() {
	globalFlags := []cli.Flag{
		&cli.StringFlag{
			Name:        "cp",
			Usage:       "Specify cloud provider: aws, gcp, azure",
			Required:    true,
			Destination: &cloudProvider,
			EnvVars:     []string{"CLOUD_PROVIDER"},
		},
		&cli.StringFlag{
			Name:        "secrets-store",
			Required:    true,
			Usage:       "Which secrets store should be used (github, gitlab)",
			Destination: &secretsStore,
			EnvVars:     []string{"SECRETS_STORE"},
		},
	}

	// Create new cli app
	app := &cli.App{
		// Flags: globalFlags,
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Victor Dorneanu",
				Email: "some e-mail",
			},
		},
		Version:  "v0.1",
		Compiled: time.Now(),
		Commands: []*cli.Command{
			{
				// list sub-command
				Name:    "list",
				Aliases: []string{"l"},
				Flags: append([]cli.Flag{
					&cli.StringFlag{
						Name:        "iam-user",
						Usage:       "Name of the IAM user",
						Destination: &iamUser,
						EnvVars:     []string{"IAM_USER"},
					},
				}, globalFlags...),
				Usage: "List available access keys",
				Action: func(c *cli.Context) error {
					rotatorApp := app.AccessKeyRotatorAppFactory(app.AccessKeyRotatorSettings{
						CloudProvider: cloudProvider,
						IamUser:       iamUser,
						SecretsStore:  "github",
					})

					keys, err := rotatorApp.ListKeys(context.Background())
					if err != nil {
						return err
					}
					// Print keys
					for _, k := range keys {
						fmt.Printf("Id: %s\n", k.ID)
					}
					return nil
				},
			},
			{
				// rotate subcommand
				Name:    "rotate",
				Aliases: []string{"r"},
				Flags: append([]cli.Flag{
					&cli.StringFlag{
						Name:        "iam-user",
						Usage:       "Name of the IAM user",
						Destination: &iamUser,
						EnvVars:     []string{"IAM_USER"},
					},
					&cli.StringFlag{
						Name:        "access-key-id",
						Usage:       "Access Key ID",
						Destination: &accessKeyID,
					},
				}, globalFlags...),
				Usage: "Rotate access key (per default all will be rotated)",
				Action: func(c *cli.Context) error {
					rotatorApp := app.AccessKeyRotatorAppFactory(
						app.AccessKeyRotatorSettings{
							CloudProvider: cloudProvider,
							IamUser:       iamUser,
							SecretsStore:  "github",
						})
					err := rotatorApp.Rotate(context.Background(), accessKeyID)
					return err
				},
			},
			{
				// upload subcommand
				Name:    "upload",
				Aliases: []string{"u"},
				Flags: append([]cli.Flag{
					&cli.StringFlag{
						Name:        "iam-user",
						Usage:       "Name of the IAM user",
						Destination: &iamUser,
						EnvVars:     []string{"IAM_USER"},
					},
					&cli.StringFlag{
						Name:        "access-key-id",
						Usage:       "Access Key ID",
						Destination: &accessKeyID,
					},
					&cli.StringFlag{
						Name:        "repo-owner",
						Usage:       "Repository owner",
						Destination: &repoOwner,
						EnvVars:     []string{"REPO_OWNER"},
					},
					&cli.StringFlag{
						Name:        "repo-name",
						Usage:       "Repository name",
						Destination: &repoName,
						EnvVars:     []string{"REPO_NAME"},
					},
					&cli.StringFlag{
						Name:        "token-path",
						Usage:       "Token path in the config store",
						Destination: &tokenPath,
						EnvVars:     []string{"TOKEN_CONFIG_STORE_PATH"},
					},
					&cli.StringFlag{
						Name:        "secret-name",
						Usage:       "Name of the secret to be created/updated",
						Destination: &secretName,
						EnvVars:     []string{"SECRET_NAME"},
					},
				}, globalFlags...),
				Usage: "Upload access key to repo store",
				Action: func(c *cli.Context) error {
					rotatorApp := app.AccessKeyRotatorAppFactory(
						app.AccessKeyRotatorSettings{
							CloudProvider:        cloudProvider,
							SecretsStore:         secretsStore,
							IamUser:              iamUser,
							RepoOwner:            repoOwner,
							RepoName:             repoName,
							SecretName:           secretName,
							ConfigStoreTokenPath: tokenPath,
						})
					err := rotatorApp.UploadSecrets(context.Background())
					return err
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
