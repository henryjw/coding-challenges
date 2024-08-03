package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name: "Simple Memcached Server",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "p",
				Value: 9999,
				Usage: "Port number to run the server",
			},
		},
		Action: func(context *cli.Context) error {
			return run(context.Int("p"))
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Error running app: %v\n", err)
	}
}

func run(portNumber int) error {
	return errors.New(fmt.Sprintf("attempted to run app on port %d. However, this functionality is not yet implemented", portNumber))
}
