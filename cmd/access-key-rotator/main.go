package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dorneanu/go-key-rotator/app"
	km "github.com/dorneanu/go-key-rotator/keymanager"
	"github.com/urfave/cli"
)

var (
	cloudProvider string
	iamUser       string
	accessKeyID   string
	repoOwner     string
	repoName      string
)

// newKeyManager returns a new KeyManager based on specified cloud provider
func newKeyManager(cloudProvider string) km.KeyManager {
	switch cloudProvider {
	case "aws":
		return km.NewAWSKeyManager(iamUser)
	case "gcp":
		return km.NewGCPKeyManager()
	case "azure":
		return km.NewAzureKeyManager()
	default:
		panic("Unknown cloud provider")
	}
}

func main() {

	globalFlags := []cli.Flag{
		cli.StringFlag{
			Name:        "cp",
			Usage:       "Specify cloud provider: aws, gcp, azure",
			Required:    true,
			Destination: &cloudProvider,
		},
	}

	// Create new cli app
	app := &cli.App{
		// Flags: globalFlags,
		Authors: []cli.Author{
			cli.Author{
				Name:  "Victor Dorneanu",
				Email: "some e-mail",
			},
		},
		Version:  "v0.1",
		Compiled: time.Now(),
		Commands: []cli.Command{
			{
				// list sub-command
				Name:    "list",
				Aliases: []string{"l"},
				Flags: append([]cli.Flag{
					&cli.StringFlag{
						Name:        "iam-user",
						Usage:       "Name of the IAM user",
						Destination: &iamUser,
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
					},
					&cli.StringFlag{
						Name:        "repo-name",
						Usage:       "Repository name",
						Destination: &repoName,
					},
				}, globalFlags...),
				Usage: "Upload access key to repo store",
				Action: func(c *cli.Context) error {
					rotatorApp := app.AccessKeyRotatorAppFactory(
						app.AccessKeyRotatorSettings{
							CloudProvider: cloudProvider,
							IamUser:       iamUser,
							SecretsStore:  "github",
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
