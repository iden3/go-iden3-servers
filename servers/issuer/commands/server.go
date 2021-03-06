package commands

import (
	"github.com/urfave/cli"

	"github.com/iden3/go-iden3-servers/cmd"
	"github.com/iden3/go-iden3-servers/config"
	"github.com/iden3/go-iden3-servers/servers/issuer/endpoint"
)

var ServerCommands = []cli.Command{
	{
		Name:    "init",
		Aliases: []string{},
		Usage:   "create keys and identity for the server",
		Action:  cmd.CmdNewIssuer,
	},
	{
		Name:    "start",
		Aliases: []string{},
		Usage:   "start the server",
		Action: cmd.WithCfg(func(c *cli.Context, cfg *config.Config) error {
			return cmd.CmdStart(c, cfg, endpoint.Serve)
		}),
	},
	{
		Name:    "sync",
		Aliases: []string{},
		Usage:   "sync the identity state with the smart contract",
		Action:  cmd.WithCfg(cmd.CmdSync),
	},
	// {
	// 	Name:    "stop",
	// 	Aliases: []string{},
	// 	Usage:   "stops the server",
	// 	Action:  cmd.WithCfg(cmd.CmdStop),
	// },
	// {
	// 	Name:    "info",
	// 	Aliases: []string{},
	// 	Usage:   "server status",
	// 	Action:  cmd.WithCfg(cmd.CmdInfo),
	// },
	{
		Name:  "eth",
		Usage: "create and manage eth wallet",
		Subcommands: []cli.Command{{
			Name:    "new",
			Aliases: []string{},
			Usage:   "create new Eth Account Address",
			Action:  cmd.CmdNewEthAccount,
		},
			{
				Name:    "import",
				Aliases: []string{},
				Usage:   "import Eth Account Private Key",
				Action:  cmd.CmdImportEthAccount,
			}},
	},
}
