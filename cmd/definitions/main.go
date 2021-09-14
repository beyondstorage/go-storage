//go:build tools
// +build tools

package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var app = &cli.App{
	Name:  "definitions",
	Usage: "definitions [service.toml]",
	Before: func(c *cli.Context) error {
		if c.Args().Len() > 1 {
			log.Fatalf("args length should be 0 or 1, actual %d", c.Args().Len())
		}
		return nil
	},
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name: "debug",
		},
	},
	Action: func(c *cli.Context) error {
		if c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
			log.SetReportCaller(true)
		}

		if c.Args().Len() == 0 {
			generateGlobal(NewData())
			return nil
		}

		data := NewData()
		filePath := c.Args().First()
		data.LoadService(filePath)
		generateService(data)

		log.Printf("%s generate finished", filePath)
		return nil
	},
}

func main() {
	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}
