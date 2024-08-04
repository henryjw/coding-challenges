package main

import (
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
