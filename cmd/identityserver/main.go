package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/iden3/go-iden3-servers/cmd/identityserver/commands"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "identityserver"
	app.Version = "0.0.1-alpha"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config"},
	}

	app.Commands = []cli.Command{}
	app.Commands = append(app.Commands, commands.ServerCommands...)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
