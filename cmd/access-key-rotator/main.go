package main

import (
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
		fmt.Println("it's aws")
		manager := km.NewAWSKeyManager("test")
		return manager
	case "gcp":
		fmt.Println("it's google cloud")
		return km.NewAWSKeyManager("test")
	case "azure":
		fmt.Println("it's azure")
		return km.NewAWSKeyManager("test")
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
					fmt.Println("added task: ", c.Args().First())
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
