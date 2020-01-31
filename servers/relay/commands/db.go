package commands

import (
	"github.com/iden3/go-iden3-servers/cmd"
	"github.com/iden3/go-iden3-servers/config"
	"github.com/urfave/cli"
)

var DbCommands = []cli.Command{{
	Name:  "db",
	Usage: "operate with database",
	Subcommands: []cli.Command{
		{
			Name:  "rawdump",
			Usage: "dump database raw key values",
			Action: cmd.WithCfg(func(c *cli.Context, cfg *config.Config) error {
				return cmd.CmdDbRawDump(c, cfg.Storage.Path)
			}),
		},
		{
			Name:  "rawimport",
			Usage: "import database raw from dumped key values",
			Action: cmd.WithCfg(func(c *cli.Context, cfg *config.Config) error {
				return cmd.CmdDbRawImport(c, cfg.Storage.Path)
			}),
		},
		{
			Name:  "ipfsexport",
			Usage: "export database values to ipfs",
			Action: cmd.WithCfg(func(c *cli.Context, cfg *config.Config) error {
				return cmd.CmdDbIPFSexport(c, cfg.Storage.Path)
			}),
		},
	},
}}
