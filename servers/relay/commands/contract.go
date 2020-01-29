package commands

import (
	"fmt"

	"github.com/iden3/go-iden3-servers/cmd"
	"github.com/iden3/go-iden3-servers/config"
	"github.com/iden3/go-iden3-servers/loaders"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var ContractCommands = []cli.Command{{
	Name:  "contract",
	Usage: "operate with contracts",
	Subcommands: []cli.Command{
		{
			Name:   "info",
			Usage:  "show information about contracts",
			Action: cmd.WithCfg(cmdContractInfo),
		},
		{
			Name:   "deploy",
			Usage:  "deploy contract",
			Action: cmd.WithCfg(cmdContractDeploy),
		},
	},
}}

func contractInfo(cfgContracts *config.ConfigContracts) map[string]config.ConfigContract {
	var info map[string]config.ConfigContract = make(map[string]config.ConfigContract)
	info["rootcommits"] = cfgContracts.RootCommits
	info["iden3impl"] = cfgContracts.Iden3Impl
	info["iden3deployer"] = cfgContracts.Iden3Deployer
	return info
}

func cmdContractInfo(c *cli.Context, cfg *config.Config) error {
	ks, acc := loaders.LoadKeyStore(&cfg.KeyStore, &cfg.Keys.Ethereum)
	client := loaders.LoadWeb3(ks, &acc, cfg.Web3.Url)

	info := func(name string, info config.ConfigContract) {
		if len(info.Address) > 0 {
			code, err := client.CodeAt(info.Address)
			if err != nil {
				log.Panic(err)
			}
			if len(code) > 0 {
				log.Info(name, ": code set at ", info.Address)
			} else {
				log.Info(name, ": code NOT set at ", info.Address)
			}
		} else {
			log.Info(name, ": address not set")
		}
	}

	for k, v := range contractInfo(&cfg.Contracts) {
		info(k, v)
	}

	return nil
}

func cmdContractDeploy(c *cli.Context, cfg *config.Config) error {
	contractid := c.Args()[0]

	ks, acc := loaders.LoadKeyStore(&cfg.KeyStore, &cfg.Keys.Ethereum)
	client := loaders.LoadWeb3(ks, &acc, cfg.Web3.Url)

	if len(c.Args()) != 1 {
		return fmt.Errorf("should specify contract")
	}

	info, ok := contractInfo(&cfg.Contracts)[contractid]
	if !ok {
		return fmt.Errorf("contract %v does not exist", contractid)
	}
	contract := loaders.LoadContract(client, info.JsonABI, nil)

	_, _, err := contract.DeploySync()
	if err != nil {
		return err
	}

	log.Info("Contract deployed at ", contract.Address().Hex())

	return nil
}
