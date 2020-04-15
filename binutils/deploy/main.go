package main

import (
	"os"

	"github.com/iden3/go-iden3-servers/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "deploy-iden3"
	app.Version = "0.0.1-alpha"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config"},
	}

	app.Commands = []cli.Command{{
		Name:    "state",
		Aliases: []string{},
		Usage:   "Deploy iden3 identity state contract",
		Action:  cmd.CmdDeployState,
	}}

	log.SetLevel(log.DebugLevel)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
