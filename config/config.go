package config

import (

	// common3 "github.com/iden3/go-iden3-core/common"
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/urfave/cli"
)

type ConfigContract struct {
	JsonABI string         `validate:"required"`
	Address common.Address `validate:"required"`
}

type ConfigServer struct {
	ServiceApi string `validate:"required"`
	AdminApi   string `validate:"required"`
}

type ConfigKeyStore struct {
	Path     string `validate:"required"`
	Password string `validate:"required"`
}

type ConfigEthKeys struct {
	KDis        common.Address `validate:"required"`
	KReen       common.Address `validate:"required"`
	KUpdateRoot common.Address `validate:"required"`
}

type ConfigContracts struct {
	RootCommits   ConfigContract `validate:"required"`
	Iden3Impl     ConfigContract `validate:"required"`
	Iden3Deployer ConfigContract `validate:"required"`
	Iden3Proxy    ConfigContract `validate:"required"`
}

type Config struct {
	Id        core.ID      `validate:"required"`
	Domain    string       `validate:"required"`
	Namespace string       `validate:"required"`
	Server    ConfigServer `validate:"required"`
	Web3      struct {
		Url string `validate:"required"`
	} `validate:"required"`
	KeyStore     ConfigKeyStore `validate:"required"`
	KeyStoreBaby ConfigKeyStore `validate:"required"`
	Keys         struct {
		Ethereum ConfigEthKeys `validate:"required"`
		BabyJub  struct {
			KOp babyjub.PublicKey `validate:"required"`
		} `validate:"required"`
	} `validate:"required"`
	Contracts ConfigContracts `validate:"required"`
	Storage   struct {
		Path string
	} `validate:"required"`
	Names struct {
		Path string `validate:"required"`
	} `validate:"required"`
	Entitites struct {
		Path string `validate:"required"`
	} `validate:"required"`
}

func LoadFromCliFlag(c *cli.Context, cfg interface{}) error {
	cfgFilePath := c.GlobalString("config")
	if cfgFilePath == "" {
		return fmt.Errorf("No config file path specified")
	}
	bs, err := ioutil.ReadFile(cfgFilePath)
	if err != nil {
		return err
	}
	return Load(string(bs), &cfg)
}

func Load(cfgToml string, cfg interface{}) error {
	if _, err := toml.Decode(cfgToml, cfg); err != nil {
		return err
	}
	validate := validator.New()
	return validate.Struct(cfg)
}
