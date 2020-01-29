package commands

import (
	"github.com/iden3/go-iden3-servers/cmd"
	"github.com/urfave/cli"
)

var ClaimCommands = []cli.Command{
	{
		Name:  "claim",
		Usage: "claim add",
		Subcommands: []cli.Command{{
			Name:   "add",
			Usage:  "claim add",
			Action: cmd.WithCfg(cmd.CmdAddClaim),
		}},
	},
	{
		Name:  "claims",
		Usage: "claims import from file",
		Subcommands: []cli.Command{{
			Name:   "fromfile",
			Usage:  "import claims from file",
			Action: cmd.WithCfg(cmd.CmdAddClaimsFromFile),
		}},
	},
}
