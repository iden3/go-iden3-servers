package main

import (
	"os"

	"github.com/iden3/go-iden3-servers/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "eth-iden3"
	app.Version = "0.0.1-alpha"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config"},
	}

	app.Commands = []cli.Command{
		{
			Name:   "new",
			Usage:  "Create new Eth Account Address",
			Action: cmd.CmdNewEthAccount,
		},
		{
			Name:   "import",
			Usage:  "Import Eth Account Private Key",
			Action: cmd.CmdImportEthAccount,
		},
		{
			Name:  "deploy",
			Usage: "Deploy smart contract",
			Subcommands: []cli.Command{
				{
					Name:   "state",
					Usage:  "Deploy iden3 identity state contract",
					Action: cmd.CmdDeployState,
				},
			},
		},
	}

	log.SetLevel(log.DebugLevel)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
