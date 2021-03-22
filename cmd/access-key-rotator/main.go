package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	km "github.com/dorneanu/go-key-rotator/keymanager"
	"github.com/urfave/cli"
)

// newKeyManager returns a new KeyManager based on specified cloud provider
func newKeyManager(cloudProvider string) km.KeyManager {
	switch cloudProvider {
	case "aws":
		return km.NewAWSKeyManager("test")
	case "gcp":
		return km.NewGCPKeyManager()
	case "azure":
		return km.NewAzureKeyManager()
	default:
		panic("Unknown cloud provider")
	}

}
func main() {
	var cloudProvider string

	app := &cli.App{
		// Global flags
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "cp",
				Usage:       "Specify cloud provider: aws, gcp, azure",
				Required:    true,
				Destination: &cloudProvider,
			},
		},
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
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "key-id",
						Usage: "Access key ID",
					},
					&cli.StringFlag{
						Name:  "iam-user",
						Usage: "Name of the IAM user",
					},
				},
				Usage: "List available access keys",
				Action: func(c *cli.Context) error {
					keyManager := newKeyManager(cloudProvider)
					keys, _ := keyManager.ListAccessKeys(context.Background())
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
				Usage:   "Rotate access key (per default all will be rotated)",
				Action: func(c *cli.Context) error {
					fmt.Println("completed task: ", c.Args().First())
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
