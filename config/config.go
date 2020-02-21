package config

import (

	// common3 "github.com/iden3/go-iden3-core/common"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/urfave/cli"
)

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalText(data []byte) error {
	duration, err := time.ParseDuration(string(data))
	if err != nil {
		return err
	}
	d.Duration = duration
	return nil
}

type Contract struct {
	JsonABI string         `validate:"required"`
	Address common.Address `validate:"required"`
}

type Server struct {
	ServiceApi string `validate:"required"`
	AdminApi   string `validate:"required"`
}

type KeyStore struct {
	Path     string `validate:"required"`
	Password string `validate:"required"`
}

type EthKeys struct {
	KDis        common.Address `validate:"required"`
	KReen       common.Address `validate:"required"`
	KUpdateRoot common.Address `validate:"required"`
}

type Contracts struct {
	IdenStates Contract `validate:"required"`
	// Iden3Impl     Contract `validate:"required"`
	// Iden3Deployer Contract `validate:"required"`
	// Iden3Proxy    Contract `validate:"required"`
}

type KeysBabyJub struct {
	KOp babyjub.PublicKey `validate:"required"`
}

type Identity struct {
	Id   core.ID `validate:"required"`
	Keys struct {
		// Ethereum EthKeys `validate:"required"`
		BabyJub KeysBabyJub `validate:"required"`
	} `validate:"required"`
}

type Web3 struct {
	Url string `validate:"required"`
}

type IdenPubOffChain struct {
	Http struct {
		Url string `validate:"required"`
	} `validate:"required"`
}

type Config struct {
	Identity Identity `validate:"required"`
	// Domain    string       `validate:"required"`
	// Namespace string       `validate:"required"`
	Server       Server    `validate:"required"`
	Web3         Web3      `validate:"required"`
	KeyStore     KeyStore  `validate:"required"`
	KeyStoreBaby KeyStore  `validate:"required"`
	Contracts    Contracts `validate:"required"`
	Account      struct {
		Address common.Address `validate:"required"`
	} `validate:"required"`
	Storage struct {
		Path string
	} `validate:"required"`
	Issuer struct {
		PublishStatePeriod        Duration `validate:"required"`
		SyncIdenStatePublicPeriod Duration `validate:"required"`
	}
	IdenPubOffChain IdenPubOffChain `validate:"required"`
	// Names struct {
	// 	Path string `validate:"required"`
	// } `validate:"required"`
	// Entitites struct {
	// 	Path string `validate:"required"`
	// } `validate:"required"`
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
	if err := Load(string(bs), cfg); err != nil {
		return fmt.Errorf("Error loading configuration from cli flag: %w", err)
	}
	return nil
}

func Load(cfgToml string, cfg interface{}) error {
	if _, err := toml.Decode(cfgToml, cfg); err != nil {
		return err
	}
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return fmt.Errorf("Error validating configuration file: %w", err)
	}
	return nil
}
