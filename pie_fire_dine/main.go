package main

import (
	"os"

	"pie_fire_dine/config"
	"pie_fire_dine/logger"
	"pie_fire_dine/server"

	"github.com/urfave/cli"
)

func setup() {
	config.Load()
	logger.SetupLogger(config.GetLogger())
}

func main() {
	setup()

	cliApp := cli.NewApp()
	cliApp.Commands = []cli.Command{
		{
			Name:        "start",
			Description: "Start the service",
			Action: func(c *cli.Context) error {
				server.Start()
				return nil
			},
		},
	}
	// Run app via command line
	if err := cliApp.Run(os.Args); err != nil {
		panic(err)
	}
}
