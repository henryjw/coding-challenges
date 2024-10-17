package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	serverModule "web-server/server"
)

func main() {
	app := &cli.App{
		Name: "Simple Web Server",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "p",
				Value: 8080,
				Usage: "Port number to run the server on",
			},
		},
		Action: func(context *cli.Context) error {
			server := serverModule.New()
			return server.Start(context.Int("p"))
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Error starting app: %v\n", err)
	}
}
