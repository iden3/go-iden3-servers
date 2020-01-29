package commands

import (
	"github.com/urfave/cli"

	"github.com/iden3/go-iden3-servers/cmd"
	"github.com/iden3/go-iden3-servers/config"
	"github.com/iden3/go-iden3-servers/servers/claimserver/endpoint"
)

var ServerCommands = []cli.Command{
	{
		Name:    "init",
		Aliases: []string{},
		Usage:   "create keys and identity for the server",
		Action:  cmd.WithCfg(cmd.CmdNewIdentity),
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
		Name:    "stop",
		Aliases: []string{},
		Usage:   "stops the server",
		Action:  cmd.WithCfg(cmd.CmdStop),
	},
	{
		Name:    "info",
		Aliases: []string{},
		Usage:   "server status",
		Action:  cmd.WithCfg(cmd.CmdInfo),
	},
}
