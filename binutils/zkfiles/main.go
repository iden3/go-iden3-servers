package main

import (
	"fmt"
	"os"

	zkutils "github.com/iden3/go-iden3-core/utils/zk"
	"github.com/iden3/go-iden3-servers/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "zkfiles-iden3"
	app.Version = "0.0.1-alpha"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "path", Required: true},
		cli.StringFlag{Name: "url", Required: false},
		cli.StringFlag{Name: "format", Required: true,
			Usage: fmt.Sprintf("Options: %v, %v, %v", zkutils.ProvingKeyFormatJSON,
				zkutils.ProvingKeyFormatBin, zkutils.ProvingKeyFormatGoBin)},
	}

	app.Commands = []cli.Command{
		{
			Name:   "download",
			Usage:  "Download the zk files from the url into path",
			Action: cmd.CmdDownloadZKFiles,
		},
		{
			Name:   "hash",
			Usage:  "Hash the zk files from the path",
			Action: cmd.CmdHashZKFiles,
		},
		{
			Name:  "downloadhash",
			Usage: "Download the zk files from url into path and hash them",
			Action: func(c *cli.Context) error {
				if err := cmd.CmdDownloadZKFiles(c); err != nil {
					return err
				}
				if err := cmd.CmdHashZKFiles(c); err != nil {
					return err
				}
				return nil
			},
		},
	}

	log.SetLevel(log.DebugLevel)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
